package actions

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"math/big"
	"sync"

	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/constants"
	"github.com/teathedev/backend-boilerplate/internal/db"
	"github.com/teathedev/backend-boilerplate/internal/ent"
	"github.com/teathedev/backend-boilerplate/internal/ent/accesstokenkey"
	"github.com/teathedev/backend-boilerplate/internal/filters"
	"github.com/teathedev/pkg/env"
	"github.com/teathedev/pkg/jwt"
	"github.com/teathedev/pkg/logger"
	"github.com/teathedev/backend-boilerplate/types"
)

var (
	activeJWKs        []types.JWK
	activeSigningKeys []*types.JWTSigningKey
	jwkMu             sync.RWMutex
)

func init() {
	log := logger.New("JWKsActions")
	if err := RefreshJWKs(context.Background()); err != nil {
		log.Fatal(
			"Failed to refresh JWKs",
			logger.LogParams{
				"Error": err,
			},
		)
	}
}

// EntAccessTokenKeyToJWK converts ent.AccessTokenKey (model) to types.JWK (general purpose).
func EntAccessTokenKeyToJWK(model *ent.AccessTokenKey) (*types.JWK, error) {
	j, err := jwt.ConvertToJWK(model.PublicPem, model.ID)
	if err != nil {
		return nil, err
	}
	return &types.JWK{
		KTY:       types.JWKKeyType(j.KTY),
		Use:       types.JWKKeyUse(j.Use),
		Algorithm: types.JWKAlgorithm(j.Algorithm),
		KeyID:     j.KeyID,
		N:         j.N,
		E:         j.E,
	}, nil
}

// getJWTKeyEncryptionKey returns the JWT key encryption key from env (base64-encoded, 32 bytes for AES-256-GCM).
func getJWTKeyEncryptionKey() ([]byte, error) {
	b64 := env.GetString(constants.JWTMasterKeyEnv, "")
	if b64 == "" {
		return nil, nil
	}
	key, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, nil
	}
	return key, nil
}

// RefreshJWKs loads keys from DB, decrypts active ones, and populates the in-memory cache.
// Call on startup or when keys are rotated.
func RefreshJWKs(ctx context.Context) error {
	masterKey, err := getJWTKeyEncryptionKey()
	if err != nil || len(masterKey) == 0 {
		return nil
	}

	keys, err := db.Client.AccessTokenKey.
		Query().
		Where(filters.VerifyTokens()).
		Order(accesstokenkey.ByState(), accesstokenkey.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return err
	}

	var jwks []types.JWK
	var signingKeys []*types.JWTSigningKey

	for _, key := range keys {
		jwk, err := EntAccessTokenKeyToJWK(key)
		if err != nil {
			continue
		}
		jwks = append(jwks, *jwk)

		if key.State == types.AccessTokenKeyStatesActive {
			privBytes, err := jwt.DecryptPrivateKey(key.PrivateKeyEncrypted, masterKey)
			if err != nil {
				continue
			}
			priv, err := x509.ParsePKCS1PrivateKey(privBytes)
			if err != nil {
				continue
			}
			signingKeys = append(signingKeys, &types.JWTSigningKey{
				KeyID:      key.ID,
				PrivateKey: priv,
				Algorithm:  types.JWKAlgorithmRS256,
			})
		}
	}

	jwkMu.Lock()
	activeJWKs = jwks
	activeSigningKeys = signingKeys
	jwkMu.Unlock()

	return nil
}

// GetRandomJWK returns a random JWK from the in-memory active set.
func GetRandomJWK() *types.JWK {
	jwkMu.RLock()
	defer jwkMu.RUnlock()

	if len(activeJWKs) == 0 {
		return nil
	}
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(activeJWKs))))
	return &activeJWKs[i.Int64()]
}

// GetRandomJWTSigningKey returns a random JWTSigningKey for JWT signing.
func GetRandomJWTSigningKey() *types.JWTSigningKey {
	jwkMu.RLock()
	defer jwkMu.RUnlock()

	if len(activeSigningKeys) == 0 {
		return nil
	}
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(activeSigningKeys))))
	return activeSigningKeys[i.Int64()]
}

// GetJWTSigningKeyByKeyID returns the JWTSigningKey for the given key ID, or nil if not found.
func GetJWTSigningKeyByKeyID(keyID uuid.UUID) *types.JWTSigningKey {
	jwkMu.RLock()
	defer jwkMu.RUnlock()

	for _, k := range activeSigningKeys {
		if k.KeyID == keyID {
			return k
		}
	}
	return nil
}

// GetJWKSetForDiscovery returns the full JWK set for OIDC discovery (/.well-known/jwks.json).
func GetJWKSetForDiscovery() types.JWKSet {
	jwkMu.RLock()
	defer jwkMu.RUnlock()

	keys := make([]types.JWK, len(activeJWKs))
	copy(keys, activeJWKs)
	return types.JWKSet{Keys: keys}
}

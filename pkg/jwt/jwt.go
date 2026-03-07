// Package jwt contains JWT and JWK related functions as utility package
package jwt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/types"
)

// EncryptPrivateKey encrypts jwt private keys (or as known JWK) for storing in database safely
// MasterKey will be used in decrypting the the key too.
func EncryptPrivateKey(
	privateKey []byte,
	masterKey []byte,
) ([]byte, error) {
	block, _ := aes.NewCipher(masterKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, privateKey, nil), nil
}

// DecryptPrivateKey decrypts jwt private keys (or as known JWK) from storage to usage on JWT
func DecryptPrivateKey(encryptedData []byte, masterKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Split the nonce and the actual ciphertext
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Decrypt and verify
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// GeneratedKey is the result of key generation - use for persistence and in-memory session.
// Model: ent.AccessTokenKey (persisted). General purpose: JWK (served), JWTSigningKey (signing).
type GeneratedKey struct {
	ID                  uuid.UUID
	PrivateKey          *rsa.PrivateKey
	PrivateKeyEncrypted []byte
	PublicPem           string
	State               types.AccessTokenKeyStates
}

// ToSigningKey returns a JWTSigningKey for in-memory JWT signing.
func (g *GeneratedKey) ToSigningKey() *types.JWTSigningKey {
	return &types.JWTSigningKey{
		KeyID:      g.ID,
		PrivateKey: g.PrivateKey,
		Algorithm:  types.JWKAlgorithmRS256,
	}
}

// GenerateKey creates a new RSA key pair for JWT signing.
// Returns GeneratedKey: persist to ent, use ToSigningKey() for signing, convert to JWK for discovery.
func GenerateKey(masterKey []byte) (*GeneratedKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	encPriv, err := EncryptPrivateKey(privBytes, masterKey)
	if err != nil {
		return nil, err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	publicPem := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))

	return &GeneratedKey{
		ID:                  uuid.New(),
		PrivateKey:          priv,
		PrivateKeyEncrypted: encPriv,
		PublicPem:           publicPem,
		State:               types.AccessTokenKeyStatesActive,
	}, nil
}

// ParsePublicPEM parses the string public key into rsa.PublicKey type
func ParsePublicPEM(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

// ConvertToJWK converts public PEM + key ID into RFC 7517 JWK.
// Use with ent.AccessTokenKey (PublicPem, ID) or GeneratedKey (PublicPem, ID).
func ConvertToJWK(publicPem string, keyID uuid.UUID) (*types.JWK, error) {
	pub, err := ParsePublicPEM(publicPem)
	if err != nil {
		return nil, err
	}

	// RFC 7518: N and E must be base64url encoded
	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes())

	return &types.JWK{
		KTY:       types.JWKKeyTypeRSA,
		Use:       types.JWKKeyUseSignature,
		Algorithm: types.JWKAlgorithmRS256,
		KeyID:     keyID,
		N:         n,
		E:         e,
	}, nil
}

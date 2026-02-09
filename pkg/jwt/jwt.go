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

// GenerateKey is a function that creates JWT keys (or as known JWK)
func GenerateKey(
	masterKey []byte,
) (*types.AccessTokenKey, error) {
	accessTokenKey := &types.AccessTokenKey{
		ID:                  uuid.New(),
		PrivateKey:          nil,
		PrivateEncryptedKey: nil,
		PublicPEM:           "",
		State:               types.AccessTokenKeyStatesActive,
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	accessTokenKey.PrivateKey = x509.MarshalPKCS1PrivateKey(priv)
	if encPriv, err := EncryptPrivateKey(accessTokenKey.PrivateKey, masterKey); err != nil {
		return nil, err
	} else {
		accessTokenKey.PrivateEncryptedKey = encPriv
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	accessTokenKey.PublicPEM = string(
		pem.EncodeToMemory(
			&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes},
		),
	)

	return accessTokenKey, nil
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

// ConvertToJWK converts AccessTokenKey object into JWK type
func ConvertToJWK(
	accessTokenKey *types.AccessTokenKey,
) (*types.JWK, error) {
	publicPem, err := ParsePublicPEM(accessTokenKey.PublicPEM)
	if err != nil {
		return nil, err
	}

	return &types.JWK{
		KTY:       types.JWKTypesRSA,
		Use:       types.JWKUseTypesSignature,
		Algorithm: "RS256",
		KeyID:     accessTokenKey.ID,
		N:         base64.RawURLEncoding.EncodeToString(publicPem.N.Bytes()),
		E: base64.RawStdEncoding.EncodeToString(
			big.NewInt(int64(publicPem.E)).Bytes(),
		),
	}, nil
}

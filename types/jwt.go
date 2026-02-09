package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWKTypes string

const JWKTypesRSA JWKTypes = "RSA"

type JWKUseTypes string

const JWKUseTypesSignature JWKUseTypes = "sig"

type JWKAlgorithms string

const JWKAlgorithmsRSA256 JWKAlgorithms = "RSA256"

func (alg JWKAlgorithms) ToJWTSigningMethod() jwt.SigningMethod {
	return jwt.SigningMethodRS256
}

type JWK struct {
	KTY       JWKTypes      `json:"kty"`
	Use       JWKUseTypes   `json:"use"`
	Algorithm JWKAlgorithms `json:"alg"`
	KeyID     uuid.UUID     `json:"kid"`
	N         string        `json:"n"`
	E         string        `json:"e"`
}

type AccessTokenKeyStates int8

const (
	AccessTokenKeyStatesActive   AccessTokenKeyStates = 0
	AccessTokenKeyStatesPrevious AccessTokenKeyStates = 1
	AccessTokenKeyStatesRedired  AccessTokenKeyStates = 2
)

func (AccessTokenKeyStates) Values() []int8 {
	return []int8{
		int8(AccessTokenKeyStatesActive),
		int8(AccessTokenKeyStatesPrevious),
		int8(AccessTokenKeyStatesRedired),
	}
}

type AccessTokenKey struct {
	ID                  uuid.UUID
	PrivateKey          []byte
	PrivateEncryptedKey []byte
	PublicPEM           string
	State               AccessTokenKeyStates
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           *time.Time
}

type AccessTokenHeader struct {
	KeyID     uuid.UUID     `json:"kid"`
	Algorithm JWKAlgorithms `json:"alg"`
	Type      string        `json:"typ"` // Always JWT
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserID   uuid.UUID `json:"sub"`
	UserRole UserRoles `json:"role"`
}

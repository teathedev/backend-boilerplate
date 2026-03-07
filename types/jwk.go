package types

import (
	"crypto/rsa"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWK types - RFC 7517 / OIDC compliant (general purpose)

// JWKKeyType (kty) - RFC 7517 Section 4.1
type JWKKeyType string

const (
	JWKKeyTypeRSA JWKKeyType = "RSA"
	JWKKeyTypeEC  JWKKeyType = "EC"
)

// JWKKeyUse (use) - RFC 7517 Section 4.2
type JWKKeyUse string

const (
	JWKKeyUseSignature  JWKKeyUse = "sig"
	JWKKeyUseEncryption JWKKeyUse = "enc"
)

// JWKAlgorithm (alg) - RFC 7518 JWA, use "RS256" not "RSA256"
type JWKAlgorithm string

const (
	JWKAlgorithmRS256 JWKAlgorithm = "RS256"
	JWKAlgorithmES256 JWKAlgorithm = "ES256"
)

// ToJWTSigningMethod returns the jwt-go SigningMethod for this algorithm.
func (a JWKAlgorithm) ToJWTSigningMethod() jwt.SigningMethod {
	switch a {
	case JWKAlgorithmRS256:
		return jwt.SigningMethodRS256
	case JWKAlgorithmES256:
		return jwt.SigningMethodES256
	default:
		return jwt.SigningMethodRS256
	}
}

// JWK is the public key representation for OIDC discovery and verification.
// RFC 7517 compliant - safe to serialize as JSON and serve at /.well-known/jwks.json
type JWK struct {
	KTY       JWKKeyType   `json:"kty"`
	Use       JWKKeyUse    `json:"use"`
	Algorithm JWKAlgorithm `json:"alg"`
	KeyID     uuid.UUID    `json:"kid"`
	N         string       `json:"n,omitempty"` // RSA modulus (base64url)
	E         string       `json:"e,omitempty"` // RSA exponent (base64url)
}

// JWKSet is the OIDC discovery format - {"keys": [...]}
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

// JWTSigningKey holds the private key for in-memory JWT signing.
// Used internally - never serialized or exposed.
type JWTSigningKey struct {
	KeyID      uuid.UUID
	PrivateKey *rsa.PrivateKey
	Algorithm  JWKAlgorithm
}

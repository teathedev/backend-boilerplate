package types

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// AccessTokenKeyStates - lifecycle state for persisted keys (ent model)
type AccessTokenKeyStates int8

const (
	AccessTokenKeyStatesActive   AccessTokenKeyStates = 0
	AccessTokenKeyStatesPrevious AccessTokenKeyStates = 1
	AccessTokenKeyStatesRetired  AccessTokenKeyStates = 2
)

func (AccessTokenKeyStates) Values() []int8 {
	return []int8{
		int8(AccessTokenKeyStatesActive),
		int8(AccessTokenKeyStatesPrevious),
		int8(AccessTokenKeyStatesRetired),
	}
}

// AccessTokenHeader - JWT header fields
type AccessTokenHeader struct {
	KeyID     uuid.UUID    `json:"kid"`
	Algorithm JWKAlgorithm `json:"alg"`
	Type      string       `json:"typ"` // Always "JWT"
}

// AccessTokenClaims - JWT payload
type AccessTokenClaims struct {
	*jwt.RegisteredClaims
	UserID   uuid.UUID `json:"sub"`
	UserRole UserRoles `json:"role"`
}

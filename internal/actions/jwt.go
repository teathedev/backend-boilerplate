package actions

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/types"
)

func CreateAccessToken(
	claims *types.AccessTokenClaims,
) (string, error) {
	if claims.RegisteredClaims == nil {
		claims.RegisteredClaims = &jwt.RegisteredClaims{
			Issuer:    "default-iss",
			Subject:   claims.UserID.String(),
			Audience:  jwt.ClaimStrings{},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		}
	}

	token := jwt.NewWithClaims(
		types.JWKAlgorithmRS256.ToJWTSigningMethod(),
		claims,
	)

	token.Header["typ"] = "JWT"

	signingKey := GetRandomJWTSigningKey()
	if signingKey == nil {
		return "", fmt.Errorf("no signing key available")
	}

	token.Header["kid"] = signingKey.KeyID.String()

	signedToken, err := token.SignedString(signingKey.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signedToken, nil
}

// DecodeAccessToken parses the JWT without validating the signature.
// Use when you need to inspect claims or kid before verification (e.g. to choose the key).
// Call VerifyAccessToken on the result before trusting the claims.
func DecodeAccessToken(tokenString string) (*jwt.Token, error) {
	claims := &types.AccessTokenClaims{}
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, fmt.Errorf("decode access token: %w", err)
	}
	return token, nil
}

// VerifyAccessToken verifies the signature of an already-decoded token.
// The token must have been produced by DecodeAccessToken. Returns an error if the signature is invalid.
func VerifyAccessToken(decoded *jwt.Token) error {
	if decoded == nil || decoded.Raw == "" {
		return fmt.Errorf("decoded token is nil or empty")
	}

	kidVal, ok := decoded.Header["kid"]
	if !ok {
		return fmt.Errorf("token header missing kid")
	}
	kidStr, ok := kidVal.(string)
	if !ok {
		return fmt.Errorf("token header kid is not a string")
	}
	keyID, err := uuid.Parse(kidStr)
	if err != nil {
		return fmt.Errorf("invalid token kid: %w", err)
	}

	signingKey := GetJWTSigningKeyByKeyID(keyID)
	if signingKey == nil {
		return fmt.Errorf("no signing key found for kid %s", keyID)
	}

	signingString, err := decoded.SigningString()
	if err != nil {
		return fmt.Errorf("signing string: %w", err)
	}

	if err := decoded.Method.Verify(signingString, decoded.Signature, signingKey.PrivateKey.PublicKey); err != nil {
		return fmt.Errorf("token signature invalid: %w", err)
	}
	return nil
}

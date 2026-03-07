package actions

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

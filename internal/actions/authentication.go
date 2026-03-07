package actions

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/teathedev/backend-boilerplate/internal/ent"
)

func HashPassword(
	password string,
	salt string,
) string {
	hash := sha512.Sum512([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

func CheckUserPassword(
	user *ent.User,
	password string,
) bool {
	hashed := HashPassword(password, user.PasswordSalt)
	return user.PasswordHash == hashed
}

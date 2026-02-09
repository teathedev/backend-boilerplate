package filters

import (
	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/internal/ent/predicate"
	"github.com/teathedev/backend-boilerplate/internal/ent/refreshtoken"
)

func UsersRefreshToken(userID uuid.UUID) predicate.RefreshToken {
	return refreshtoken.UserIDEQ(userID)
}

func ActiveRefreshToken() predicate.RefreshToken {
	return refreshtoken.IsClaimed(false)
}

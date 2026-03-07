package actions

import (
	"context"

	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/internal/db"
	"github.com/teathedev/backend-boilerplate/internal/ent"
)

func CreateRefreshToken(
	ctx context.Context,
	user *ent.User,
) (*ent.RefreshToken, error) {
	return db.Client.RefreshToken.Create().
		SetUserID(user.ID).
		SetIsClaimed(false).
		SetToken(uuid.New().String()).
		Save(ctx)
}

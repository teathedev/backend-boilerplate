package actions

import (
	"github.com/teathedev/backend-boilerplate/internal/ent"
	"github.com/teathedev/backend-boilerplate/types"
)

// EntUserToTypesUser maps ent.User (model) to types.User (domain).
func EntUserToTypesUser(entUser *ent.User) *types.User {
	if entUser == nil {
		return nil
	}

	return &types.User{
		ID:          entUser.ID,
		PhoneNumber: entUser.PhoneNumber,
		Email:       entUser.Email,
		Username:    entUser.Username,
		Role:        entUser.Role,
		State:       entUser.State,
		FirstName:   entUser.FirstName,
		LastName:    entUser.LastName,
		CreatedAt:   entUser.CreatedAt,
		UpdatedAt:   entUser.UpdatedAt,
	}
}

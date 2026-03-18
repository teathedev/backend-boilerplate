package usecases

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/constants"
	"github.com/teathedev/backend-boilerplate/internal/actions"
	"github.com/teathedev/backend-boilerplate/internal/db"
	"github.com/teathedev/backend-boilerplate/internal/ent/user"
	"github.com/teathedev/backend-boilerplate/internal/filters"
	"github.com/teathedev/backend-boilerplate/types"
	"github.com/teathedev/pkg/errors"
	"github.com/teathedev/pkg/logger"
	"github.com/teathedev/pkg/validation"
)

type userUseCase struct {
	log logger.Logger
}

var User userUseCase

func init() {
	User = userUseCase{
		log: logger.New("User"),
	}
}

func (uc *userUseCase) UpdateMe(ctx context.Context, userID uuid.UUID, req *types.UpdateMe) (*types.User, error) {
	if err := validation.ValidateStruct(req); err != nil {
		return nil, err
	}

	updater := db.Client.User.UpdateOneID(userID)
	if req.FirstName != nil {
		updater.SetFirstName(*req.FirstName)
	}
	if req.LastName != nil {
		updater.SetLastName(*req.LastName)
	}
	if req.Email != nil {
		email := strings.ToLower(*req.Email)
		existing, err := db.Client.User.Query().Where(user.EmailEQ(email)).First(ctx)
		if err == nil && existing.ID != userID {
			return nil, errors.NewBadInput("User.UpdateMe", []errors.BadInputField{
				{Field: "email", Condition: "not_valid", Value: *req.Email},
			})
		}
		updater.SetEmail(email)
	}
	if req.Username != nil {
		username := strings.ToLower(*req.Username)
		existing, err := db.Client.User.Query().Where(user.UsernameEQ(username)).First(ctx)
		if err == nil && existing.ID != userID {
			return nil, errors.NewBadInput("User.UpdateMe", []errors.BadInputField{
				{Field: "username", Condition: "not_valid", Value: *req.Username},
			})
		}
		updater.SetUsername(username)
	}
	if req.PhoneNumber != nil {
		existing, err := db.Client.User.Query().Where(user.PhoneNumberEQ(*req.PhoneNumber)).First(ctx)
		if err == nil && existing.ID != userID {
			return nil, errors.NewBadInput("User.UpdateMe", []errors.BadInputField{
				{Field: "phoneNumber", Condition: "not_valid", Value: *req.PhoneNumber},
			})
		}
		updater.SetPhoneNumber(*req.PhoneNumber)
	}

	entUser, err := updater.Save(ctx)
	if err != nil {
		uc.log.Error("Failed to update user", logger.LogParams{
			"RequestID": ctx.Value(constants.RequestID),
			"UserID":    userID,
			"Error":     err,
		})
		return nil, errors.New("User.UpdateMe", "failed to update user")
	}
	return actions.EntUserToTypesUser(entUser), nil
}

func (uc *userUseCase) UpdatePassword(ctx context.Context, userID uuid.UUID, req *types.UpdatePassword) error {
	if req == nil {
		return errors.NewBadInput("User.UpdatePassword", []errors.BadInputField{
			{Field: "body", Condition: "not_valid", Value: ""},
		})
	}
	if err := validation.ValidateStruct(req); err != nil {
		return err
	}

	entUser, err := db.Client.User.
		Query().
		Where(
			user.And(
				user.IDEQ(userID),
				filters.ActiveUsers(),
			),
		).
		First(ctx)
	if err != nil {
		return errors.New("User.UpdatePassword", "user not found")
	}

	if !actions.CheckUserPassword(entUser, req.CurrentPassword) {
		return errors.NewBadInput("User.UpdatePassword", []errors.BadInputField{
			{Field: "currentPassword", Condition: "not_valid", Value: ""},
		})
	}

	salt := uuid.New().String()
	hash := actions.HashPassword(req.NewPassword, salt)
	_, err = db.Client.User.UpdateOneID(userID).
		SetPasswordSalt(salt).
		SetPasswordHash(hash).
		Save(ctx)
	if err != nil {
		uc.log.Error("Failed to update password", logger.LogParams{
			"RequestID": ctx.Value(constants.RequestID),
			"UserID":    userID,
			"Error":     err,
		})
		return errors.New("User.UpdatePassword", "failed to update password")
	}
	return nil
}

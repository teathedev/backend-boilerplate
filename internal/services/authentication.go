package services

import (
	"context"

	"github.com/teathedev/backend-boilerplate/constants"
	"github.com/teathedev/backend-boilerplate/internal/actions"
	"github.com/teathedev/backend-boilerplate/internal/db"
	"github.com/teathedev/backend-boilerplate/internal/ent/user"
	"github.com/teathedev/backend-boilerplate/internal/filters"
	"github.com/teathedev/backend-boilerplate/pkg/errors"
	"github.com/teathedev/backend-boilerplate/pkg/logger"
	"github.com/teathedev/backend-boilerplate/types"
)

type authenticationService struct {
	log logger.Logger
}

var AuthenticationService authenticationService

func init() {
	AuthenticationService = authenticationService{
		log: logger.New("AuthenticationService"),
	}
}

func (svc *authenticationService) Login(
	ctx context.Context,
	request *types.Login,
) (*types.AuthenticationResult, error) {
	userItem, err := db.Client.User.
		Query().
		Where(
			user.And(
				filters.UserByIdentifier(request.Identifier),
				filters.ActiveUsers(),
			),
		).
		First(ctx)
	if err != nil {
		if errors.IsNotFound(err) {
			svc.log.Trace(
				"User not found by identifier!",
				logger.LogParams{
					"Identifier": request.Identifier,
					"RequestID":  ctx.Value(constants.RequestID),
					"Error":      err,
				},
			)
		} else {
			svc.log.Error(
				"Failed to find user by identifier!",
				logger.LogParams{
					"Identifier": request.Identifier,
					"RequestID":  ctx.Value(constants.RequestID),
					"Error":      err,
				},
			)

			return nil, errors.New(
				"AuthenticationService.Login",
				"failed to find user by identifier",
			)
		}

		return nil, errors.NewBadInput("AuthenticationService", []errors.BadInputField{
			{
				Field:     "identifier",
				Condition: "not_valid",
				Value:     request.Identifier,
			},
			{
				Field:     "password",
				Condition: "not_valid",
				Value:     request.Identifier,
			},
		})
	}

	if userItem == nil {
		svc.log.Trace(
			"User not found by identifier!",
			logger.LogParams{
				"Identifier": request.Identifier,
				"RequestID":  ctx.Value(constants.RequestID),
				"Error":      err,
			},
		)
		return nil, errors.NewBadInput(
			"AuthenticationService.Login",
			[]errors.BadInputField{
				{
					Field:     "identifier",
					Condition: "not_valid",
					Value:     request.Identifier,
				},
			},
		)
	}

	if !actions.CheckUserPassword(
		userItem,
		request.Password,
	) {
		svc.log.Trace(
			"Failed to match password!",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Password":  request.Password,
			},
		)
		return nil, errors.NewBadInput(
			"AuthenticationService",
			[]errors.BadInputField{
				{
					Field:     "identifier",
					Condition: "not_valid",
					Value:     request.Identifier,
				},
				{
					Field:     "password",
					Condition: "not_valid",
					Value:     request.Identifier,
				},
			},
		)
	}

	accessToken, err := actions.CreateAccessToken(&types.AccessTokenClaims{
		UserID:   userItem.ID,
		UserRole: userItem.Role,
	})
	if err != nil {
		svc.log.Error(
			"Failed to create access token",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Error":     err,
			},
		)
		return nil, errors.New("AuthenticationService.Login", "failed to create access token")
	}

	return &types.AuthenticationResult{
		AccessToken: accessToken,
		User:        actions.EntUserToTypesUser(userItem),
	}, nil
}

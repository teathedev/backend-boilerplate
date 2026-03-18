package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/teathedev/backend-boilerplate/constants"
	"github.com/teathedev/backend-boilerplate/internal/actions"
	"github.com/teathedev/backend-boilerplate/internal/db"
	"github.com/teathedev/backend-boilerplate/internal/ent"
	"github.com/teathedev/backend-boilerplate/internal/ent/refreshtoken"
	"github.com/teathedev/backend-boilerplate/internal/ent/user"
	"github.com/teathedev/backend-boilerplate/internal/filters"
	"github.com/teathedev/backend-boilerplate/types"
	"github.com/teathedev/pkg/errors"
	"github.com/teathedev/pkg/logger"
	"github.com/teathedev/pkg/validation"
)

type authenticationUseCase struct {
	log logger.Logger
}

var Authentication authenticationUseCase

func init() {
	Authentication = authenticationUseCase{
		log: logger.New("Authentication"),
	}
}

func (uc *authenticationUseCase) Login(
	ctx context.Context,
	request *types.Login,
) (*types.AuthenticationResult, error) {
	if request == nil {
		return nil, errors.NewBadInput(
			"Authentication",
			[]errors.BadInputField{
				{
					Field:     "body",
					Condition: "not_valid",
					Value:     "",
				},
			})
	}

	if err := validation.ValidateStruct(request); err != nil {
		uc.log.Error(
			"Failed to validate request",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Error":     err,
			},
		)
		return nil, err
	}

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
		if ent.IsNotFound(err) {
			uc.log.Trace(
				"User not found by identifier!",
				logger.LogParams{
					"Identifier": request.Identifier,
					"RequestID":  ctx.Value(constants.RequestID),
					"Error":      err,
				},
			)
		} else {
			uc.log.Error(
				"Failed to find user by identifier!",
				logger.LogParams{
					"Identifier": request.Identifier,
					"RequestID":  ctx.Value(constants.RequestID),
					"Error":      err,
				},
			)

			return nil, errors.New(
				"Authentication.Login",
				"failed to find user by identifier",
			)
		}

		return nil, errors.NewBadInput("Authentication", []errors.BadInputField{
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
		uc.log.Trace(
			"User not found by identifier!",
			logger.LogParams{
				"Identifier": request.Identifier,
				"RequestID":  ctx.Value(constants.RequestID),
				"Error":      err,
			},
		)
		return nil, errors.NewBadInput(
			"Authentication.Login",
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
		uc.log.Trace(
			"Failed to match password!",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Password":  request.Password,
			},
		)
		return nil, errors.NewBadInput(
			"Authentication",
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

	accessToken, err := actions.CreateAccessToken(
		&types.AccessTokenClaims{
			UserID:   userItem.ID,
			UserRole: userItem.Role,
		},
	)
	if err != nil {
		uc.log.Error(
			"Failed to create access token",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Error":     err,
			},
		)
		return nil, errors.New("Authentication.Login", "failed to create access token")
	}

	refreshToken, err := actions.CreateRefreshToken(ctx, userItem)
	if err != nil {
		uc.log.Error(
			"Failed to create refresh token",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Error":     err,
			},
		)
		return nil, errors.New("Authentication.Login", "failed to create refresh token")

	}

	u := actions.EntUserToTypesUser(userItem)
	return &types.AuthenticationResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		User:         *u,
	}, nil
}

func (uc *authenticationUseCase) Register(
	ctx context.Context,
	request *types.Register,
) (*types.AuthenticationResult, error) {
	if request == nil {
		return nil, errors.NewBadInput(
			"Authentication",
			[]errors.BadInputField{
				{
					Field:     "body",
					Condition: "not_valid",
					Value:     "",
				},
			})
	}

	if err := validation.ValidateStruct(request); err != nil {
		return nil, errors.NewBadInput(
			"Authentication.Register",
			[]errors.BadInputField{
				{
					Field:     "body",
					Condition: errors.BadInputConditionNotValid,
					Value:     err.Error(),
				},
			},
		)
	}

	existingUser, err := db.Client.User.
		Query().
		Where(
			user.Or(
				user.PhoneNumberEQ(request.PhoneNumber),
				user.EmailEQ(request.Email),
				user.UsernameEQ(request.Username),
			),
		).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			uc.log.Trace(
				"User not found by identifier!",
				logger.LogParams{
					"PhoneNumber": request.PhoneNumber,
					"Email":       request.Email,
					"Username":    request.Username,
					"RequestID":   ctx.Value(constants.RequestID),
					"Error":       err,
				},
			)
		} else {
			uc.log.Error(
				"Failed to find user by identifier!",
				logger.LogParams{
					"PhoneNumber": request.PhoneNumber,
					"Email":       request.Email,
					"Username":    request.Username,
					"RequestID":   ctx.Value(constants.RequestID),
					"Error":       err,
				},
			)

			return nil, errors.New(
				"Authentication.Register",
				"failed to find user by identifier",
			)
		}
	}

	if existingUser != nil {
		uc.log.Trace(
			"User already exists!",
			logger.LogParams{
				"PhoneNumber": request.PhoneNumber,
				"Email":       request.Email,
				"Username":    request.Username,
				"RequestID":   ctx.Value(constants.RequestID),
				"Error":       err,
			},
		)
		return nil, errors.NewBadInput(
			"Authentication.Register",
			[]errors.BadInputField{
				{
					Field:     "phoneNumber",
					Condition: "not_valid",
					Value:     request.PhoneNumber,
				},
				{
					Field:     "email",
					Condition: "not_valid",
					Value:     request.Email,
				},
				{
					Field:     "username",
					Condition: "not_valid",
					Value:     request.Username,
				},
			},
		)
	}

	passwordSalt := uuid.New().String()
	passwordHash := actions.HashPassword(request.Password, passwordSalt)

	userItem, err := db.Client.User.
		Create().
		SetPhoneNumber(request.PhoneNumber).
		SetEmail(request.Email).
		SetUsername(request.Username).
		SetRole(request.Role).
		SetFirstName(request.FirstName).
		SetLastName(request.LastName).
		SetPasswordSalt(passwordSalt).
		SetPasswordHash(passwordHash).
		SetState(types.UserStatesActive).
		Save(ctx)
	if err != nil {
		uc.log.Error(
			"Failed to create user",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Error":     err,
			},
		)
		return nil, errors.New("Authentication.Register", "failed to create user")
	}

	accessToken, err := actions.CreateAccessToken(&types.AccessTokenClaims{
		UserID:   userItem.ID,
		UserRole: userItem.Role,
	})
	if err != nil {
		return nil, errors.New("Authentication.Register", "failed to create access token")
	}

	refreshToken, err := actions.CreateRefreshToken(ctx, userItem)
	if err != nil {
		return nil, errors.New("Authentication.Register", "failed to create refresh token")
	}

	u := actions.EntUserToTypesUser(userItem)
	return &types.AuthenticationResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		User:         *u,
	}, nil
}

func (uc *authenticationUseCase) GetUserByToken(
	ctx context.Context,
	token string,
) (*types.User, error) {
	decoded, err := actions.DecodeAccessToken(token)
	if err != nil {
		return nil, errors.New("Authentication.GetUserByToken", "failed to decode access token")
	}
	if err := actions.VerifyAccessToken(decoded); err != nil {
		return nil, errors.New("Authentication.GetUserByToken", "failed to verify access token")
	}
	tokenClaims, ok := decoded.Claims.(*types.AccessTokenClaims)
	if !ok {
		return nil, errors.New("Authentication.GetUserByToken", "invalid access token claims")
	}

	userItem, err := db.Client.User.
		Query().
		Where(
			user.And(
				filters.ActiveUsers(),
				user.IDEQ(tokenClaims.UserID),
			),
		).
		First(ctx)
	if err != nil {
		return nil, errors.New("Authentication.GetUserByToken", "failed to find user by token")
	}

	return &types.User{
		ID:          userItem.ID,
		PhoneNumber: userItem.PhoneNumber,
		Email:       userItem.Email,
		Username:    userItem.Username,
		Role:        userItem.Role,
		State:       userItem.State,
		FirstName:   userItem.FirstName,
		LastName:    userItem.LastName,
	}, nil
}

func (uc *authenticationUseCase) Refresh(
	ctx context.Context,
	request *types.RefreshTokenRequest,
) (*types.AuthenticationResult, error) {
	if request == nil {
		return nil, errors.NewBadInput(
			"Authentication.Refresh",
			[]errors.BadInputField{
				{Field: "body", Condition: "not_valid", Value: ""},
			})
	}
	if err := validation.ValidateStruct(request); err != nil {
		return nil, err
	}

	rt, err := db.Client.RefreshToken.
		Query().
		Where(
			refreshtoken.And(
				refreshtoken.TokenEQ(request.RefreshToken),
				filters.ActiveRefreshToken(),
			),
		).
		WithUser().
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			uc.log.Trace(
				"Refresh token not found or already claimed",
				logger.LogParams{
					"RequestID": ctx.Value(constants.RequestID),
				},
			)
		}
		return nil, errors.NewBadInput(
			"Authentication.Refresh",
			[]errors.BadInputField{
				{
					Field:     "refreshToken",
					Condition: "not_valid",
					Value:     "",
				},
			})
	}

	userItem, err := rt.Edges.UserOrErr()
	if err != nil || userItem == nil {
		return nil, errors.New("Authentication.Refresh", "failed to load user for refresh token")
	}

	if userItem.State != types.UserStatesActive && userItem.State != types.UserStatesInvited {
		return nil, errors.NewBadInput(
			"Authentication.Refresh",
			[]errors.BadInputField{
				{
					Field:     "refreshToken",
					Condition: "not_valid",
					Value:     "",
				},
			})
	}

	// Mark old token as claimed (one-time use)
	_, err = db.Client.RefreshToken.UpdateOne(rt).SetIsClaimed(true).Save(ctx)
	if err != nil {
		uc.log.Error(
			"Failed to claim refresh token",
			logger.LogParams{
				"RequestID": ctx.Value(constants.RequestID),
				"Error":     err,
			},
		)
		return nil, errors.New("Authentication.Refresh", "failed to process refresh token")
	}

	accessToken, err := actions.CreateAccessToken(&types.AccessTokenClaims{
		UserID:   userItem.ID,
		UserRole: userItem.Role,
	})
	if err != nil {
		return nil, errors.New("Authentication.Refresh", "failed to create access token")
	}

	newRefreshToken, err := actions.CreateRefreshToken(ctx, userItem)
	if err != nil {
		return nil, errors.New("Authentication.Refresh", "failed to create refresh token")
	}

	u := actions.EntUserToTypesUser(userItem)
	return &types.AuthenticationResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
		User:         *u,
	}, nil
}

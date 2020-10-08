package resolver

import (
	"context"
	"time"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/entities"
	internalErrors "github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/graph/model"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	internalTime "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	pkgErrors "github.com/pkg/errors"
)

func (r *mutationResolver) createUser(ctx context.Context, input model.CreateUserInput) (*model.AuthOutput, error) {
	validationResult := r.validator.ValidateCreateUserInput(input, r.languageTagFromContext(ctx))
	if validationResult.HasError {
		r.addValidationErrors(ctx, validationResult)
		return nil, internalErrors.ErrValidationError
	}

	hashedPassword, err := r.passwordService.HashPassword(input.Password)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "could not hash password"))
		return nil, internalErrors.ErrInternalServerError
	}

	user := entities.User{
		ID:        id.New(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = r.db.UserRepository().Store(user)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "cannot save user in the database"))
		return nil, internalErrors.ErrInternalServerError
	}

	token, err := r.jwtService.GenerateTokenForUserID(user.ID)
	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrapf(err, "cannot generate jwt token for user with ID: %s", user.ID.String()))
		return nil, internalErrors.ErrInternalServerError
	}

	return &model.AuthOutput{
		User: &model.User{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(internalTime.DefaultFormat),
			UpdatedAt: user.UpdatedAt.Format(internalTime.DefaultFormat),
		},
		Token: &model.Token{
			Value: token,
		},
	}, nil
}

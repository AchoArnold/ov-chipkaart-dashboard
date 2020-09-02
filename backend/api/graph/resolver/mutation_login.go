package resolver

import (
	"context"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"
	internalErrors "github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/model"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/validator"
	internalTime "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	pkgErrors "github.com/pkg/errors"
)

func (r *mutationResolver) login(ctx context.Context, input model.LoginInput) (*model.AuthOutput, error) {
	validationResult := r.validator.ValidateLoginInput(input, r.languageTagFromContext(ctx))
	if validationResult.HasError {
		r.addValidationErrors(ctx, validationResult)
		return nil, internalErrors.ErrValidationError
	}

	user, err := r.db.UserRepository().FindByEmail(input.Email)
	if err == database.ErrEntityNotFound {
		r.addError(ctx, fieldEmail, validator.ErrInvalidEmailOrPassword.Error(), CodeValidationError)
		r.addError(ctx, fieldPassword, validator.ErrInvalidEmailOrPassword.Error(), CodeValidationError)
		return nil, internalErrors.ErrValidationError
	}

	if err != nil {
		r.errorHandler.CaptureError(ctx, pkgErrors.Wrap(err, "cannot find user by email"))
		return nil, internalErrors.ErrInternalServerError
	}

	passwordIsValid := r.passwordService.CheckPasswordHash(input.Password, user.Password)
	if !passwordIsValid {
		r.addError(ctx, fieldEmail, validator.ErrInvalidEmailOrPassword.Error(), CodeValidationError)
		r.addError(ctx, fieldPassword, validator.ErrInvalidEmailOrPassword.Error(), CodeValidationError)
		return nil, internalErrors.ErrValidationError
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

package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/entities"
	internalErrors "github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/generated"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/model"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/validator"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"
	internalTime "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	pkgErrors "github.com/pkg/errors"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.AuthOutput, error) {
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthOutput, error) {
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

func (r *mutationResolver) CancelToken(ctx context.Context, input model.CancelTokenInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) StoreRequest(ctx context.Context, input *model.StoreAnalyzeRequestInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context) (*model.User, error) {
	return &model.User{}, nil
}

func (r *queryResolver) AnalyzeRequests(ctx context.Context, skip *int, take *int, orderBy *string) ([]*model.AnalzyeRequestDetails, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
var (
	fieldPassword = "password"
	fieldEmail    = "email"
)

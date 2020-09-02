package resolver

import (
	"context"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/id"

	"github.com/99designs/gqlgen/graphql"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/validator"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/middlewares"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/services/jwt"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/services/password"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errorhandler"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/logger"
	"github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/text/language"
)

const (
	// CodeValidationError is the code that is returned on a validation error message
	CodeValidationError = "VALIDATION_ERROR"
)

var (
	// ErrUnauthorizedRequest is the error when the user is unauthorized
	ErrUnauthorizedRequest = errors.New("401 Unauthorized")
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver resolves
type Resolver struct {
	db              database.DB
	validator       validator.Validator
	passwordService password.Service
	errorHandler    errorhandler.ErrorHandler
	logger          logger.Logger
	jwtService      jwt.Service
}

// NewResolver creates a new instance of the resolver
func NewResolver(
	db database.DB,
	validator validator.Validator,
	passwordService password.Service,
	errorHandler errorhandler.ErrorHandler,
	logger logger.Logger,
	jwtService jwt.Service,
) *Resolver {

	return &Resolver{
		db:              db,
		validator:       validator,
		passwordService: passwordService,
		errorHandler:    errorHandler,
		logger:          logger,
		jwtService:      jwtService,
	}
}

func (r *Resolver) languageTagFromContext(ctx context.Context) language.Tag {
	tag, ok := ctx.Value(middlewares.ContextKeyLanguageTag).(*language.Tag)
	if tag == nil || !ok {
		r.errorHandler.CaptureError(ctx, errors.New("cannot fetch language tag from resolver"))
		return language.English
	}

	return *tag
}

func (r *Resolver) userIDFromContext(ctx context.Context) (userID id.ID, err error) {
	userID, err = id.FromInterface(ctx.Value(middlewares.ContextKeyUserID))

	r.errorHandler.CaptureError(ctx, errors.Wrap(err, "cannot fetch user id from context"))

	return userID, err
}

func (r *Resolver) addValidationErrors(ctx context.Context, result validator.ValidationResult) {
	for field, fieldErrors := range result.Errors {
		for _, err := range fieldErrors {
			r.addError(ctx, field, err, CodeValidationError)
		}
	}
}

func (r *Resolver) addError(ctx context.Context, pathName string, err string, code string) {
	graphql.AddError(ctx, &gqlerror.Error{
		Message: err,
		Path:    append(graphql.GetFieldContext(ctx).Path(), ast.PathName(pathName)),
		Extensions: map[string]interface{}{
			"code": code,
		},
	})
}

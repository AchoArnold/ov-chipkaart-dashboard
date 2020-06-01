package validator

import (
	"net/url"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/model"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

var (
	//ErrInvalidEmailOrPassword is thrown when the user's email/password is wrong.
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
)

// ValidationResult stores the result of a validation
type ValidationResult struct {
	HasError bool
	Errors   url.Values
}

// Validator represents a validator
type Validator interface {
	ValidateCreateUserInput(input model.CreateUserInput, localeTag language.Tag) ValidationResult
	ValidateLoginInput(input model.LoginInput, localeTag language.Tag) ValidationResult
}

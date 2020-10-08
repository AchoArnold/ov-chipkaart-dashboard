package validator

import (
	"net/url"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api-service/graph/model"
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
	ValidateStoreAnalzyeRequest(input model.StoreAnalyzeRequestInput, localTag language.Tag) ValidationResult
	ValidateAnalzyeRequestsInput(skip *int, take *int, orderBy *string, orderDirection *string, localTag language.Tag) ValidationResult
}

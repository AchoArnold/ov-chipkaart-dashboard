package govalidator

import (
	"context"
	"fmt"
	"net/url"

	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/database"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/model"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/api/graph/validator"
	"github.com/NdoleStudio/ov-chipkaart-dashboard/backend/shared/errorhandler"
	"github.com/pkg/errors"
	"github.com/thedevsaddam/govalidator"
	"golang.org/x/text/language"
)

const (
	ruleUserEmailIsUnique = "user_email_is_unique"
)

// GoValidator is a validator using the govalidator package
type GoValidator struct {
	db           database.DB
	errorHandler errorhandler.ErrorHandler
}

// New creates a new go validator
func New(db database.DB, errorHandler errorhandler.ErrorHandler) validator.Validator {
	service := &GoValidator{db, errorHandler}
	service.init()

	return service
}

// ValidateCreateUserInput validates the create user input request
func (service GoValidator) ValidateCreateUserInput(input model.CreateUserInput, _ language.Tag) (result validator.ValidationResult) {
	v := govalidator.New(govalidator.Options{
		Data: &input,
		Rules: govalidator.MapData{
			"firstName": []string{"required", "min:1", "max:50"},
			"lastName":  []string{"required", "min:1", "max:50"},
			"email":     []string{"required", "email", ruleUserEmailIsUnique},
			"password":  []string{"required"},
			"reCaptcha": []string{"required"},
		},
	})

	return service.urlValuesToResult(v.ValidateStruct())
}

// ValidateLoginInput validates the login input object
func (service GoValidator) ValidateLoginInput(input model.LoginInput, _ language.Tag) (result validator.ValidationResult) {
	v := govalidator.New(govalidator.Options{
		Data: &input,
		Rules: govalidator.MapData{
			"email":      []string{"required", "email"},
			"password":   []string{"required"},
			"rememberMe": []string{"bool"},
			"reCaptcha":  []string{"required"},
		},
	})

	return service.urlValuesToResult(v.ValidateStruct())
}

func (service GoValidator) urlValuesToResult(value url.Values) validator.ValidationResult {
	return validator.ValidationResult{
		HasError: len(value) > 0,
		Errors:   value,
	}
}

func (service GoValidator) init() {
	govalidator.AddCustomRule(ruleUserEmailIsUnique, service.ruleUserEmailExists())
}

func (service GoValidator) ruleUserEmailExists() func(field string, rule string, message string, value interface{}) error {
	return func(field string, rule string, message string, value interface{}) error {
		email := value.(string)

		_, err := service.db.UserRepository().FindByEmail(email)
		if err == database.ErrEntityNotFound {
			return nil
		}

		if err != nil {
			service.errorHandler.CaptureError(context.Background(), errors.Wrapf(err, "cannot fetch user by email"))
		}

		if message != "" {
			return errors.New(message)
		}

		return fmt.Errorf("A user already exist with the %s '%s'", field, value)
	}
}

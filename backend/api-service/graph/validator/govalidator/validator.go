package govalidator

import (
	"context"
	"fmt"
	"net/url"
	internalTime "time"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/ovchipkaart"

	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/time"
	internalErrors "github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errors"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/database"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/model"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/api/graph/validator"
	"github.com/AchoArnold/ov-chipkaart-dashboard/backend/shared/errorhandler"
	"github.com/palantir/stacktrace"
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
	helpers      validator.Helpers
	errorHandler errorhandler.ErrorHandler
}

// New creates a new go validator
func New(db database.DB, helpers validator.Helpers, errorHandler errorhandler.ErrorHandler) validator.Validator {
	service := &GoValidator{db, helpers, errorHandler}
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

// ValidateStoreAnalzyeRequest validates the store analyze request input
func (service GoValidator) ValidateStoreAnalzyeRequest(input model.StoreAnalyzeRequestInput, _ language.Tag) validator.ValidationResult {
	v := govalidator.New(govalidator.Options{
		Data: &input,
		Rules: govalidator.MapData{
			"ovChipkaartUsername": []string{"min:6"},
			"ovChipkaartPassword": []string{"max:6"},
			"travelHistoryFile":   []string{"mime:text/csv"},
			"startDate":           []string{"required", "date:yyyy-mm-dd"},
			"endDate":             []string{"required", "date:yyyy-mm-dd"},
			"ovChipkaartNumber":   []string{"required", "min:16", "max:16", "numeric"},
		},
	})

	values := v.ValidateStruct()

	if input.OvChipkaartPassword == nil && input.OvChipkaartUsername == nil && input.TravelHistoryFile == nil {
		values.Add("ovChipkaartUsername", "You must provide either the username and password or the travel history csv file")
		values.Add("ovChipkaartPassword", "You must provide either the username and password or the travel history csv file")
		values.Add("travelHistoryFile", "You must provide either the username and password or the travel history csv file")
	}

	if len(values) > 0 {
		return service.urlValuesToResult(values)
	}

	startDate, _ := time.FromDate(input.StartDate)
	endDate, _ := time.FromDate(input.EndDate)
	if startDate.Unix() >= endDate.Unix() {
		values.Add("startDate", "The start date must be before the end date")
		values.Add("endDate", "The end date must be after the start date")
	}
	if len(values) > 0 {
		return service.urlValuesToResult(values)
	}

	// hours in 6 months
	sixMonthsInHours := 25 * 7 * 24 * internalTime.Hour
	if endDate.Sub(startDate) > sixMonthsInHours {
		values.Add("startDate", "The start date must be maximum 6 months before the end date")
		values.Add("endDate", "The end date must be maximum 6 months after the start date")
	}

	if len(values) > 0 {
		return service.urlValuesToResult(values)
	}

	if input.OvChipkaartPassword != nil || input.OvChipkaartUsername != nil {
		err := service.helpers.ValidateOvChipkaartCredentials(*input.OvChipkaartUsername, *input.OvChipkaartPassword)
		if err != nil {
			service.errorHandler.CaptureError(context.Background(), err)
			if stacktrace.GetCode(err) == ovchipkaart.ErrCodeUnauthorized {
				values.Add("ovChipkaartUsername", "Invalid OV Chipkaart username or password")
				values.Add("ovChipkaartPassword", "Invalid OV Chipkaart username or password")
			} else {
				values.Add("ovChipkaartUsername", "Internal error while verifying your ov chipkaart username")
				values.Add("ovChipkaartPassword", "Internal error while verifying your ov chipkaart password")
			}
		}
	}

	return service.urlValuesToResult(values)
}

// ValidateAnalzyeRequestsInput validates the analyze request query inputs
func (service GoValidator) ValidateAnalzyeRequestsInput(skip *int, take *int, orderBy *string, orderDirection *string, _ language.Tag) validator.ValidationResult {
	input := struct {
		skip, take              *int
		orderBy, orderDirection *string
	}{
		skip,
		take,
		orderBy,
		orderDirection,
	}

	v := govalidator.New(govalidator.Options{
		Data: &input,
		Rules: govalidator.MapData{
			"skip":           []string{"numeric", "min:0"},
			"take":           []string{"numeric", "min:0"},
			"orderBy":        []string{"in:createdAt,ovChipkaartNumber,startDate,EndDate"},
			"orderDirection": []string{"in:DESC,desc,ASC,asc"},
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
		if err == internalErrors.ErrEntityNotFound {
			return nil
		}

		if err != nil {
			service.errorHandler.CaptureError(context.Background(), stacktrace.Propagate(err, "cannot fetch user by email"))
			return fmt.Errorf("internal error when verifying the %s '%s'", field, value)
		}

		if message != "" {
			return errors.New(message)
		}

		return fmt.Errorf("a user already exist with the %s '%s'", field, value)
	}
}

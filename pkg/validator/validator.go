package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/tbe-team/raybot/pkg/sort"
)

var (
	AlphaNumberSpaceRegex = regexp.MustCompile("^[a-zA-Z0-9 ]+$")
)

func newValidator() *validator.Validate {
	v10 := validator.New()

	// Register custom validators here
	_ = v10.RegisterValidation("alphanumspace", func(fl validator.FieldLevel) bool {
		return AlphaNumberSpaceRegex.MatchString(fl.Field().String())
	})

	_ = v10.RegisterValidation("enum", validateEnum)
	_ = v10.RegisterValidation("sort", validateSortFields)
	return v10
}

func validateEnum(fl validator.FieldLevel) bool {
	type Enum interface {
		Validate() error
	}

	value, ok := fl.Field().Interface().(Enum)
	if !ok {
		return false
	}

	return value.Validate() == nil
}

func validateSortFields(fl validator.FieldLevel) bool {
	allowedFields := make(map[string]bool)
	for _, field := range strings.Split(fl.Param(), " ") {
		allowedFields[field] = true
	}

	value := fl.Field()
	if value.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < value.Len(); i++ {
		sort := value.Index(i).Interface().(sort.Sort)
		if !allowedFields[sort.Col] {
			return false
		}
	}
	return true
}

// IsValidationError checks if the given error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

func ValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "uuid":
		return "must be a valid UUID"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters long", fe.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", fe.Param())
	case "alphanumspace":
		return "must contain only alphanumeric characters and spaces"
	case "ip":
		return "must be a valid IP address"
	case "enum":
		return fmt.Sprintf("invalid enum value: %s", fe.Value())
	case "sort":
		return fmt.Sprintf("must contain only allowed sort fields: [%s]", fe.Param())
	default:
		return "is invalid"
	}
}

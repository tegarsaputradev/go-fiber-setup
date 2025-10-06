package helper

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

func ParseValidationError(err error) map[string]string {
	errs := make(map[string]string)

	if verrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range verrs {
			field := strings.ToLower(e.Field())
			var msg string

			switch e.Tag() {
			case "required":
				msg = field + " is required"
			case "email":
				msg = field + " must be a valid email"
			case "min":
				msg = field + " must be at least " + e.Param() + " characters"
			case "max":
				msg = field + " must be at most " + e.Param() + " characters"
			default:
				msg = field + " is not valid"
			}

			errs[field] = msg
		}
	}

	return errs
}

func ParseDuplicateError(err error) map[string]string {
	if err == nil || !strings.Contains(err.Error(), "Duplicate entry") {
		return nil
	}

	re := regexp.MustCompile(`for key '.*\.(.+)'`)
	match := re.FindStringSubmatch(err.Error())

	field := "field"
	if len(match) > 1 {
		parts := strings.Split(match[1], "_")
		field = parts[len(parts)-1]
	}

	return map[string]string{
		field: field + " already exists",
	}
}

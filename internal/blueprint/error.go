package blueprint

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ErrorMessages map[string]string

var errorMessages = ErrorMessages{
	"required":     "Field is required",
	"no-condition": "At least one of the conditions must not be empty",
}

type fieldError struct {
	Message   string
	Tag       string
	Value     string
	Namespace string
}

type Errors []fieldError

func handleErrors(err []validator.FieldError) error {

	var errors Errors
	for _, e := range err {
		message, ok := errorMessages[e.Tag()]
		if !ok {
			message = fmt.Sprintf("%s value is not valid", e.Field())
		}

		errors = append(errors, fieldError{
			Message:   message,
			Tag:       e.Tag(),
			Value:     fmt.Sprintf("\"%v\"", e.Value()),
			Namespace: e.Namespace(),
		})
	}

	return format(errors)
}

func format(errors Errors) error {
	formattedErrors := fmt.Sprintf("Validation errors [count=%d]:\n", len(errors))
	for _, e := range errors {
		formattedErrors += fmt.Sprintf("  - %s\n\tValue: %s\n\tField(s): %s\n\tCode: %s\n", e.Message, e.Value, e.Namespace, e.Tag)
	}

	return fmt.Errorf("%s", formattedErrors)
}

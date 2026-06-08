package validator

import (
	"net/http"

	"go.leoweyr.com/go-http-api-contract-io/v3/response"
)

const validationFailedMessage string = "VALIDATION_FAILED"
const validationFailedStatusCode int = http.StatusUnprocessableEntity

var _ response.RespondableError = (*ValidationError)(nil)

type ValidationError struct {
	details map[string]string
}

func NewValidationError(details map[string]string) *ValidationError {
	var validationError *ValidationError = &ValidationError{
		details: details,
	}

	return validationError
}

func (validationError *ValidationError) Error() string {
	return validationFailedMessage
}

func (validationError *ValidationError) StatusCode() int {
	return validationFailedStatusCode
}

func (validationError *ValidationError) Message() string {
	return validationFailedMessage
}

func (validationError *ValidationError) Details() map[string]string {
	return validationError.details
}

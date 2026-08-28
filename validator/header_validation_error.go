package validator

import (
	"net/http"

	"go.leoweyr.com/go-http-api-contract-io/v4/response"
)

const headerValidationFailedMessage string = "HEADER_VALIDATION_FAILED"
const headerValidationFailedStatusCode int = http.StatusBadRequest

var _ response.RespondableError = (*HeaderValidationError)(nil)

type HeaderValidationError struct {
	details map[string]string
}

func NewHeaderValidationError(details map[string]string) *HeaderValidationError {
	var headerValidationError *HeaderValidationError = &HeaderValidationError{
		details: details,
	}

	return headerValidationError
}

func (headerValidationError *HeaderValidationError) Error() string {
	return headerValidationFailedMessage
}

func (headerValidationError *HeaderValidationError) StatusCode() int {
	return headerValidationFailedStatusCode
}

func (headerValidationError *HeaderValidationError) Message() string {
	return headerValidationFailedMessage
}

func (headerValidationError *HeaderValidationError) Details() any {
	return headerValidationError.details
}

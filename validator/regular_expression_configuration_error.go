package validator

import (
	"fmt"
	"net/http"

	"go.leoweyr.com/go-http-api-contract-io/v5/response"
)

var _ response.RespondableError = (*RegularExpressionConfigurationError)(nil)

type RegularExpressionConfigurationError struct {
	cause error
}

func NewRegularExpressionConfigurationError(cause error) *RegularExpressionConfigurationError {
	return &RegularExpressionConfigurationError{cause: cause}
}

func (configurationError *RegularExpressionConfigurationError) Error() string {
	return fmt.Sprintf("regular expression configuration: %v", configurationError.cause)
}

func (configurationError *RegularExpressionConfigurationError) StatusCode() int {
	return http.StatusInternalServerError
}

func (configurationError *RegularExpressionConfigurationError) Message() string {
	return "INTERNAL_ERROR"
}

func (configurationError *RegularExpressionConfigurationError) Details() any {
	return nil
}

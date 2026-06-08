package response

type ErrorResponseBuilder struct{}

func NewErrorResponseBuilder() *ErrorResponseBuilder {
	var errorResponseBuilder *ErrorResponseBuilder = &ErrorResponseBuilder{}

	return errorResponseBuilder
}

func (errorResponseBuilder *ErrorResponseBuilder) buildErrorBody(message string, details map[string]string) ErrorBody {
	if details == nil {
		details = map[string]string{}
	}

	var errorBody ErrorBody = ErrorBody{
		Message: message,
		Details: details,
	}

	return errorBody
}

func (errorResponseBuilder *ErrorResponseBuilder) BuildErrorResponse(message string, details map[string]string) ErrorResponse {
	var errorResponse ErrorResponse = ErrorResponse{
		Error: errorResponseBuilder.buildErrorBody(message, details),
	}

	return errorResponse
}

func (errorResponseBuilder *ErrorResponseBuilder) BuildErrorResponseFromError(respondableError RespondableError) ErrorResponse {
	return errorResponseBuilder.BuildErrorResponse(respondableError.Message(), respondableError.Details())
}

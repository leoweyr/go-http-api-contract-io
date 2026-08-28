package response

import "reflect"

type ErrorResponseBuilder struct{}

var sharedErrorResponseBuilder *ErrorResponseBuilder = &ErrorResponseBuilder{}

// SharedErrorResponseBuilder returns the shared ErrorResponseBuilder singleton.
func SharedErrorResponseBuilder() *ErrorResponseBuilder {
	return sharedErrorResponseBuilder
}

func (errorResponseBuilder *ErrorResponseBuilder) buildErrorBody(message string, details any) ErrorBody {
	if details == nil {
		details = map[string]any{}
	} else {
		var detailsValue reflect.Value = reflect.ValueOf(details)

		switch detailsValue.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
			if detailsValue.IsNil() {
				details = map[string]any{}
			}
		}
	}

	var errorBody ErrorBody = ErrorBody{
		Message: message,
		Details: details,
	}

	return errorBody
}

func (errorResponseBuilder *ErrorResponseBuilder) BuildErrorResponse(message string, details any) ErrorResponse {
	var errorResponse ErrorResponse = ErrorResponse{
		Error: errorResponseBuilder.buildErrorBody(message, details),
	}

	return errorResponse
}

func (errorResponseBuilder *ErrorResponseBuilder) BuildErrorResponseFromError(respondableError RespondableError) ErrorResponse {
	return errorResponseBuilder.BuildErrorResponse(respondableError.Message(), respondableError.Details())
}

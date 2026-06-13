package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

type ResponseWriter struct{}

func NewResponseWriter() *ResponseWriter {
	var responseWriter *ResponseWriter = &ResponseWriter{}

	return responseWriter
}

func (responseWriter *ResponseWriter) WriteJSON(httpResponseWriter http.ResponseWriter, statusCode int, responseBody any) error {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(statusCode)

	var jsonEncoder *json.Encoder = json.NewEncoder(httpResponseWriter)
	var encodeError error = jsonEncoder.Encode(responseBody)

	return encodeError
}

func (responseWriter *ResponseWriter) WriteHeaders(httpResponseWriter http.ResponseWriter, payload any) error {
	var payloadValue reflect.Value = reflect.ValueOf(payload)

	if payloadValue.Kind() == reflect.Pointer {
		if payloadValue.IsNil() {
			return errors.New("header payload must not be a nil pointer")
		}

		payloadValue = payloadValue.Elem()
	}

	if payloadValue.Kind() != reflect.Struct {
		return errors.New("header payload must be a struct or a pointer to a struct")
	}

	var payloadType reflect.Type = payloadValue.Type()
	var fieldCount int = payloadType.NumField()
	var fieldIndex int = 0

	for fieldIndex = 0; fieldIndex < fieldCount; fieldIndex++ {
		var structField reflect.StructField = payloadType.Field(fieldIndex)
		var headerName string = structField.Tag.Get("header")

		if headerName == "" || headerName == "-" {
			continue
		}

		var fieldValue reflect.Value = payloadValue.Field(fieldIndex)

		if fieldValue.Kind() != reflect.String {
			continue
		}

		var headerValue string = fieldValue.String()

		if headerValue == "" {
			continue
		}

		httpResponseWriter.Header().Set(headerName, headerValue)
	}

	return nil
}

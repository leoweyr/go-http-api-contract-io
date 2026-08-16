package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

type ResponseReader struct{}

func NewResponseReader() *ResponseReader {
	var responseReader *ResponseReader = &ResponseReader{}

	return responseReader
}

func (responseReader *ResponseReader) ReadStatus(httpResponse *http.Response) int {
	return httpResponse.StatusCode
}

func (responseReader *ResponseReader) ReadJSON(httpResponse *http.Response, payload any) (int, error) {
	var jsonDecoder *json.Decoder = json.NewDecoder(httpResponse.Body)
	var decodeError error = jsonDecoder.Decode(payload)

	return httpResponse.StatusCode, decodeError
}

func (responseReader *ResponseReader) ReadHeaders(httpResponse *http.Response, payload any) error {
	var payloadValue reflect.Value = reflect.ValueOf(payload)

	if payloadValue.Kind() != reflect.Pointer || payloadValue.IsNil() {
		return errors.New("header payload must be a non-nil pointer to a struct")
	}

	var payloadElement reflect.Value = payloadValue.Elem()

	if payloadElement.Kind() != reflect.Struct {
		return errors.New("header payload must point to a struct")
	}

	var payloadType reflect.Type = payloadElement.Type()
	var fieldCount int = payloadType.NumField()
	var fieldIndex int = 0

	for fieldIndex = 0; fieldIndex < fieldCount; fieldIndex++ {
		var structField reflect.StructField = payloadType.Field(fieldIndex)
		var headerName string = structField.Tag.Get("header")

		if headerName == "" || headerName == "-" {
			continue
		}

		var fieldValue reflect.Value = payloadElement.Field(fieldIndex)

		if !fieldValue.CanSet() || fieldValue.Kind() != reflect.String {
			continue
		}

		fieldValue.SetString(httpResponse.Header.Get(headerName))
	}

	return nil
}

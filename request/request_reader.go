package request

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

type RequestReader struct{}

func NewRequestReader() *RequestReader {
	var requestReader *RequestReader = &RequestReader{}

	return requestReader
}

func (requestReader *RequestReader) ReadJSON(httpRequest *http.Request, payload any) error {
	var jsonDecoder *json.Decoder = json.NewDecoder(httpRequest.Body)

	jsonDecoder.DisallowUnknownFields()

	var decodeError error = jsonDecoder.Decode(payload)

	return decodeError
}

func (requestReader *RequestReader) ReadHeader(httpRequest *http.Request, payload any) error {
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

		fieldValue.SetString(httpRequest.Header.Get(headerName))
	}

	return nil
}

package request

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
)

type RequestWriter struct{}

func NewRequestWriter() *RequestWriter {
	var requestWriter *RequestWriter = &RequestWriter{}

	return requestWriter
}

func (requestWriter *RequestWriter) WriteJSON(httpRequest *http.Request, payload any) error {
	var bodyBuffer *bytes.Buffer = new(bytes.Buffer)
	var jsonEncoder *json.Encoder = json.NewEncoder(bodyBuffer)
	var encodeError error = jsonEncoder.Encode(payload)

	if encodeError != nil {
		return encodeError
	}

	var bodySnapshot []byte = bodyBuffer.Bytes()

	httpRequest.Body = io.NopCloser(bytes.NewReader(bodySnapshot))
	httpRequest.ContentLength = int64(len(bodySnapshot))

	httpRequest.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(bodySnapshot)), nil
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	return nil
}

func (requestWriter *RequestWriter) WriteHeader(httpRequest *http.Request, payload any) error {
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

		httpRequest.Header.Set(headerName, headerValue)
	}

	return nil
}

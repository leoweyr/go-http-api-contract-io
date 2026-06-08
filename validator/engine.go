package validator

import (
	"reflect"
	"strings"

	validation "github.com/go-playground/validator/v10"

	"go.leoweyr.com/go-http-api-contract-io/v3/response"
)

type Engine struct {
	validationLibraryEngine *validation.Validate
	errorFormatter          *ErrorFormatter
}

func extractJSONTagName(structField reflect.StructField) string {
	var rawJSONTagName string = structField.Tag.Get("json")

	if rawJSONTagName == "" {
		return structField.Name
	}

	var jsonTagNameSegments []string = strings.Split(rawJSONTagName, ",")
	var primaryJSONTagName string = jsonTagNameSegments[0]

	if primaryJSONTagName == "" || primaryJSONTagName == "-" {
		return structField.Name
	}

	return primaryJSONTagName
}

func registerJSONTagName(validationLibraryEngine *validation.Validate) {
	validationLibraryEngine.RegisterTagNameFunc(extractJSONTagName)
}

func NewEngine() *Engine {
	var validationLibraryEngine *validation.Validate = validation.New()
	registerJSONTagName(validationLibraryEngine)

	var errorFormatter *ErrorFormatter = NewErrorFormatter()
	var engine *Engine = &Engine{
		validationLibraryEngine: validationLibraryEngine,
		errorFormatter:          errorFormatter,
	}

	return engine
}

func (engine *Engine) ValidateStruct(payload any) error {
	return engine.validationLibraryEngine.Struct(payload)
}

func (engine *Engine) FormatValidationError(validationError error) response.RespondableError {
	var details map[string]string = engine.errorFormatter.FormatValidationErrors(validationError)

	return NewValidationError(details)
}

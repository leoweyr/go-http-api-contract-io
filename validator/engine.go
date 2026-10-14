package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"unicode/utf8"

	validation "github.com/go-playground/validator/v10"

	"go.leoweyr.com/go-http-api-contract-io/v5/response"
)

type Engine struct {
	validationLibraryEngine             *validation.Validate
	errorFormatter                      *ErrorFormatter
	patternMutex                        sync.RWMutex
	compiledPatterns                    map[string]*regexp.Regexp
	regularExpressionPreparationResults map[reflect.Type]error
}

var sharedEngine *Engine = newEngine()

// SharedEngine returns the shared Engine singleton.
func SharedEngine() *Engine {
	return sharedEngine
}

func extractFieldDisplayName(structField reflect.StructField) string {
	var rawJSONTagName string = structField.Tag.Get("json")

	if rawJSONTagName == "" {
		var rawHeaderTagName string = structField.Tag.Get("header")

		if rawHeaderTagName != "" && rawHeaderTagName != "-" {
			return rawHeaderTagName
		}

		return structField.Name
	}

	var jsonTagNameSegments []string = strings.Split(rawJSONTagName, ",")
	var primaryJSONTagName string = jsonTagNameSegments[0]

	if primaryJSONTagName == "" || primaryJSONTagName == "-" {
		return structField.Name
	}

	return primaryJSONTagName
}

func registerFieldDisplayName(validationLibraryEngine *validation.Validate) {
	validationLibraryEngine.RegisterTagNameFunc(extractFieldDisplayName)
}

func validateUTF8String(field validation.FieldLevel) bool {
	var fieldValue reflect.Value = field.Field()

	return fieldValue.Kind() == reflect.String && utf8.ValidString(fieldValue.String())
}

func newEngine() *Engine {
	var validationLibraryEngine *validation.Validate = validation.New()
	registerFieldDisplayName(validationLibraryEngine)

	var errorFormatter *ErrorFormatter = NewErrorFormatter()
	var engine *Engine = &Engine{
		validationLibraryEngine:             validationLibraryEngine,
		errorFormatter:                      errorFormatter,
		compiledPatterns:                    make(map[string]*regexp.Regexp),
		regularExpressionPreparationResults: make(map[reflect.Type]error),
	}

	var registrationError error = validationLibraryEngine.RegisterValidation("regexp", engine.matchRegularExpression)

	if registrationError != nil {
		panic(fmt.Errorf("register regexp validation: %w", registrationError))
	}

	registrationError = validationLibraryEngine.RegisterValidation("utf8", validateUTF8String)

	if registrationError != nil {
		panic(fmt.Errorf("register utf8 validation: %w", registrationError))
	}

	return engine
}

func (engine *Engine) matchRegularExpression(field validation.FieldLevel) bool {
	var fieldValue reflect.Value = field.Field()

	if fieldValue.Kind() != reflect.String {
		return false
	}

	engine.patternMutex.RLock()
	var pattern *regexp.Regexp = engine.compiledPatterns[field.Param()]
	engine.patternMutex.RUnlock()

	return pattern != nil && pattern.MatchString(fieldValue.String())
}

func (engine *Engine) precompileRegularExpressionsForType(modelType reflect.Type, visitedTypes map[reflect.Type]bool) error {
	if modelType == nil || visitedTypes[modelType] {
		return nil
	}

	visitedTypes[modelType] = true

	switch modelType.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Array:
		return engine.precompileRegularExpressionsForType(modelType.Elem(), visitedTypes)

	case reflect.Map:
		var preparationError error = engine.precompileRegularExpressionsForType(modelType.Key(), visitedTypes)

		if preparationError != nil {
			return preparationError
		}

		return engine.precompileRegularExpressionsForType(modelType.Elem(), visitedTypes)

	case reflect.Struct:
		var index int

		for index = 0; index < modelType.NumField(); index++ {
			var field reflect.StructField = modelType.Field(index)
			var declaration string = field.Tag.Get("validate")

			if field.PkgPath != "" || declaration == "-" {
				continue
			}

			var rule string

			for _, rule = range strings.Split(declaration, ",") {
				var alternative string

				for _, alternative = range strings.Split(rule, "|") {
					var ruleName string
					var expression string
					var hasExpression bool
					ruleName, expression, hasExpression = strings.Cut(alternative, "=")

					if ruleName != "regexp" {
						continue
					}

					if !hasExpression {
						return NewRegularExpressionConfigurationError(fmt.Errorf("field %s requires a regular expression parameter", field.Name))
					}

					expression = strings.ReplaceAll(strings.ReplaceAll(expression, "0x2C", ","), "0x7C", "|")

					if engine.compiledPatterns[expression] != nil {
						continue
					}

					var pattern *regexp.Regexp
					var compilationError error
					pattern, compilationError = regexp.Compile(expression)

					if compilationError != nil {
						return NewRegularExpressionConfigurationError(fmt.Errorf("field %s: %w", field.Name, compilationError))
					}

					engine.compiledPatterns[expression] = pattern
				}
			}

			var preparationError error = engine.precompileRegularExpressionsForType(field.Type, visitedTypes)

			if preparationError != nil {
				return preparationError
			}
		}
	}

	return nil
}

// PrecompileRegularExpressionsForStruct precompiles and caches regular expressions declared on a struct and its nested types.
func (engine *Engine) PrecompileRegularExpressionsForStruct(payload any) error {
	var modelType reflect.Type = reflect.TypeOf(payload)
	engine.patternMutex.RLock()
	var preparationError error
	var prepared bool
	preparationError, prepared = engine.regularExpressionPreparationResults[modelType]
	engine.patternMutex.RUnlock()

	if prepared {
		return preparationError
	}

	engine.patternMutex.Lock()

	defer engine.patternMutex.Unlock()

	preparationError, prepared = engine.regularExpressionPreparationResults[modelType]

	if prepared {
		return preparationError
	}

	preparationError = engine.precompileRegularExpressionsForType(modelType, make(map[reflect.Type]bool))
	engine.regularExpressionPreparationResults[modelType] = preparationError

	return preparationError
}

// ValidateStruct precompiles declared regular expressions and validates a DTO's field values.
func (engine *Engine) ValidateStruct(payload any) error {
	var preparationError error = engine.PrecompileRegularExpressionsForStruct(payload)

	if preparationError != nil {
		return preparationError
	}

	return engine.validationLibraryEngine.Struct(payload)
}

func (engine *Engine) FormatValidationError(validationError error) response.RespondableError {
	var configurationError *RegularExpressionConfigurationError

	if errors.As(validationError, &configurationError) {
		return configurationError
	}

	var details map[string]string = engine.errorFormatter.FormatValidationErrors(validationError)

	return NewValidationError(details)
}

func (engine *Engine) FormatHeaderValidationError(validationError error) response.RespondableError {
	var configurationError *RegularExpressionConfigurationError

	if errors.As(validationError, &configurationError) {
		return configurationError
	}

	var details map[string]string = engine.errorFormatter.FormatValidationErrors(validationError)

	return NewHeaderValidationError(details)
}

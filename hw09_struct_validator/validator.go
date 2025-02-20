package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidLen    = errors.New("invalid length")
	ErrInvalidMin    = errors.New("invalid min")
	ErrInvalidMax    = errors.New("invalid max")
	ErrInvalidIn     = errors.New("invalid in")
	ErrInvalidRegexp = errors.New("invalid regexp")
	ErrInvalidType   = errors.New("invalid type")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var errorMessages []string

	for _, ve := range v {
		errorMessages = append(errorMessages, fmt.Sprintf("Field '%s': %s", ve.Field, ve.Err))
	}

	return strings.Join(errorMessages, "; ")
}

func Validate(v interface{}) error {
	var validationErrors ValidationErrors
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return errors.New("input is not struct")
	}

	valType := val.Type()

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)
		tag := field.Tag.Get("validate")

		if tag != "" {
			tagParts := strings.Split(tag, "|")
			fieldErr := ValidationError{
				Field: field.Name,
				Err:   nil,
			}

			for _, part := range tagParts {
				value := val.Field(i)

				if value.Kind() == reflect.Slice {
					for j := 0; j < value.Len(); j++ {
						err := validateField(part, value.Index(j).Interface(), &fieldErr)

						if err != nil {
							return err
						}
					}
				} else {
					err := validateField(part, value.Interface(), &fieldErr)

					if err != nil {
						return err
					}
				}
			}

			if fieldErr.Err != nil {
				validationErrors = append(validationErrors, fieldErr)
			}
		}
	}

	return validationErrors
}

func validateField(tag string, value interface{}, fieldErr *ValidationError) error {
	tagPairs := strings.Split(tag, ":")
	tagName, tagValue := tagPairs[0], tagPairs[1]
	var validError error

	switch tagName {
	case "len":
		val, err := strconv.Atoi(tagValue)
		if err != nil {
			return err
		}
		if str, ok := value.(string); ok {
			if len(str) != val {
				validError = fmt.Errorf("%s length should be %v: %w", str, tagValue, ErrInvalidLen)
			}
		} else {
			return fmt.Errorf("type of value must be string: %w", ErrInvalidType)
		}
	case "min":
		val, err := strconv.Atoi(tagValue)
		if err != nil {
			return err
		}
		if valInt, ok := value.(int); ok {
			if valInt < val {
				validError = fmt.Errorf("%v less then %v: %w", valInt, tagValue, ErrInvalidMin)
			}
		} else {
			return fmt.Errorf("type of value must be int: %w", ErrInvalidType)
		}
	case "max":
		val, err := strconv.Atoi(tagValue)
		if err != nil {
			return err
		}
		if valInt, ok := value.(int); ok {
			if valInt > val {
				validError = fmt.Errorf("%v more then %v: %w", valInt, tagValue, ErrInvalidMax)
			}
		} else {
			return fmt.Errorf("type of value must be int: %w", ErrInvalidType)
		}
	case "in":
		options := strings.Split(tagValue, ",")
		valueStr := fmt.Sprint(value)

		for _, option := range options {
			if option == valueStr {
				return nil
			}
		}

		validError = fmt.Errorf("%s not in %v: %w", valueStr, tagValue, ErrInvalidIn)
	case "regexp":
		regex, err := regexp.Compile(tagValue)
		if err != nil {
			return fmt.Errorf("invalid regex pattern")
		}
		if str, ok := value.(string); ok {
			if str != "" && !regex.MatchString(str) {
				validError = fmt.Errorf("value '%s' does not match pattern '%s': %w", str, tagValue, ErrInvalidRegexp)
			}
		} else {
			return fmt.Errorf("type of value must be string: %w", ErrInvalidType)
		}

	}

	if validError != nil {
		if fieldErr.Err != nil {
			fieldErr.Err = fmt.Errorf("%v: %w", validError, fieldErr.Err)
		} else {
			fieldErr.Err = validError
		}
	}

	return nil
}

package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

var (
	errorUnsupportedType              = errors.New("unsupported element type")
	errorValidateStringLength         = errors.New("bad string length")
	errorValidateStringRegexp         = errors.New("bad regexp matching")
	errorValidateIn                   = errors.New("no matches")
	errorValidateUnsupportedValueType = errors.New("unsupported value type")
	errorValidateMin                  = errors.New("value is less than allowed")
	errorValidateMax                  = errors.New("value is greater than allowed")

	errorValidationErrors = ValidationErrors{
		{
			Field: "test",
			Err:   nil,
		},
	}
)

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, errRes := range v {
		sb.WriteString(fmt.Sprintf("%s -> %v", errRes.Field, errRes.Err))
	}
	return sb.String()
}

func contains(value interface{}, s ...interface{}) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}
	return false
}

func getInternalTag(tag string) map[string]string {
	res := make(map[string]string)
	elements := strings.Split(tag, "|")
	for _, el := range elements {
		buf := strings.Split(el, ":")
		switch len(buf) {
		case 0:
			return res
		case 1:
			res[buf[0]] = ""
		default:
			res[buf[0]] = strings.Join(buf[1:], ":")
		}
	}
	return res
}

func validateString(name string, value string, tag map[string]string) error {
	validationErrors := ValidationErrors{}
	count := tag["len"]
	if count != "" {
		countInt, err := strconv.Atoi(count)
		if err != nil {
			return err
		}
		if countInt < len(value) {
			validationError := ValidationError{}
			validationError.Field = name
			validationError.Err = errorValidateStringLength
			validationErrors = append(validationErrors, validationError)
		}
	}

	stringPattern := tag["regexp"]
	if stringPattern != "" {
		pattern, err := regexp.Compile(stringPattern)
		if err != nil {
			return err
		}
		if !pattern.MatchString(value) {
			validationError := ValidationError{}
			validationError.Field = name
			validationError.Err = errorValidateStringRegexp
			validationErrors = append(validationErrors, validationError)
		}
	}

	variantsStr := tag["in"]
	if variantsStr != "" {
		variants := strings.Split(variantsStr, ",")
		if !contains(value, variants) {
			validationError := ValidationError{}
			validationError.Field = name
			validationError.Err = errorValidateIn
			validationErrors = append(validationErrors, validationError)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateInt(name string, value int, tag map[string]string) error {
	validationErrors := ValidationErrors{}
	minSt := tag["min"]
	if minSt != "" {
		min, err := strconv.Atoi(minSt)
		if err != nil {
			return err
		}
		if min > value {
			validationError := ValidationError{}
			validationError.Field = name
			validationError.Err = errorValidateMin
			validationErrors = append(validationErrors, validationError)
		}
	}

	maxSt := tag["max"]
	if maxSt != "" {
		max, err := strconv.Atoi(maxSt)
		if err != nil {
			return err
		}
		if max < value {
			validationError := ValidationError{}
			validationError.Field = name
			validationError.Err = errorValidateMax
			validationErrors = append(validationErrors, validationError)
		}
	}

	variantsStr := tag["in"]
	if variantsStr != "" {
		variantsSt := strings.Split(variantsStr, ",")
		variants := make([]int, len(variantsSt))
		for ind, el := range variantsSt {
			j, err := strconv.Atoi(el)
			if err != nil {
				return err
			}
			variants[ind] = j
		}
		if !contains(variants, value) {
			validationError := ValidationError{}
			validationError.Field = name
			validationError.Err = errorValidateIn
			validationErrors = append(validationErrors, validationError)
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateFromTag(name string, value interface{}, tag string) error {
	validationErrors := ValidationErrors{}

	validatorTag := getInternalTag(tag)

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		v := reflect.ValueOf(value)
		for i := 0; i < v.Len(); i++ {
			err := validateFromTag(name, v.Index(i).Interface(), tag)
			if errors.As(err, &errorValidationErrors) {
				validationErrors = append(validationErrors, err.(ValidationErrors)...) //nolint:errorlint
			} else if err != nil {
				return err
			}
		}
		return validationErrors
	}

	switch data := value.(type) {
	case string:
		err := validateString(name, data, validatorTag)
		if errors.As(err, &errorValidationErrors) {
			validationErrors = append(validationErrors, err.(ValidationErrors)...) //nolint:errorlint
		} else if err != nil {
			return err
		}
	case int:
		err := validateInt(name, data, validatorTag)
		if errors.As(err, &errorValidationErrors) {
			validationErrors = append(validationErrors, err.(ValidationErrors)...) //nolint:errorlint
		} else if err != nil {
			return err
		}
	default:
		_, found := validatorTag["nested"]
		if !found {
			validationError := ValidationError{}
			validationError.Field = name
			validationError.Err = errorValidateUnsupportedValueType
			validationErrors = append(validationErrors, validationError)
		} else {
			err := validateStruct(value)
			if errors.As(err, &errorValidationErrors) {
				validationErrors = append(validationErrors, err.(ValidationErrors)...) //nolint:errorlint
			} else if err != nil {
				return err
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateField(name string, value interface{}, tag reflect.StructTag) error {
	validationErrors := ValidationErrors{}
	if tag.Get("validate") != "" {
		validate, ok := tag.Lookup("validate")
		if ok {
			err := validateFromTag(name, value, validate)
			if errors.As(err, &errorValidationErrors) {
				validationErrors = append(validationErrors, err.(ValidationErrors)...) //nolint:errorlint
			} else if err != nil {
				return err
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateStruct(v interface{}) error {
	vE := ValidationErrors{}
	el := reflect.ValueOf(v)
	for i := 0; i < el.NumField(); i++ {
		fieldDetail := el.Type().Field(i)
		valF := el.Field(i)
		if valF.IsValid() && fieldDetail.IsExported() {
			fE := validateField(fieldDetail.Name, valF.Interface(), fieldDetail.Tag)
			if errors.As(fE, &errorValidationErrors) {
				vE = append(vE, fE.(ValidationErrors)...) //nolint:errorlint
			} else if fE != nil {
				return fE
			}
		}
	}
	if len(vE) > 0 {
		return vE
	}
	return nil
}

func Validate(v interface{}) error {
	elem := reflect.ValueOf(v)
	if elem.Kind() != reflect.Struct {
		return errorUnsupportedType
	}
	err := validateStruct(v)
	if errors.As(err, &errorValidationErrors) || err != nil {
		return err
	}

	return nil
}

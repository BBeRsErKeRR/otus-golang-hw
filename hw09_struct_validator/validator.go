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
	errorUnsupportedType              = errors.New("unsupported element type")
	errorValidateStringLength         = errors.New("bad string length")
	errorValidateStringRegexp         = errors.New("bad regexp matching")
	errorValidateIn                   = errors.New("no matches")
	errorValidateUnsupportedValueType = errors.New("unsupported value type")
	errorValidateMin                  = errors.New("value is less than allowed")
	errorValidateMax                  = errors.New("value is greater than allowed")
	errorValidateNoSpaces             = errors.New("value is contains spaces")
	errorValidateOdd                  = errors.New("value is even")
	errorValidateEven                 = errors.New("value is odd")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type Validator struct {
	validationErrors ValidationErrors
}

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, errRes := range v {
		sb.WriteString(fmt.Sprintf("%s -> %v", errRes.Field, errRes.Err))
	}
	return sb.String()
}

func contains[T comparable](s []T, value T) bool {
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

func (v *Validator) getErrors() ValidationErrors {
	return v.validationErrors
}

func (v *Validator) appendError(name string, e error) {
	validationError := ValidationError{}
	validationError.Field = name
	validationError.Err = e
	v.validationErrors = append(v.validationErrors, validationError)
}

func (v *Validator) validateString(name, value string, tag map[string]string) error {
	count := tag["len"]
	if count != "" {
		countInt, err := strconv.Atoi(count)
		if err != nil {
			return err
		}
		if countInt < len(value) {
			v.appendError(name, errorValidateStringLength)
		}
	}

	stringPattern := tag["regexp"]
	if stringPattern != "" {
		pattern, err := regexp.Compile(stringPattern)
		if err != nil {
			return err
		}
		if !pattern.MatchString(value) {
			v.appendError(name, errorValidateStringRegexp)
		}
	}

	variantsStr := tag["in"]
	if variantsStr != "" {
		variants := strings.Split(variantsStr, ",")
		if !contains(variants, value) {
			v.appendError(name, errorValidateIn)
		}
	}

	_, isNoSpaces := tag["nospaces"]
	if isNoSpaces {
		if strings.Contains(value, " ") {
			v.appendError(name, errorValidateNoSpaces)
		}
	}

	return nil
}

func (v *Validator) validateInt64(name string, value int64, tag map[string]string) error {
	minSt := tag["min"]
	if minSt != "" {
		min, err := strconv.Atoi(minSt)
		if err != nil {
			return err
		}
		if int64(min) > value {
			v.appendError(name, errorValidateMin)
		}
	}

	maxSt := tag["max"]
	if maxSt != "" {
		max, err := strconv.Atoi(maxSt)
		if err != nil {
			return err
		}
		if int64(max) < value {
			v.appendError(name, errorValidateMax)
		}
	}

	variantsStr := tag["in"]
	if variantsStr != "" {
		variantsSt := strings.Split(variantsStr, ",")
		variants := make([]int64, len(variantsSt))
		for ind, el := range variantsSt {
			j, err := strconv.Atoi(el)
			if err != nil {
				return err
			}
			variants[ind] = int64(j)
		}
		if !contains(variants, value) {
			v.appendError(name, errorValidateIn)
		}
	}

	_, odd := tag["odd"]
	if odd {
		if value%2 == 0 {
			v.appendError(name, errorValidateOdd)
		}
	}

	_, even := tag["even"]
	if even {
		if value%2 != 0 {
			v.appendError(name, errorValidateEven)
		}
	}

	return nil
}

func (v *Validator) validateInt(name string, value int, tag map[string]string) error {
	return v.validateInt64(name, int64(value), tag)
}

func (v *Validator) validateFromTag(name string, value interface{}, tag string) error {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		val := reflect.ValueOf(value)
		for i := 0; i < val.Len(); i++ {
			err := v.validateFromTag(name, val.Index(i).Interface(), tag)
			if err != nil {
				return err
			}
		}
		return nil
	}

	validatorTag := getInternalTag(tag)
	switch data := value.(type) {
	case string:
		err := v.validateString(name, data, validatorTag)
		if err != nil {
			return err
		}
	case int:
		err := v.validateInt(name, data, validatorTag)
		if err != nil {
			return err
		}
	case int64:
		err := v.validateInt64(name, data, validatorTag)
		if err != nil {
			return err
		}
	default:
		_, found := validatorTag["nested"]
		if !found {
			v.appendError(name, errorValidateUnsupportedValueType)
		} else {
			err := v.ValidateStruct(value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *Validator) validateField(name string, value interface{}, tag reflect.StructTag) error {
	validate, ok := tag.Lookup("validate")
	if ok {
		err := v.validateFromTag(name, value, validate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Validator) ValidateStruct(st interface{}) error {
	el := reflect.ValueOf(st)
	for i := 0; i < el.NumField(); i++ {
		fieldDetail := el.Type().Field(i)
		valF := el.Field(i)
		if fieldDetail.IsExported() {
			fE := v.validateField(fieldDetail.Name, valF.Interface(), fieldDetail.Tag)
			if fE != nil {
				return fE
			}
		}
	}
	return nil
}

func Validate(v interface{}) error {
	elem := reflect.ValueOf(v)
	if elem.Kind() != reflect.Struct {
		return errorUnsupportedType
	}
	validator := Validator{}
	err := validator.ValidateStruct(v)
	if err != nil {
		return err
	}
	valErrors := validator.getErrors()
	if len(valErrors) > 0 {
		return valErrors
	}

	return nil
}

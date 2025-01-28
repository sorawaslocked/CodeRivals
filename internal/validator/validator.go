package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) addFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.addFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinChars(value string, min uint32) bool {
	return utf8.RuneCountInString(value) >= int(min)
}

func MaxChars(value string, max uint32) bool {
	return utf8.RuneCountInString(value) <= int(max)
}

func Equal(value1, value2 string) bool {
	return value1 == value2
}

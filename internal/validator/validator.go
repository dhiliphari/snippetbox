package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

package validator

import (
	"strings"
	"unicode/utf8"
)

type FormValidator struct {
	FieldErrors map[string]string
}

func (v *FormValidator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *FormValidator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exist := v.FieldErrors[field]; !exist {
		v.FieldErrors[field] = message
	}
}

func (v *FormValidator) CheckField(ok bool, field string, message string) {
	if !ok {
		v.AddFieldError(field, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(strings.TrimSpace(value)) <= max
}

func ValueIn(value int, permitted ...int) bool {
	for _, v := range permitted {
		if value == v {
			return true
		}
	}
	return false
}

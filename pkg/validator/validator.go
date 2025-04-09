package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRGX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	if v.NonFieldErrors == nil {
		v.NonFieldErrors = []string{}
	}

	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exist := v.FieldErrors[field]; !exist {
		v.FieldErrors[field] = message
	}
}

func (v *Validator) CheckField(ok bool, field string, message string) {
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

func ValueIn[T comparable](value T, permitted ...T) bool {
	for _, v := range permitted {
		if value == v {
			return true
		}
	}
	return false
}

func MinChars(value string, min int) bool {
	return utf8.RuneCountInString(strings.TrimSpace(value)) >= min
}

func Matches(value string, rgx *regexp.Regexp) bool {
	return rgx.MatchString(strings.TrimSpace(value))
}

func Equal[T comparable](value1, value2 T) bool {
	return value1 == value2
}

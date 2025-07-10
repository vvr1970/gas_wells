package validation

import (
	"regexp"
	"strings"
	"unicode"
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Правила валидации
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MinLength(value string, n int) bool {
	return len(value) >= n
}

func MaxLength(value string, n int) bool {
	return len(value) <= n
}

func Between(value, min, max float64) bool {
	return value >= min && value <= max
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func IsEmail(value string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(value)
}

func IsPassword(value string) bool {
	var (
		hasMinLen  = len(value) >= 8
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range value {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func In[T comparable](value T, allowedValues ...T) bool {
	for _, v := range allowedValues {
		if value == v {
			return true
		}
	}
	return false
}

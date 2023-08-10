package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"unicode"
	"unicode/utf8"
)

func passwordValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		value := field.String()
		if n := utf8.RuneCountInString(value); n < 8 || n > 63 {
			return false
		}

		var (
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

		return hasUpper && hasLower && hasNumber && hasSpecial
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func wifiEncryptionValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		value := field.String()
		for _, encryption := range []string{"psk", "psk2"} {
			if value == encryption {
				return true
			}
		}

		return false
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func wifiTTLValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Uint:
		value := fl.Field().Uint()
		for _, ttl := range []uint64{5, 10, 15} {
			if value == ttl {
				return true
			}
		}

		return false
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func wifiChannelValidation(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		value := field.String()
		for _, channel := range []string{"auto", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10"} {
			if value == channel {
				return true
			}
		}

		return false
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

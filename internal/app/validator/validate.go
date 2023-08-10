package validator

import (
	"fmt"
	"github.com/Sergii-Kirichok/DTekSpeachParser/internal/app/translator"
	ut "github.com/go-playground/universal-translator"
	enT "github.com/go-playground/validator/v10/translations/en"
	ruT "github.com/go-playground/validator/v10/translations/ru"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func Struct(s interface{}) error { return validate.Struct(s) }

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	registerDefaultTranslations()
	registerConcreteTranslations()
	registerAliases()
	registerPersonalValidations()
}

func registerDefaultTranslations() {
	for locale, recorder := range map[string]func(v *validator.Validate, trans ut.Translator) error{
		"en": enT.RegisterDefaultTranslations,
		"ru": ruT.RegisterDefaultTranslations,
	} {
		trans, found := translator.Get(locale)
		if !found {
			log.Panicf("translator for locale %q not found", locale)
		}

		panicIf(recorder(validate, trans))
	}
}

func registerConcreteTranslations() {
	for tag, trans := range map[string]translator.Translation{
		"ipv4": {
			"en": "Please provide a valid IPv4 address",
			"ru": "Пожалуйста укажите действительный адрес IPv4",
		},
	} {
		registerTranslations(tag, trans)
	}
}

func registerAliases() {
	for alias, st := range map[string]struct {
		tags  string
		trans translator.Translation
	}{
		"wifi_ssid": {
			tags: "required,min=1,max=32",
			trans: translator.Translation{
				"en": "SSID must be from 1 to 32 characters long",
				"ru": "SSID должен иметь длину от 1 до 32 символов",
			},
		},
	} {
		validate.RegisterAlias(alias, st.tags)
		registerTranslations(alias, st.trans)
	}
}

func registerPersonalValidations() {
	for tag, st := range map[string]struct {
		fn    validator.Func
		trans translator.Translation
	}{
		"password": {
			fn: passwordValidation,
			trans: translator.Translation{
				"en": "Password must be at least 8 characters but less than 63 with 1 upper case letter, 1 number and 1 special character",
				"ru": "Пароль должен содержать не менее 8 символов но не более 63, включая заглавные буквы, цифры и специальные символы",
			},
		},
		"wifi_channel": {
			fn: wifiChannelValidation,
			trans: translator.Translation{
				"en": fmt.Sprintf(errNotFromListEnFormat, "channel"),
				"ru": fmt.Sprintf(errNotFromListRuFormat, "канал"),
			},
		},
		"wifi_ttl": {
			fn: wifiTTLValidation,
			trans: translator.Translation{
				"en": fmt.Sprintf(errNotFromListEnFormat, "ttl"),
				"ru": fmt.Sprintf(errNotFromListRuFormat, "время соединения"),
			},
		},
		"wifi_encryption": {
			fn: wifiEncryptionValidation,
			trans: translator.Translation{
				"en": fmt.Sprintf(errNotFromListEnFormat, "encryption"),
				"ru": fmt.Sprintf(errNotFromListRuFormat, "тип шифрования"),
			},
		},
	} {
		panicIf(validate.RegisterValidation(tag, st.fn))
		registerTranslations(tag, st.trans)
	}
}

func registerTranslations(tag string, translations translator.Translation) {
	for locale, translation := range translations {
		trans, found := translator.Get(locale)
		if !found {
			log.Panicf("translator for locale %q not found", locale)
		}

		registerTranslation(tag, trans, translation)
	}
}

func registerTranslation(tag string, trans ut.Translator, translation string) {
	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, translation, true)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
		if err != nil {
			return fe.(error).Error()
		}

		return t
	}

	panicIf(validate.RegisterTranslation(tag, trans, registerFn, transFn))
}

func panicIf(err error) {
	if err != nil {
		log.Panic(err)
	}
}

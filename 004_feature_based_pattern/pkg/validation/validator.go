package validation

import (
	"log"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/ru"
	"github.com/go-playground/locales/tr"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	ruTranslations "github.com/go-playground/validator/v10/translations/ru"
	trTranslations "github.com/go-playground/validator/v10/translations/tr"
)

var translator ut.Translator

func Init(lang string) *validator.Validate {
	// Override default validator engine
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("Validator engine is not found")
	}

	// tag alanlarini register et
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// json tag'ini al
		tag := fld.Tag.Get("json")

		if tag == "" {
			return fld.Name
		}

		// json tag varsa, virgule kadar olan kismi al
		name := strings.Split(tag, ",")[0]
		if name == "-" || name == "" {
			return fld.Name
		}

		return name
	})

	// Translator init
	trLocale := tr.New()
	enLocale := en.New()
	ruLocale := ru.New()

	// fallback locale: en, supported locales: tr, ru
	uni := ut.New(enLocale, trLocale, enLocale, ruLocale)

	var found bool
	translator, found = uni.GetTranslator(lang)

	if !found {
		translator, _ = uni.GetTranslator("en")
	}

	var err error
	switch lang {
	case "tr":
		err = trTranslations.RegisterDefaultTranslations(v, translator)
	case "ru":
		err = ruTranslations.RegisterDefaultTranslations(v, translator)
	default: // "en" veya diğer durumlar için varsayılan İngilizce
		err = enTranslations.RegisterDefaultTranslations(v, translator)
	}

	if err != nil {
		log.Printf("Error registering translations: %v", err)
	}

	return v
}

func MapValidationErrors(ve validator.ValidationErrors) map[string][]string {
	out := make(map[string][]string) // ram de map olusturuldu

	for _, fe := range ve {
		field := fe.Field() // Name, Email, Age
		msg := fe.Translate(translator)

		out[field] = append(out[field], msg)
	}
	return out
}

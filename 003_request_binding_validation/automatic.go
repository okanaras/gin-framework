package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
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

func main() {

	cfg := loadConfig()

	initValidator(cfg)

	r := gin.Default()

	r.POST("/users", createUserHandler)

	r.Run(":8080")
}

type APIErrorResponse struct {
	Message string              `json:"message"`          // struct tag
	Errors  map[string][]string `json:"errors,omitempty"` // omitempty: alan boşsa JSON çıktısında yer almaz

	// Example:
	// {
	//   "message": "Validation Failed",
	//   "errors": {
	//     "email": ["Email is required", "Email must be valid"],
	//     "password": ["Password is required"]
	//   }
	// }
}

type APISuccessResponse struct {
	Message string      `json:"messsage"`
	Data    interface{} `json:"data,omitempty"`
}

// Request DTO (Data Transfer Object) : bind + validation
// required, email, min, max, len vb. business kuralları
// Tip uyumsuzlukları (int beklerken string geldiysa), bind parsing hataları
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=100"` // Kullanimi : `` arasina binding kurallari yazilir
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
}

type Config struct {
	Lang string // "tr", "en", "ru"
}

func loadConfig() Config {
	lang := strings.TrimSpace(strings.ToLower(os.Getenv("APP_LANG"))) // "tr", "en", "ru"
	if lang == "" {
		lang = "en"
	}

	switch lang {
	case "tr", "en", "ru":
		// bunlardan biri gelirse, aynen kullan
	default:
		lang = "en"
	}

	return Config{Lang: lang}
}

func initValidator(cfg Config) *validator.Validate {
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
	translator, found = uni.GetTranslator(cfg.Lang)

	if !found {
		translator, _ = uni.GetTranslator("en")
	}

	var err error
	switch cfg.Lang {
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

func createUserHandler(c *gin.Context) {
	var req CreateUserRequest // Request DTO instance

	// ShouldBindJSON: gelen JSON verisini req struct'ina bind eder
	// 1-) json parse eder
	// 2-) struct alanlarina map'ler
	// 3-) binding tag'larina gore validation yapar
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			// ve : validation hatalari icerisinde olur
			// ok : dogrulama basarili mi

			c.JSON(http.StatusUnprocessableEntity, APIErrorResponse{
				Message: "Validation Failed",
				Errors:  mapValidationErrors(ve),
			})
			return
		}

		c.JSON(http.StatusBadRequest, APIErrorResponse{
			Message: "Invalid Request Payload",
		})
		return
	}

	// Kod buraya kadar gelirse, validation basarili demektir. Fakat business kurallari kontrol edilmemistir. Ornegin: DB'ye kayit eklenirken hata olusabilir.
	// Business hata varsayayimi:
	if err := pretendDBInsert(req); err != nil {
		// db insert hatasi, duplicate email hatasi, timeout hatasi vs gibi.
		c.JSON(http.StatusInternalServerError, APIErrorResponse{
			Message: "Internal Server Error",
			//Message: err.Error(),
		})
		return
	}

	// Basarili response
	c.JSON(http.StatusCreated, APISuccessResponse{
		Message: "User Created Successfully",
		Data: gin.H{
			"name":  req.Name,
			"email": req.Email,
			"age":   req.Age,
		},
	})
}

func mapValidationErrors(ve validator.ValidationErrors) map[string][]string {
	out := make(map[string][]string) // ram de map olusturuldu

	for _, fe := range ve {
		field := fe.Field() // Name, Email, Age
		msg := fe.Translate(translator)

		out[field] = append(out[field], msg)
	}
	return out
}

func pretendDBInsert(req CreateUserRequest) error {
	// Simule edilen DB insert islemi
	if len(req.Email) >= 5 && req.Email[:5] == "fail@" {
		return errors.New("simulated database error")
	}

	return nil
}

package main

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
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
				Errors:  mapValidationErrors(CreateUserRequest{}, ve),
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

func mapValidationErrors(req any, ve validator.ValidationErrors) map[string][]string {
	out := make(map[string][]string) // ram de map olusturuldu

	reqType := reflect.TypeOf(req)

	for _, fe := range ve {
		// fe.Field(): hatali alani verir (struct field ismi, örn: "Name, Email, Age")

		// fe.StructField(): direct struct field bilgisini verir. Direkt struct in adina erisir. (orn: Name, Email, Age)
		field := toJSONFieldName(reqType, fe.StructField())

		out[field] = append(out[field], validationMessage(fe))
	}
	return out
}

func toJSONFieldName(reqType reflect.Type, structField string) string {
	// reflect: Go dilinde runtime'da tip bilgilerine erismek icin kullanilir
	// reqType pointer ise ornek: &CreateUserRequest{}

	// fe.Kind(): hatali alandaki degerin tipini verir (örn: string, int, struct, ptr vb)
	if reqType.Kind() == reflect.Ptr {
		// Elem(): pointer'in gosterdigi tip bilgisini verir
		reqType = reqType.Elem()
	}

	if reqType.Kind() == reflect.Struct {
		return lowerFirst(structField)
	}

	// reqType.FieldByName: struct field bilgisini verir
	f, ok := reqType.FieldByName(structField)
	if !ok {
		return lowerFirst(structField)
	}

	// json tag'ini al
	tag := f.Tag.Get("json")
	if tag == "" { // json tag yoksa
		return lowerFirst(structField)
	}

	// json tag varsa, virgule kadar olan kismi al
	jsonName := strings.Split(tag, ",")[0]
	if jsonName == "-" || jsonName == "" {
		return lowerFirst(structField)
	}

	return jsonName
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func validationMessage(fe validator.FieldError) string {
	// fe.Tag(): hatali alanda hangi kuralin ihlal edildigini verir (örn: "required", "email", "min" vb)
	// fe.Param(): kural parametresini verir (örn: min=2 ise "2" degeri)

	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Minimum value is " + fe.Param()
	case "max":
		return "Maximum value is " + fe.Param()
	default:
		return "Invalid value"
	}
}

func pretendDBInsert(req CreateUserRequest) error {
	// Simule edilen DB insert islemi
	if len(req.Email) >= 5 && req.Email[:5] == "fail@" {
		return errors.New("simulated database error")
	}

	return nil
}

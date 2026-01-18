package user

import (
	"errors"
	"feature-base-starter-kit/pkg/api"
	"feature-base-starter-kit/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func CreateUserHandler(c *gin.Context) {
	var req CreateUserRequest // Request DTO instance

	// ShouldBindJSON: gelen JSON verisini req struct'ina bind eder
	// 1-) json parse eder
	// 2-) struct alanlarina map'ler
	// 3-) binding tag'larina gore validation yapar
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			// ve : validation hatalari icerisinde olur
			// ok : dogrulama basarili mi

			c.JSON(http.StatusUnprocessableEntity, api.APIErrorResponse{
				Message: "Validation Failed",
				Errors:  validation.MapValidationErrors(ve),
			})
			return
		}

		c.JSON(http.StatusBadRequest, api.APIErrorResponse{
			Message: "Invalid Request Payload",
		})
		return
	}

	// Kod buraya kadar gelirse, validation basarili demektir. Fakat business kurallari kontrol edilmemistir. Ornegin: DB'ye kayit eklenirken hata olusabilir.
	// Business hata varsayayimi:
	if err := pretendDBInsert(req); err != nil {
		// db insert hatasi, duplicate email hatasi, timeout hatasi vs gibi.
		c.JSON(http.StatusInternalServerError, api.APIErrorResponse{
			Message: "Internal Server Error",
			//Message: err.Error(),
		})
		return
	}

	format := c.Query("format")
	response := api.APISuccessResponse{
		Message: "User Created Successfully",
		Data: gin.H{
			"name":  req.Name,
			"email": req.Email,
			"age":   req.Age,
		},
	}

	// Basarili response
	switch format {
	case "xml":
		c.XML(http.StatusCreated, response)
	case "yaml", "yml":
		c.YAML(http.StatusCreated, response)
	default:
		c.JSON(http.StatusCreated, response)
	}

}

func pretendDBInsert(req CreateUserRequest) error {
	// Simule edilen DB insert islemi
	if len(req.Email) >= 5 && req.Email[:5] == "fail@" {
		return errors.New("simulated database error")
	}

	return nil
}

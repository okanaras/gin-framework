package api

import "github.com/gin-gonic/gin"

type APIErrorResponse struct {
	Message string              `json:"message" xml:"message" yaml:"message"`                            // struct tag
	Errors  map[string][]string `json:"errors,omitempty" xml:"errors,omitempty" yaml:"errors,omitempty"` // omitempty: alan boşsa JSON, XML veya YAML çıktısında göstermez

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
	Message string      `json:"message" xml:"message" yaml:"message"`
	Data    interface{} `json:"data,omitempty" xml:"data,omitempty" yaml:"data,omitempty"`
}

func SendError(ctx *gin.Context, status int, message string, errs map[string][]string) {
	ctx.JSON(status, APIErrorResponse{
		Message: message,
		Errors:  errs,
	})
}

func SendSuccess(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(status, APISuccessResponse{
		Message: message,
		Data:    data,
	})
}

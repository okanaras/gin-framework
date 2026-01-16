package router

import (
	"feature-base-starter-kit/internal/modules/user"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()

	r.POST("/users", user.CreateUserHandler)

	return r
}

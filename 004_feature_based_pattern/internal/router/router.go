package router

import (
	"feature-base-starter-kit/internal/config"
	"feature-base-starter-kit/internal/middleware"
	"feature-base-starter-kit/internal/modules/user"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	// r.Use(mid1, mid2) // Global Middleware eklenebilir
	r.Use(middleware.LoggerMiddleware())

	protectedRoute := r.Group("/api")
	protectedRoute.Use(middleware.AuthMiddleware(cfg.API_SECRET_KEY))
	protectedRoute.POST("/users", user.CreateUserHandler)

	//r.POST("/users", user.CreateUserHandler)

	return r
}

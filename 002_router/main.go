package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Query Param : ?active=true&role=admin
	// Query Param zorunlu değildir. Eğer query parametre gönderilmezse boş string döner.
	router.GET("/users", func(ctx *gin.Context) {
		active := ctx.Query("active")
		role := ctx.Query("role")

		ctx.JSON(http.StatusOK, gin.H{
			"endpoint": ctx.FullPath(),
			"method":   ctx.Request.Method,
			"active":   active,
			"role":     role,
			"query":    ctx.Request.URL.RawQuery,
			"message":  "List of users (Query Param example)",
		})
	})

	// Path Param : /users/:id (örnek: /users/123)
	// Path Param genellikle zorunludur.
	router.GET("/users/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		ctx.JSON(http.StatusOK, gin.H{
			"endpoint": ctx.FullPath(),
			"method":   ctx.Request.Method,
			"id":       id,
			"message":  "User details (Path Param example)",
		})
	})

	// Path Param ve Query Param birlikte kullanımı
	// Örnek: /users/123/profile?is_active=true
	router.GET("/users/:id/profile", func(ctx *gin.Context) {
		id := ctx.Param("id")
		isActive := ctx.Query("is_active")

		ctx.JSON(http.StatusOK, gin.H{
			"endpoint":  ctx.FullPath(),
			"method":    ctx.Request.Method,
			"id":        id,
			"is_active": isActive,
			"query":     ctx.Request.URL.RawQuery,
			"message":   "User profile (Path and Query Param example)",
		})
	})

	router.Run(":8080")
}

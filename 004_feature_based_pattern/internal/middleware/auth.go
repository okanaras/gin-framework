package middleware

import (
	"feature-base-starter-kit/pkg/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(validSecretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")

		if apiKey != validSecretKey {
			// Abort : Request zincirini durdurur ve belirtilen yanıtı gönderir.
			// Yalnizca return vermek, zinciri durdurmaz. Bu nedenle Abort kullanilir.

			api.SendError(ctx, http.StatusUnauthorized, "Unauthorized access", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

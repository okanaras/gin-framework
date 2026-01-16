package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path

		log.Printf("[GIRIS] %s istegi %s adresinden %s yoluna geldi.", clientIP, method, path)

		// İşlem öncesi
		// Zincirde sonraki middleware veya handler'a geçiş
		ctx.Next()
		// İşlem sonrası

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		statusCode := ctx.Writer.Status()

		log.Printf("[CIKIS] Durum: %d | Sure: %v | IP: %s | YOL: %s | METHOD: %s\n", statusCode, latency, clientIP, path, method)
	}
}

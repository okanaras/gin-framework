package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Alternatif olarak, middleware'leri manuel olarak ekleyebilirsiniz. Default() fonksiyonu ile aynı işlevselliği sağlar.
	//router2 := gin.New()
	//router2.Use(gin.Logger())   // istekleri loglamak için middleware ekleniyor
	//router2.Use(gin.Recovery()) // panik durumunda 500 hatası döndürmek için middleware ekleniyor

	router.GET("/", func(ctx *gin.Context) {
		// Gelen istekler body, header, query parametreleri ile ctx üzerinden erişilebilir

		// Content-Type otomatik olarak application/json olarak ayarlanır
		// gin.H, map[string]interface{} türünün kısaltmasıdır
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Merhaba, Gin Framework!",
		})
	})

	router.GET("/panic", func(ctx *gin.Context) {
		// Bu handler bir panik oluşturur
		panic("Bu bir paniktir!, Recover middleware tarafından yakalanacaktır.")
	})

	//port := ":8080" // Varsayılan port
	port := os.Getenv("PORT")
	// Isletim sistemi env değişkenlerinden PORT'u alir.
	// Eğer PORT env değişkeni ayarlanmamışsa boş string döner.
	// macOS/Linux için terminalde: "PORT=9090 go run ." şeklinde ayarlayabilirsiniz.

	if port == "" {
		port = "8080" // Eğer env değişkeni yoksa varsayılan portu kullan
	}

	//addr := ":" + port
	// Sprintf ile de aynı işlemi yapabiliriz
	addr := fmt.Sprintf(":%s", port)

	if err := router.Run(addr); err != nil {
		fmt.Printf("Sunucu başlatılamadı: %v\n", err)
		panic(err)
	}
}

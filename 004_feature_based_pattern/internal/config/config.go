package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Lang           string // "tr", "en", "ru"
	API_SECRET_KEY string
	Port           string
}

func LoadConfig() Config {

	if err := godotenv.Load(); err != nil {
		// .env dosyasi bulunamazsa, hata vermez, ortam degiskenlerinden devam eder
		log.Println("No .env file found, continuing with environment variables")
	}

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

	apiSecretKey := os.Getenv("API_SECRET_KEY")
	if apiSecretKey == "" {
		log.Fatal("API_SECRET_KEY is not set in environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Lang:           lang,
		API_SECRET_KEY: apiSecretKey,
		Port:           port,
	}
}

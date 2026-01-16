package config

import (
	"os"
	"strings"
)

type Config struct {
	Lang string // "tr", "en", "ru"
}

func LoadConfig() Config {
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

	return Config{Lang: lang}
}

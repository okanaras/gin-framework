package main

import (
	"feature-base-starter-kit/internal/config"
	"feature-base-starter-kit/internal/router"
	"feature-base-starter-kit/pkg/validation"
)

func main() {
	cfg := config.LoadConfig()

	validation.Init(cfg.Lang)

	r := router.Setup()
	r.Run(":8080")
}

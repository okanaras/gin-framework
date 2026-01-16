package main

import (
	"feature-base-starter-kit/internal/config"
	"feature-base-starter-kit/internal/router"
	"feature-base-starter-kit/pkg/validation"
)

func main() {
	cfg := config.LoadConfig()

	validation.Init(cfg.Lang)

	r := router.Setup(&cfg)
	// pointer olarak gonderdik cunku config yapisi buyuk olabilir. Yani cfg.Lang gibi kullanmak yerine, pointer ile gonderip, icinde istedigimiz yere erisebiliriz.

	r.Run(":" + cfg.Port)
}

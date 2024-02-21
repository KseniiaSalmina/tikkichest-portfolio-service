package main

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"

	app "github.com/KseniiaSalmina/tikkichest-portfolio-service/internal"
	"github.com/KseniiaSalmina/tikkichest-portfolio-service/internal/config"
)

var (
	cfg config.Application
)

func init() {
	_ = godotenv.Load(".env")
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
}

func main() {
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	application.Run()
}

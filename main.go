package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	_ "github.com/KseniiaSalmina/tikkichest-portfolio-service/docs"
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

// @title Tikkichest portfolio service
// @version 1.1.0
// @description part of tikkichest
// @host localhost:8088
// @BasePath /
func main() {
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}
	application.Run()
}

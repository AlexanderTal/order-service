package config

import (
	"log"

	"github.com/AlexanderTal/order-service/internal/app/config/section"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Root Config

type Config struct {
	Repository section.Repository
	Processor  section.Processor
	Monitor    section.Monitor
}

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	if err := envconfig.Process("APP", &Root); err != nil {
		log.Fatalf("failed to process env config: %v", err)
	}
}

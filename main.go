package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AlexanderTal/order-service/internal/app/config"
	rhealth "github.com/AlexanderTal/order-service/internal/app/handler/http/health"
	rprocessor "github.com/AlexanderTal/order-service/internal/app/processor/http"
	rcpostgres "github.com/AlexanderTal/order-service/internal/app/repository/conn/postgres"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	config.Load()
	cfg := config.Root

	db, err := rcpostgres.NewClient(context.Background(), cfg.Repository.Postgres)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close database connection: %v", err)
		}
	}()

	log.Printf("Connected to database: %s", cfg.Repository.Postgres.Name)

	hHealth := rhealth.NewHandler()
	httpProcessor := rprocessor.NewHTTP(hHealth, cfg.Processor.WebServer)

	if err := httpProcessor.Serve(); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}

	return nil
}

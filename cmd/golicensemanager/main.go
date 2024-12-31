package main

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/LywwKkA-aD/golicensemanager/internal/app"
	"github.com/LywwKkA-aD/golicensemanager/internal/config"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		sugar.Warnf("Error loading .env file: %v", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		sugar.Fatalf("Failed to load configuration: %v", err)
	}

	// Create application instance
	application, err := app.NewApp(cfg, sugar)
	if err != nil {
		sugar.Fatalf("Failed to create application: %v", err)
	}

	// Start the application
	if err := application.Start(); err != nil {
		sugar.Fatalf("Failed to start application: %v", err)
	}
}

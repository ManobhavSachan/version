package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"version-backend/internal/api"
	"version-backend/internal/config"
	"version-backend/internal/db"
	"version-backend/internal/osquery"
	"version-backend/pkg/logger"
)

func main() {
	// Initialize logger
	logger.SetLevel("info")
	log := logger.GetLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	database, err := db.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize osquery client
	osqueryClient, err := osquery.NewClient(cfg.Osquery.SocketPath)
	if err != nil {
		log.Fatalf("Failed to create osquery client: %v", err)
	}
	defer osqueryClient.Close()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Collect initial data at startup
	log.Info("Collecting initial system information...")
	if err := collectAndSaveData(ctx, osqueryClient, database); err != nil {
		log.Errorf("Failed to collect initial data: %v", err)
	}

	// Start data collection in background
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.Osquery.QueryInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := collectAndSaveData(ctx, osqueryClient, database); err != nil {
					log.Errorf("Failed to collect and save data: %v", err)
				}
			}
		}
	}()

	// Initialize and start HTTP server
	router := api.NewRouter(database)
	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info("Shutting down server...")
		cancel()
	}()

	log.Infof("Server starting on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// collectAndSaveData collects system information and saves it to the database
func collectAndSaveData(ctx context.Context, client *osquery.Client, db *db.DB) error {
	// Collect system information
	sysInfo, err := client.GetSystemInfo(ctx)
	if err != nil {
		return fmt.Errorf("failed to get system info: %w", err)
	}

	// Save to database
	if err := db.SaveSystemInfo(sysInfo); err != nil {
		return fmt.Errorf("failed to save system info: %w", err)
	}

	return nil
}

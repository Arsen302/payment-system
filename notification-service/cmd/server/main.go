package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arsen302/payment-system/notification-service/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "NOTIFICATION-SERVICE: ", log.LstdFlags)
	logger.Println("Starting Notification Service...")

	// TODO: Initialize Kafka consumer
	// TODO: Set up event handlers

	logger.Printf("Notification service started with Kafka bootstrap servers: %s", cfg.KafkaBootstrapServers)

	// Set up graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-sigCh
	logger.Println("Shutting down service...")

	// TODO: Properly close all connections and clean up resources

	logger.Println("Service stopped")
} 
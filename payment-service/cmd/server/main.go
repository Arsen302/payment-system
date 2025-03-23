package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Arsen302/payment-system/payment-service/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "PAYMENT-SERVICE: ", log.LstdFlags)
	logger.Println("Starting Payment Service...")

	// TODO: Initialize database connection (PostgreSQL)
	// TODO: Initialize Kafka producer
	// TODO: Initialize gRPC server and register services

	// Create a TCP listener on the configured port
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}
	
	// Note: lis will be used when we implement the gRPC server
	_ = lis

	logger.Printf("Server started at port %d", cfg.Port)

	// Set up graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-sigCh
	logger.Println("Shutting down server...")

	// TODO: Properly close all connections and clean up resources

	logger.Println("Server stopped")
} 
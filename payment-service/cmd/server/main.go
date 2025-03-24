package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Arsen302/payment-system/payment-service/internal/config"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Define Prometheus metrics
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_request_count_total",
			Help: "Total number of payment requests processed",
		},
		[]string{"status"},
	)
	
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_request_duration_seconds",
			Help:    "Histogram of payment request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Create router
	router := mux.NewRouter()
	
	// Setup API routes
	router.HandleFunc("/api/payments", createPaymentHandler).Methods("POST")
	router.HandleFunc("/api/payments/{id}", getPaymentHandler).Methods("GET")
	router.HandleFunc("/api/payments", listPaymentsHandler).Methods("GET")
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start metrics server if enabled
	if cfg.Metrics.Enabled {
		go startMetricsServer(cfg.Metrics.Port)
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	log.Printf("Payment service started on port %s", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

// Start a separate HTTP server for Prometheus metrics
func startMetricsServer(port string) {
	metricsRouter := mux.NewRouter()
	metricsRouter.Handle("/metrics", promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: metricsRouter,
	}

	log.Printf("Metrics server started on port %s", port)

	if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Metrics server failed: %v", err)
	}
}

// Example handlers
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK")
}

func createPaymentHandler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.With(prometheus.Labels{"endpoint": "/api/payments"}))
	defer timer.ObserveDuration()

	// Process payment...
	// In a real implementation, you would validate the request and create the payment

	// Record success metric
	requestCounter.With(prometheus.Labels{"status": "success"}).Inc()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id": "payment123", "status": "processing"}`)
}

func getPaymentHandler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.With(prometheus.Labels{"endpoint": "/api/payments/{id}"}))
	defer timer.ObserveDuration()

	vars := mux.Vars(r)
	id := vars["id"]

	// In a real implementation, you would fetch the payment from database

	// Record success metric
	requestCounter.With(prometheus.Labels{"status": "success"}).Inc()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"id": "%s", "status": "completed"}`, id)
}

func listPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	timer := prometheus.NewTimer(requestDuration.With(prometheus.Labels{"endpoint": "/api/payments"}))
	defer timer.ObserveDuration()

	// In a real implementation, you would list payments from database

	// Record success metric
	requestCounter.With(prometheus.Labels{"status": "success"}).Inc()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"payments": [{"id": "payment123", "status": "completed"}]}`)
} 
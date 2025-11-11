package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahernandez9/rockets/internal/api"
	"github.com/ahernandez9/rockets/internal/pubsub"
	"github.com/ahernandez9/rockets/internal/repository/inmemory"
	"github.com/ahernandez9/rockets/internal/service"
)

// @title Rockets API
// @version 1.0
// @description REST API for rocket telemetry system with async message processing via pub/sub channels
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@rockets.example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8088
// @BasePath /
// @schemes http

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}

	// Initialize repository (in-memory storage with mutex protection)
	repo := inmemory.NewInMemoryRepository()
	log.Println("Repository initialized (in-memory)")

	// Initialize pub/sub (channel-based with buffer)
	ps := pubsub.NewChannelPubSub(100) // Buffer size of 100 messages
	log.Println("Pub/Sub initialized (channel-based, buffer=100)")

	// Create rocket service with async message processing
	rocketService := service.NewRocketService(repo, ps)
	log.Println("Rocket service started with async message processor")

	// Setup HTTP router
	router := api.SetupRouter(rocketService)

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in goroutine
	addr := fmt.Sprintf(":%s", port)
	go func() {
		log.Printf("Starting Rockets API server on %s", addr)
		log.Printf("Swagger documentation available at http://localhost:%s/swagger/index.html", port)
		log.Println("Architecture: HTTP -> Handler -> Pub/Sub (channel) -> Async Processor -> Repository")

		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	rocketService.Stop()
	log.Println("Server stopped gracefully")
}

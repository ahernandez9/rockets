package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ahernandez9/rockets/internal/api"
	"github.com/ahernandez9/rockets/internal/pubsub/channel"
	"github.com/ahernandez9/rockets/internal/repository/inmemory"
	"github.com/ahernandez9/rockets/internal/service"
)

// @title Rockets API
// @version 1.0
// @description REST API for rocket system with message processing

func main() {
	port := os.Getenv("PORT") // We could use a more advanced approach to load env vars, ex: viper
	if port == "" {
		port = "8088"
	}

	// initialize observability here (logging, tracing, metrics)

	// Dependencies
	repo := inmemory.NewInMemoryRepository()
	pubsub := channel.NewPubSub(1000)

	// Services
	rocketService := service.NewRocketService(repo)
	messageService := service.NewMessageService(pubsub, repo)

	router := api.SetupRouter(messageService, rocketService)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start async message processor
	go messageService.Start()
	defer messageService.Stop()

	// Start HTTP server
	addr := fmt.Sprintf(":%s", port)
	go func() {
		log.Printf("Starting Rockets API server on %s", addr)

		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Server stopped")
}

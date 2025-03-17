package main

import (
	"log"
	"net/http"
	"time"

	"money-transfer-system/api"
	"money-transfer-system/service"
	"money-transfer-system/store"
)

func main() {
	// Create the account store and initialize with default accounts
	accountStore := store.NewInMemoryStore()
	accountStore.Setup()

	// Create services
	transferService := service.NewTransferService(accountStore)

	// Create API and set up routes
	apiHandler := api.NewAPI(transferService, accountStore)
	router := apiHandler.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         ":8081",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Start the server
	log.Println("Money Transfer System starting on port 8081...")
	log.Fatal(server.ListenAndServe())
} 
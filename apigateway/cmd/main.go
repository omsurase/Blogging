package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/omsurase/Blogging/apigateway/internal/config"
	"github.com/omsurase/Blogging/apigateway/internal/handlers"
	"github.com/omsurase/Blogging/apigateway/internal/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize router
	r := mux.NewRouter()

	// Initialize proxy service
	proxyService := service.NewProxyService()

	// Initialize handlers
	gatewayHandler := handlers.NewGatewayHandler(proxyService)

	// Define routes
	r.PathPrefix("/auth").HandlerFunc(gatewayHandler.AuthHandler)
	r.PathPrefix("/posts").HandlerFunc(gatewayHandler.PostsHandler)
	r.PathPrefix("/users").HandlerFunc(gatewayHandler.UsersHandler)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting API Gateway on %s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

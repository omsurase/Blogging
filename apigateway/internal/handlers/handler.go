package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/omsurase/Blogging/apigateway/internal/service"
)

type GatewayHandler struct {
	proxyService *service.ProxyService
}

func NewGatewayHandler(proxyService *service.ProxyService) *GatewayHandler {
	return &GatewayHandler{
		proxyService: proxyService,
	}
}

func (h *GatewayHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Gateway: Received request for path: %s", r.URL.Path)
	authEndpoint := strings.TrimPrefix(r.URL.Path, "/auth")
	targetURL := "http://localhost:8081/auth" + authEndpoint
	log.Printf("API Gateway: Forwarding request to %s", targetURL)
	h.proxyService.ProxyRequest(targetURL, w, r)
}

func (h *GatewayHandler) PostsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("API Gateway: Received request for path: %s", r.URL.Path)
	postsEndpoint := r.URL.Path
	targetURL := "http://localhost:8082" + postsEndpoint
	log.Printf("API Gateway: Forwarding request to %s", targetURL)
	h.proxyService.ProxyRequest(targetURL, w, r)
}

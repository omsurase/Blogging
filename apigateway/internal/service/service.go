package service

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ProxyService struct{}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

func (s *ProxyService) ProxyRequest(target string, w http.ResponseWriter, r *http.Request) {
	url, err := url.Parse(target)
	if err != nil {
		log.Printf("Error parsing target URL: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Proxying request to: %s", target)

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Path = r.URL.Path
		req.URL.RawQuery = r.URL.RawQuery
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Host = url.Host
	}

	proxy.ServeHTTP(w, r)

	log.Printf("Request proxied successfully")
}

package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/omsurase/Blogging/blog-service/internal/models"
	authpb "github.com/omsurase/Blogging/blog-service/internal/pb"
	"github.com/omsurase/Blogging/blog-service/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type BlogHandler struct {
	service    *service.BlogService
	authClient authpb.AuthServiceClient
}

func NewBlogHandler(service *service.BlogService, authConn *grpc.ClientConn) *BlogHandler {
	log.Println("Initializing BlogHandler")
	return &BlogHandler{
		service:    service,
		authClient: authpb.NewAuthServiceClient(authConn),
	}
}

func (h *BlogHandler) ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("ValidateToken: Starting token validation")

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			log.Println("ValidateToken: No Authorization header provided")
			http.Error(w, "No Authorization header provided", http.StatusUnauthorized)
			return
		}

		// Split the header and take the second part (the actual token)
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			log.Println("ValidateToken: Invalid Authorization header format")
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]
		log.Printf("ValidateToken: Extracted token: %s\n", token)

		response, err := h.authClient.ValidateToken(r.Context(), &authpb.ValidateTokenRequest{Token: token})

		if err != nil {
			log.Printf("ValidateToken: Error validating token: %v\n", err)
			http.Error(w, "Error validating token", http.StatusInternalServerError)
			return
		}

		if !response.GetValid() {
			log.Println("ValidateToken: Invalid token")
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		log.Println("ValidateToken: Token validated successfully")
		next.ServeHTTP(w, r)
	}
}

func (h *BlogHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	log.Println("CreatePost: Starting post creation")
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Printf("CreatePost: Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreatePost(&post); err != nil {
		log.Printf("CreatePost: Error creating post: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("CreatePost: Post created successfully with ID: %s\n", post.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *BlogHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAllPosts: Fetching all posts")
	posts, err := h.service.GetAllPosts()
	if err != nil {
		log.Printf("GetAllPosts: Error fetching posts: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GetAllPosts: Successfully fetched %d posts\n", len(posts))
	json.NewEncoder(w).Encode(posts)
}

func (h *BlogHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("GetPost: Fetching post with ID: %s\n", id)

	post, err := h.service.GetPost(id)
	if err != nil {
		log.Printf("GetPost: Error fetching post: %v\n", err)
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("GetPost: Successfully fetched post with ID: %s\n", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *BlogHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("UpdatePost: Starting update for post with ID: %s\n", id)

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		log.Printf("UpdatePost: Error decoding request body: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePost(id, &post); err != nil {
		log.Printf("UpdatePost: Error updating post: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("UpdatePost: Successfully updated post with ID: %s\n", id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

func (h *BlogHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("DeletePost: Attempting to delete post with ID: %s\n", id)

	if err := h.service.DeletePost(id); err != nil {
		log.Printf("DeletePost: Error deleting post: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("DeletePost: Successfully deleted post with ID: %s\n", id)
	w.WriteHeader(http.StatusNoContent)
}

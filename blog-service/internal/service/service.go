package service

import (
	"github.com/omsurase/Blogging/blog-service/internal/models"
	"github.com/omsurase/Blogging/blog-service/internal/repository"
)

type BlogService struct {
	repo *repository.MongoRepository
}

func NewBlogService(repo *repository.MongoRepository) *BlogService {
	return &BlogService{
		repo: repo,
	}
}

func (s *BlogService) CreatePost(post *models.Post) error {
	return s.repo.CreatePost(post)
}

func (s *BlogService) GetAllPosts() ([]models.Post, error) {
	return s.repo.GetAllPosts()
}

func (s *BlogService) GetPost(id string) (*models.Post, error) {
	return s.repo.GetPost(id)
}

func (s *BlogService) UpdatePost(id string, post *models.Post) error {
	return s.repo.UpdatePost(id, post)
}

func (s *BlogService) DeletePost(id string) error {
	return s.repo.DeletePost(id)
}

package service

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/omsurase/Blogging/auth-service/internal/config"
	"github.com/omsurase/Blogging/auth-service/internal/models"
	"github.com/omsurase/Blogging/auth-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, config: cfg}
}

func (s *AuthService) Register(user *models.User) (string, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	log.Printf("%s", hashedPassword)
	user.Password = string(hashedPassword)

	err = s.repo.Create(user)
	if err != nil {
		return "", "", err
	}

	// Generate token after successful registration
	token, err := s.generateToken(user.Username)
	if err != nil {
		return "", "", err
	}

	return user.ID, token, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.repo.GetByUsername(username)
	//log.Printf("request recieved.3")
	if err != nil {
		//log.Printf("request recieved.4")
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return s.generateToken(username)
}

func (s *AuthService) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, nil
}

func (s *AuthService) generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(s.config.TokenExpiryHours)).Unix(),
	})

	return token.SignedString([]byte(s.config.JWTSecret))
}

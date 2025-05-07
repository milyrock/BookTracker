package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"booktracker/internal/models"
	"booktracker/internal/repository"
)

type AuthService struct {
	Repo      *repository.PostgresRepo
	JWTSecret string
}

func NewAuthService(r *repository.PostgresRepo, secret string) *AuthService {
	return &AuthService{Repo: r, JWTSecret: secret}
}

func (s *AuthService) Register(u *models.User) error {
	existing, err := s.Repo.GetUserByUsername(u.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("username already exists")
	}
	_, err = s.Repo.CreateUser(u)
	return err
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}
	if user.PasswordHash != password { // TODO: proper password hashing
		return "", errors.New("invalid password")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

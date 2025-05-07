package service

import (
	"booktracker/internal/models"
	"booktracker/internal/repository"
)

type BookService struct {
	Repo *repository.PostgresRepo
}

func NewBookService(r *repository.PostgresRepo) *BookService {
	return &BookService{Repo: r}
}

func (s *BookService) CreateBook(b *models.Book) (int, error) {
	return s.Repo.CreateBook(b)
}

func (s *BookService) ListBooks(userID int) ([]*models.Book, error) {
	return s.Repo.ListBooks(userID)
}

func (s *BookService) GetBook(id, userID int) (*models.Book, error) {
	return s.Repo.GetBook(id, userID)
}

func (s *BookService) UpdateBook(b *models.Book) error {
	return s.Repo.UpdateBook(b)
}

func (s *BookService) DeleteBook(id, userID int) error {
	return s.Repo.DeleteBook(id, userID)
}

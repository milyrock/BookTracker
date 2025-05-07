package repository

import (
	"database/sql"
	"errors"

	"booktracker/internal/models"
)

const (
	createUserQuery = `
		INSERT INTO users(username, password_hash) 
		VALUES($1, $2) 
		RETURNING id`

	getUserByUsernameQuery = `
		SELECT id, username, password_hash 
		FROM users 
		WHERE username = $1`

	createBookQuery = `
		INSERT INTO books(user_id, title, author, status, finished_at) 
		VALUES($1, $2, $3, $4, $5) 
		RETURNING id`

	listBooksQuery = `
		SELECT id, user_id, title, author, status, finished_at 
		FROM books 
		WHERE user_id = $1 
		ORDER BY id DESC`

	getBookQuery = `
		SELECT id, user_id, title, author, status, finished_at 
		FROM books 
		WHERE id = $1 AND user_id = $2`

	updateBookQuery = `
		UPDATE books 
		SET title = $1, author = $2, status = $3, finished_at = $4 
		WHERE id = $5 AND user_id = $6`

	deleteBookQuery = `
		DELETE FROM books 
		WHERE id = $1 AND user_id = $2`
)

type PostgresRepo struct {
	DB *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{DB: db}
}

// CreateUser сохраняет пользователя и возвращает его ID
func (r *PostgresRepo) CreateUser(u *models.User) (int, error) {
	var id int
	err := r.DB.QueryRow(createUserQuery, u.Username, u.PasswordHash).Scan(&id)
	return id, err
}

// GetUserByUsername возвращает пользователя по имени
func (r *PostgresRepo) GetUserByUsername(username string) (*models.User, error) {
	u := &models.User{}
	err := r.DB.QueryRow(getUserByUsernameQuery, username).Scan(&u.ID, &u.Username, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

// CreateBook сохраняет книгу и возвращает её ID
func (r *PostgresRepo) CreateBook(b *models.Book) (int, error) {
	var id int
	err := r.DB.QueryRow(createBookQuery, b.UserID, b.Title, b.Author, b.Status, b.FinishedAt).Scan(&id)
	return id, err
}

// ListBooks возвращает список книг для указанного пользователя
func (r *PostgresRepo) ListBooks(userID int) ([]*models.Book, error) {
	rows, err := r.DB.Query(listBooksQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		b := &models.Book{}
		err := rows.Scan(&b.ID, &b.UserID, &b.Title, &b.Author, &b.Status, &b.FinishedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

// GetBook возвращает книгу по ID
func (r *PostgresRepo) GetBook(id, userID int) (*models.Book, error) {
	b := &models.Book{}
	err := r.DB.QueryRow(getBookQuery, id, userID).Scan(&b.ID, &b.UserID, &b.Title, &b.Author, &b.Status, &b.FinishedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return b, err
}

// UpdateBook обновляет информацию о книге
func (r *PostgresRepo) UpdateBook(b *models.Book) error {
	_, err := r.DB.Exec(updateBookQuery, b.Title, b.Author, b.Status, b.FinishedAt, b.ID, b.UserID)
	return err
}

// DeleteBook удаляет книгу
func (r *PostgresRepo) DeleteBook(id, userID int) error {
	_, err := r.DB.Exec(deleteBookQuery, id, userID)
	return err
}

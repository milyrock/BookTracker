package models

type Book struct {
	ID         int     `db:"id"`
	UserID     int     `db:"user_id"`
	Title      string  `db:"title"`
	Author     *string `db:"author"`
	Status     string  `db:"status"`
	FinishedAt *string `db:"finished_at"`
}

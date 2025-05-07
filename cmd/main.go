package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"booktracker/internal/middleware"
	"booktracker/internal/repository"
	"booktracker/internal/service"
	v1 "booktracker/internal/v1"
)

func main() {
	e := echo.New()
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())

	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		e.Logger.Fatal(err)
	}

	m, err := migrate.New("file://db", dbURL)
	if err != nil {
		e.Logger.Fatalf("Migration init failed: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		e.Logger.Fatalf("Could not run migrations: %v", err)
	}

	repo := repository.NewPostgresRepo(db)

	authSvc := service.NewAuthService(repo, os.Getenv("JWT_SECRET"))
	bookSvc := service.NewBookService(repo)

	authHandler := v1.NewAuthHandler(authSvc)
	bookHandler := v1.NewBookHandler(bookSvc)

	api := e.Group("/api")
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	books := api.Group("/books")
	books.Use(middleware.JWTMiddleware())
	books.GET("", bookHandler.List)
	books.GET(":id", bookHandler.Get)
	books.POST("", bookHandler.Create)
	books.PUT(":id", bookHandler.Update)
	books.DELETE(":id", bookHandler.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}

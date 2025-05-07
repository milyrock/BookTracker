package v1

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"booktracker/internal/models"
	"booktracker/internal/service"
)

type BookHandler struct {
	Service *service.BookService
}

func NewBookHandler(s *service.BookService) *BookHandler {
	return &BookHandler{Service: s}
}

func (h *BookHandler) List(c echo.Context) error {
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	books, err := h.Service.ListBooks(userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, books)
}

func (h *BookHandler) Get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book id")
	}
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	book, err := h.Service.GetBook(id, userID)
	if err != nil {
		return err
	}
	if book == nil {
		return echo.NewHTTPError(http.StatusNotFound, "book not found")
	}
	return c.JSON(http.StatusOK, book)
}

func (h *BookHandler) Create(c echo.Context) error {
	var book models.Book
	if err := c.Bind(&book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	book.UserID = userID
	id, err := h.Service.CreateBook(&book)
	if err != nil {
		return err
	}
	book.ID = id
	return c.JSON(http.StatusCreated, book)
}

func (h *BookHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book id")
	}
	var book models.Book
	if err := c.Bind(&book); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	book.ID = id
	book.UserID = userID
	if err := h.Service.UpdateBook(&book); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *BookHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book id")
	}
	userID := int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(float64))
	if err := h.Service.DeleteBook(id, userID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

package v1

import (
	"net/http"

	"booktracker/internal/models"
	"booktracker/internal/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{Service: s}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(struct{ Username, Password string })
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user := &models.User{
		Username:     req.Username,
		PasswordHash: req.Password, // TODO: proper password hashing
	}
	if err := h.Service.Register(user); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(struct{ Username, Password string })
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	token, err := h.Service.Login(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

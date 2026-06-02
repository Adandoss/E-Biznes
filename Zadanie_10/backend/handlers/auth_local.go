package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"golang.org/x/crypto/bcrypt"

	"sklep/models"
)

// POST /auth/register
func (h *AuthHandler) Register(c *echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawne dane"})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email i hasło są wymagane"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd serwera"})
	}

	user := models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Provider: "local",
	}

	if err := h.DB.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return c.JSON(http.StatusConflict, map[string]string{"error": "Użytkownik o tym emailu już istnieje"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się utworzyć użytkownika"})
	}

	tokenString, err := generateJWT(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd generowania tokenu"})
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"token": tokenString,
		"user": map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"name":     user.Name,
			"provider": user.Provider,
		},
	})
}

// POST /auth/login
func (h *AuthHandler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawne dane"})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email i hasło są wymagane"})
	}

	var user models.User
	if err := h.DB.Where("email = ? AND provider = ?", req.Email, "local").First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Niepoprawny email lub hasło"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Niepoprawny email lub hasło"})
	}

	tokenString, err := generateJWT(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd generowania tokenu"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"token": tokenString,
		"user": map[string]any{
			"id":       user.ID,
			"email":    user.Email,
			"name":     user.Name,
			"provider": user.Provider,
		},
	})
}

// GET /auth/me
func (h *AuthHandler) Me(c *echo.Context) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Nie zalogowano"})
	}

	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Użytkownik nie znaleziony"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"id":       user.ID,
		"email":    user.Email,
		"name":     user.Name,
		"provider": user.Provider,
	})
}

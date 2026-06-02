package handlers

import (
	"net/http"
	"sklep/models"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	DB *gorm.DB
}

func (h *PaymentHandler) CreatePayment(c *echo.Context) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Nie zalogowano"})
	}

	payment := new(models.Payment)
	if err := c.Bind(payment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Niepoprawne dane płatności"})
	}

	if payment.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Kwota musi być większa od zera"})
	}
	
	payment.UserID = userID
	
	if err := h.DB.Create(payment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Nie udało się utworzyć płatności"})
	}

	return c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) GetPayments(c *echo.Context) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Nie zalogowano"})
	}

	var payments []models.Payment
	if err := h.DB.Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Błąd pobierania płatności"})
	}
	return c.JSON(http.StatusOK, payments)
}

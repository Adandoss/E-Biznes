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
	payment := new(models.Payment)
	if err := c.Bind(payment); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	
	if err := h.DB.Create(payment).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) GetPayments(c *echo.Context) error {
	var payments []models.Payment
	if err := h.DB.Find(&payments).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, payments)
}

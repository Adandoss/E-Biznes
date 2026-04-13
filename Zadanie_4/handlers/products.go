package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"

	"sklep/models"
)

type ProductHandler struct {
	DB *gorm.DB
}

func (h *ProductHandler) CreateProduct(c *echo.Context) error {
	var product models.Product

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawne dane"})
	}

	if err := h.DB.Create(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się utworzyć produktu"})
	}

	return c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProducts(c *echo.Context) error {
	var products []models.Product

	err := h.DB.Scopes(
		models.WithCategory,          
		models.PriceGreaterThan(2000),  
		models.OrderByPriceDesc,      
	).Find(&products).Error

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd podczas pobierania produktów"})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *echo.Context) error {
	id := c.Param("id")
	var product models.Product

	if err := h.DB.Preload("Category").First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Produkt nie znaleziony"})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *echo.Context) error {
	id := c.Param("id")
	var product models.Product

	if err := h.DB.First(&product, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Produkt nie znaleziony"})
	}

	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawne dane"})
	}

	h.DB.Save(&product)
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *echo.Context) error {
	id := c.Param("id")

	if err := h.DB.Delete(&models.Product{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się usunąć produktu"})
	}

	return c.NoContent(http.StatusNoContent)
}
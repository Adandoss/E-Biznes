package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"

	"sklep/models"
)

const (
	ErrKey             = "error"
	ProductNotFoundMsg = "Produkt nie został znaleziony"
	InvalidDataMsg     = "Niepoprawne dane"
	CreateProductMsg   = "Nie udało się utworzyć produktu"
	GetProductsMsg     = "Błąd podczas pobierania produktów"
)

type ProductHandler struct {
	DB *gorm.DB
}

func jsonError(c *echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{
		ErrKey: message,
	})
}

func (h *ProductHandler) CreateProduct(c *echo.Context) error {
	var product models.Product

	if err := c.Bind(&product); err != nil {
		return jsonError(c, http.StatusBadRequest, InvalidDataMsg)
	}

	if err := h.DB.Create(&product).Error; err != nil {
		return jsonError(c, http.StatusInternalServerError, CreateProductMsg)
	}

	return c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProducts(c *echo.Context) error {
	var products []models.Product
	if err := h.DB.Preload("Category").Order("price DESC").Find(&products).Error; err != nil {
		return jsonError(c, http.StatusInternalServerError, GetProductsMsg)
	}
	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *echo.Context) error {
	var product models.Product
	if err := h.DB.Preload("Category").First(&product, c.Param("id")).Error; err != nil {
		return jsonError(c, http.StatusNotFound, ProductNotFoundMsg)
	}
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) GetCategories(c *echo.Context) error {
	var categories []models.Category
	if err := h.DB.Order("name ASC").Find(&categories).Error; err != nil {
		return jsonError(c, http.StatusInternalServerError, "Błąd podczas pobierania kategorii")
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *ProductHandler) UpdateProduct(c *echo.Context) error {
	var product models.Product
	if err := h.DB.First(&product, c.Param("id")).Error; err != nil {
		return jsonError(c, http.StatusNotFound, ProductNotFoundMsg)
	}
	if err := c.Bind(&product); err != nil {
		return jsonError(c, http.StatusBadRequest, InvalidDataMsg)
	}
	if err := h.DB.Save(&product).Error; err != nil {
		return jsonError(c, http.StatusInternalServerError, "Nie udało się zaktualizować produktu")
	}
	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *echo.Context) error {
	var product models.Product
	if err := h.DB.First(&product, c.Param("id")).Error; err != nil {
		return jsonError(c, http.StatusNotFound, ProductNotFoundMsg)
	}
	if err := h.DB.Delete(&product).Error; err != nil {
		return jsonError(c, http.StatusInternalServerError, "Nie udało się usunąć produktu")
	}
	return c.NoContent(http.StatusNoContent)
}
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"

	"sklep/models"
)

type CartHandler struct {
	DB *gorm.DB
}

func (h *CartHandler) CreateCart(c *echo.Context) error {
	cart := models.Cart{Status: "aktywny"}

	if err := h.DB.Create(&cart).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się utworzyć koszyka"})
	}

	return c.JSON(http.StatusCreated, cart)
}

func (h *CartHandler) GetCart(c *echo.Context) error {
	id := c.Param("id")
	var cart models.Cart

	if err := h.DB.Preload("Items.Product").First(&cart, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Koszyk nie znaleziony"})
	}

	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddItem(c *echo.Context) error {
	cartID := c.Param("id")

	var cart models.Cart
	if err := h.DB.First(&cart, cartID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Koszyk nie znaleziony"})
	}

	var item models.CartItem
	if err := c.Bind(&item); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawne dane"})
	}
	item.CartID = cart.ID

	var product models.Product
	if err := h.DB.First(&product, item.ProductID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Produkt nie znaleziony"})
	}

	if err := h.DB.Create(&item).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się dodać produktu do koszyka"})
	}

	h.DB.Preload("Product").First(&item, item.ID)

	return c.JSON(http.StatusCreated, item)
}

func (h *CartHandler) RemoveItem(c *echo.Context) error {
	cartID := c.Param("id")
	itemID := c.Param("itemId")

	var item models.CartItem
	if err := h.DB.Where("cart_id = ? AND id = ?", cartID, itemID).First(&item).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Element nie znaleziony w koszyku"})
	}

	h.DB.Delete(&item)
	return c.NoContent(http.StatusNoContent)
}

func (h *CartHandler) DeleteCart(c *echo.Context) error {
	id := c.Param("id")

	var cart models.Cart
	if err := h.DB.First(&cart, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Koszyk nie znaleziony"})
	}

	h.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	h.DB.Delete(&cart)

	return c.NoContent(http.StatusNoContent)
}
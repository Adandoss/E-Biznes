package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"

	"sklep/models"
)

const CartNotFoundMsg string = "Koszyk nie został znaleziony"

type CartHandler struct {
	DB *gorm.DB
}

func (h *CartHandler) GetMyCart(c *echo.Context) error {
	userID, err := GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Nie zalogowano"})
	}

	var cart models.Cart
	result := h.DB.Where("user_id = ?", userID).Preload("Items.Product").First(&cart)

	if result.Error != nil {
		cart = models.Cart{UserID: userID, Status: "aktywny"}
		if err := h.DB.Create(&cart).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się utworzyć koszyka"})
		}
	}

	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) CreateCart(c *echo.Context) error {
	userID, _ := GetUserID(c)
	cart := models.Cart{UserID: userID, Status: "aktywny"}

	if err := h.DB.Create(&cart).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się utworzyć koszyka"})
	}

	return c.JSON(http.StatusCreated, cart)
}

func (h *CartHandler) GetCart(c *echo.Context) error {
	var cart models.Cart
	if err := h.DB.Preload("Items.Product").First(&cart, c.Param("id")).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": CartNotFoundMsg})
	}
	return c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddItem(c *echo.Context) error {
	var cart models.Cart
	if err := h.DB.First(&cart, c.Param("id")).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": CartNotFoundMsg})
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
	var item models.CartItem
	if err := h.DB.Where("cart_id = ? AND id = ?", c.Param("id"), c.Param("itemId")).First(&item).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Element nie znaleziony w koszyku"})
	}

	h.DB.Delete(&item)
	return c.NoContent(http.StatusNoContent)
}

func (h *CartHandler) DeleteCart(c *echo.Context) error {
	var cart models.Cart
	if err := h.DB.First(&cart, c.Param("id")).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": CartNotFoundMsg})
	}

	h.DB.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	h.DB.Delete(&cart)
	return c.NoContent(http.StatusNoContent)
}
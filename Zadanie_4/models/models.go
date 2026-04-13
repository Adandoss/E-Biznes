package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name"`
	Products []Product `json:"products,omitempty"` 
}

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	
	CategoryID  uint     `json:"category_id"` 
	Category    Category `json:"category,omitempty"` 
}

type Cart struct {
	gorm.Model
	Status string     `json:"status"`
	Items  []CartItem `json:"items,omitempty"`
}

type CartItem struct {
	gorm.Model
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product,omitempty"`
	Quantity  uint    `json:"quantity"`
}
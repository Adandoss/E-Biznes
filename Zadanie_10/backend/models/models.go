package models

import "gorm.io/gorm"


type Category struct {
	gorm.Model
	Name string `json:"name"`
}

type Product struct {
	gorm.Model
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	CategoryID  uint     `json:"category_id"`
	Category    Category `json:"category,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Cart struct {
	gorm.Model
	UserID uint       `json:"user_id"`
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

type Payment struct {
	gorm.Model
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}

type User struct {
	gorm.Model
	Email       string `json:"email" gorm:"uniqueIndex:idx_email_provider"`
	Password    string `json:"-"`
	Provider    string `json:"provider" gorm:"uniqueIndex:idx_email_provider"`
	ProviderID  string `json:"provider_id"`
	Name        string `json:"name"`
	AccessToken string `json:"-"`
}
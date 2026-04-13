package models

import "gorm.io/gorm"

func WithCategory(db *gorm.DB) *gorm.DB {
	return db.Preload("Category")
}

func PriceGreaterThan(minPrice float64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("price > ?", minPrice)
	}
}

func OrderByPriceDesc(db *gorm.DB) *gorm.DB {
	return db.Order("price DESC")
}
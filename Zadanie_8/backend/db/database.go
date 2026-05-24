package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	
	"sklep/models" 
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Nie połączono z bazą:", err)
	}

	err = db.AutoMigrate(
		&models.Category{},
		&models.Product{}, 
		&models.Cart{},
		&models.CartItem{},
		&models.Payment{},
		&models.User{},
	)
	
	if err != nil {
		log.Fatal("Błąd podczas migracji:", err)
	}

	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count == 0 {
		cats := []models.Category{
			{Name: "Elektronika"},
			{Name: "Książki"},
			{Name: "Ubrania"},
		}
		for i := range cats {
			db.Create(&cats[i])
		}

		products := []models.Product{
			{Name: "Laptop", Description: "laptop", Price: 3500.0, CategoryID: cats[0].ID},
			{Name: "Książka", Description: "książka", Price: 69.90, CategoryID: cats[1].ID},
			{Name: "Koszulka", Description: "koszulka", Price: 49.99, CategoryID: cats[2].ID},
		}
		for i := range products {
			db.Create(&products[i])
		}
	}

	log.Println("Pomyślnie połączono z bazą.")
	return db
}
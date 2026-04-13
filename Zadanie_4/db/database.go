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
	)
	
	if err != nil {
		log.Fatal("Błąd podczas migracji:", err)
	}

	log.Println("Pomyślnie połączono z bazą.")
	return db
}
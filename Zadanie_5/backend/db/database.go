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
		&models.Product{}, 
		&models.Cart{},
		&models.CartItem{},
		&models.Payment{},
	)
	
	if err != nil {
		log.Fatal("Błąd podczas migracji:", err)
	}

	var count int64
	db.Model(&models.Cart{}).Count(&count)
	if count == 0 {
		db.Create(&models.Cart{Status: "aktywny"})
		log.Println("Utworzono domyślny koszyk.")
	}

	log.Println("Pomyślnie połączono z bazą.")
	return db
}
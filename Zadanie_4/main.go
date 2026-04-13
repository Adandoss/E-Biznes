package main

import (
	"log"

	"github.com/labstack/echo/v5"
	
	"sklep/db"
	"sklep/handlers"
)

func main() {

	database := db.InitDB()
	e := echo.New()

	productHandler := &handlers.ProductHandler{DB: database}
	cartHandler := &handlers.CartHandler{DB: database}

	e.POST("/products", productHandler.CreateProduct)     
	e.GET("/products", productHandler.GetProducts)        
	e.GET("/products/:id", productHandler.GetProduct)      
	e.PUT("/products/:id", productHandler.UpdateProduct)   
	e.DELETE("/products/:id", productHandler.DeleteProduct) 

	e.POST("/carts", cartHandler.CreateCart)                     
	e.GET("/carts/:id", cartHandler.GetCart)                        
	e.POST("/carts/:id/items", cartHandler.AddItem)                  
	e.DELETE("/carts/:id/items/:itemId", cartHandler.RemoveItem)     
	e.DELETE("/carts/:id", cartHandler.DeleteCart)                    

	log.Fatal(e.Start(":8080"))
}
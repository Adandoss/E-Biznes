package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	
	"sklep/db"
	"sklep/handlers"
)

func main() {

	database := db.InitDB()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins: []string{"http://localhost:5173"},
	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	productHandler := &handlers.ProductHandler{DB: database}
	cartHandler := &handlers.CartHandler{DB: database}
	paymentHandler := &handlers.PaymentHandler{DB: database}

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

	e.GET("/payments", paymentHandler.GetPayments)
	e.POST("/payments", paymentHandler.CreatePayment)

	log.Fatal(e.Start(":8080"))
}
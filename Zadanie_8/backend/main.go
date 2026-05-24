package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	
	"sklep/db"
	"sklep/handlers"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("Brak pliku .env — używam zmiennych środowiskowych systemu")
	}

	database := db.InitDB()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOrigins:     []string{"http://localhost:5173"},
	AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	AllowHeaders:     []string{"Content-Type", "Authorization"},
	AllowCredentials: true,
	}))

	productHandler := &handlers.ProductHandler{DB: database}
	cartHandler := &handlers.CartHandler{DB: database}
	paymentHandler := &handlers.PaymentHandler{DB: database}
	authHandler := &handlers.AuthHandler{DB: database}

	products := e.Group("/products")
	{
		products.POST("", productHandler.CreateProduct)
		products.GET("", productHandler.GetProducts)
		products.GET("/:id", productHandler.GetProduct)
		products.PUT("/:id", productHandler.UpdateProduct)
		products.DELETE("/:id", productHandler.DeleteProduct)
	}

	e.GET("/categories", productHandler.GetCategories)

	// Register /carts/mine on the root echo instance BEFORE the parameterized
	// carts group so Echo's router matches the static path first.
	e.GET("/carts/mine", cartHandler.GetMyCart, handlers.JWTMiddleware)

	carts := e.Group("/carts", handlers.JWTMiddleware)
	{
		carts.POST("", cartHandler.CreateCart)
		carts.GET("/:id", cartHandler.GetCart)
		carts.POST("/:id/items", cartHandler.AddItem)
		carts.DELETE("/:id/items/:itemId", cartHandler.RemoveItem)
		carts.DELETE("/:id", cartHandler.DeleteCart)
	}

	payments := e.Group("/payments", handlers.JWTMiddleware)
	{
		payments.GET("", paymentHandler.GetPayments)
		payments.POST("", paymentHandler.CreatePayment)
	}

	auth := e.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/google", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)
		auth.GET("/github", authHandler.GithubLogin)
		auth.GET("/github/callback", authHandler.GithubCallback)
		auth.GET("/me", authHandler.Me, handlers.JWTMiddleware)
	}

	log.Fatal(e.Start(":8080"))
}
package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
)

func GetUserID(c *echo.Context) (uint, error) {
	if claims, ok := c.Get("user").(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(float64); ok {
			return uint(userID), nil
		}
	}
	return 0, fmt.Errorf("invalid user in context")
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Brak tokenu autoryzacyjnego"})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Niepoprawny format tokenu"})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Niepoprawny lub wygasły token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Niepoprawne dane tokenu"})
		}

		c.Set("user", claims)
		return next(c)
	}
}

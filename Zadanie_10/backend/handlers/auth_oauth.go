package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"

	"sklep/models"
)

func (h *AuthHandler) oauthLogin(c *echo.Context, provider string) error {
	state := generateState()
	setStateCookie(c, state)
	url := getOAuthConfig(provider).AuthCodeURL(state)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// GET /auth/google 
func (h *AuthHandler) GoogleLogin(c *echo.Context) error {
	return h.oauthLogin(c, "google")
}

func (h *AuthHandler) handleOAuthCallback(c *echo.Context, provider, providerID, email, name, token string) error {
	var user models.User
	result := h.DB.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user)

	if result.Error != nil {
		user = models.User{
			Email:       email,
			Name:        name,
			Provider:    provider,
			ProviderID:  providerID,
			AccessToken: token,
		}
		if err := h.DB.Create(&user).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się utworzyć użytkownika"})
		}
	} else {
		if err := h.DB.Model(&user).Updates(models.User{
			AccessToken: token,
			Name:        name,
			Email:       email,
		}).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się zaktualizować użytkownika"})
		}
	}

	jwtToken, err := generateJWT(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd generowania tokenu"})
	}

	redirectURL := fmt.Sprintf("http://localhost:5173/auth/callback?token=%s", jwtToken)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GET /auth/google/callback
func (h *AuthHandler) GoogleCallback(c *echo.Context) error {
	if err := validateState(c); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawny stan OAuth"})
	}

	config := getOAuthConfig("google")
	code := c.QueryParam("code")

	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Brak kodu autoryzacyjnego"})
	}

	token, err := config.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się wymienić kodu na token"})
	}

	client := config.Client(c.Request().Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się pobrać danych użytkownika"})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd odczytu danych"})
	}

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd parsowania danych"})
	}

	return h.handleOAuthCallback(c, "google", googleUser.ID, googleUser.Email, googleUser.Name, token.AccessToken)
}

// GET /auth/github 
func (h *AuthHandler) GithubLogin(c *echo.Context) error {
	return h.oauthLogin(c, "github")
}

// GET /auth/github/callback
func (h *AuthHandler) GithubCallback(c *echo.Context) error {
	if err := validateState(c); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Niepoprawny stan OAuth"})
	}

	config := getOAuthConfig("github")
	code := c.QueryParam("code")

	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Brak kodu autoryzacyjnego"})
	}

	token, err := config.Exchange(c.Request().Context(), code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się wymienić kodu na token"})
	}

	client := config.Client(c.Request().Context(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Nie udało się pobrać danych użytkownika"})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd odczytu danych"})
	}

	var githubUser struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &githubUser); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Błąd parsowania danych"})
	}

	if githubUser.Email == "" {
		emailResp, emailErr := client.Get("https://api.github.com/user/emails")
		if emailErr == nil {
			defer emailResp.Body.Close()
			emailBody, readErr := io.ReadAll(emailResp.Body)
			if readErr == nil {
				var emails []struct {
					Email   string `json:"email"`
					Primary bool   `json:"primary"`
				}
				if json.Unmarshal(emailBody, &emails) == nil {
					for _, e := range emails {
						if e.Primary {
							githubUser.Email = e.Email
							break
						}
					}
				}
			}
		}
	}

	displayName := githubUser.Name
	if displayName == "" {
		displayName = githubUser.Login
	}

	return h.handleOAuthCallback(c, "github", fmt.Sprintf("%d", githubUser.ID), githubUser.Email, displayName, token.AccessToken)
}

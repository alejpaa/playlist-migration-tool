package handlers

import (
	"net/http"

	"github.com/alejpaa/playlist-migration-tool/internal/models"
	"github.com/alejpaa/playlist-migration-tool/internal/services"
	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// AuthenticateYouTube handles YouTube authentication
func (h *AuthHandler) AuthenticateYouTube(c *gin.Context) {
	response, err := h.authService.AuthenticateWithYouTube()
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		} else {
			apiErr := models.NewInternalServerError("Authentication failed", err)
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetYouTubeAuthURL obtiene la URL de autenticación de YouTube
func (h *AuthHandler) GetYouTubeAuthURL(c *gin.Context) {
	authURL, err := h.authService.GetYouTubeAuthURL()
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		} else {
			apiErr := models.NewInternalServerError("Failed to generate auth URL", err)
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"message":  "Visita esta URL para autorizar la aplicación con tu cuenta de YouTube",
	})
}

// CompleteYouTubeAuth completa la autenticación con el código de autorización
func (h *AuthHandler) CompleteYouTubeAuth(c *gin.Context) {
	var request struct {
		AuthCode string `json:"auth_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := models.NewBadRequestError("Auth code is required", err)
		c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		return
	}

	response, err := h.authService.CompleteYouTubeAuth(request.AuthCode)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		} else {
			apiErr := models.NewInternalServerError("Authentication failed", err)
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

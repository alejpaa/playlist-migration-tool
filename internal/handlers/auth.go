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

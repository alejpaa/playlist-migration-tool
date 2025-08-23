package handlers

import (
	"net/http"

	"github.com/alejpaa/playlist-migration-tool/internal/models"
	"github.com/alejpaa/playlist-migration-tool/internal/services"
	"github.com/gin-gonic/gin"
)

// ExportHandler handles export endpoints
type ExportHandler struct {
	exportService *services.ExportService
}

// NewExportHandler creates a new ExportHandler
func NewExportHandler(exportService *services.ExportService) *ExportHandler {
	return &ExportHandler{
		exportService: exportService,
	}
}

// ExportPlaylist handles POST /export/:id
func (h *ExportHandler) ExportPlaylist(c *gin.Context) {
	// Get access token from context
	accessToken, exists := c.Get("access_token")
	if !exists {
		apiErr := models.NewUnauthorizedError("Access token not found", nil)
		c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		return
	}

	accessTokenStr, ok := accessToken.(string)
	if !ok {
		apiErr := models.NewUnauthorizedError("Invalid access token format", nil)
		c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		return
	}

	// Get playlist ID from URL parameter
	playlistID := c.Param("id")
	if playlistID == "" {
		apiErr := models.NewBadRequestError("Playlist ID is required", nil)
		c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		return
	}

	// Parse request body
	var request models.ExportRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := models.NewBadRequestError("Invalid request body", err)
		c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		return
	}

	// Validate format
	if request.Format == "" {
		request.Format = "json" // default
	}

	// Export playlist
	response, err := h.exportService.ExportPlaylist(accessTokenStr, playlistID, &request)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		} else {
			apiErr := models.NewInternalServerError("Failed to export playlist", err)
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

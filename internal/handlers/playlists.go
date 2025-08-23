package handlers

import (
	"net/http"
	"strconv"

	"github.com/alejpaa/playlist-migration-tool/internal/models"
	"github.com/alejpaa/playlist-migration-tool/internal/services"
	"github.com/gin-gonic/gin"
)

// PlaylistHandler handles playlist endpoints
type PlaylistHandler struct {
	playlistService *services.PlaylistService
}

// NewPlaylistHandler creates a new PlaylistHandler
func NewPlaylistHandler(playlistService *services.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{
		playlistService: playlistService,
	}
}

// GetPlaylists handles GET /playlists
func (h *PlaylistHandler) GetPlaylists(c *gin.Context) {
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

	// Parse query parameters
	maxResults := 25 // default
	if mr := c.Query("max_results"); mr != "" {
		if parsed, err := strconv.Atoi(mr); err == nil && parsed > 0 && parsed <= 50 {
			maxResults = parsed
		}
	}

	pageToken := c.Query("page_token")

	// Get playlists
	response, err := h.playlistService.GetPlaylists(accessTokenStr, maxResults, pageToken)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		} else {
			apiErr := models.NewInternalServerError("Failed to fetch playlists", err)
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPlaylistByID handles GET /playlists/:id
func (h *PlaylistHandler) GetPlaylistByID(c *gin.Context) {
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

	// Get playlist
	response, err := h.playlistService.GetPlaylistByID(accessTokenStr, playlistID)
	if err != nil {
		if apiErr, ok := err.(*models.APIError); ok {
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		} else {
			apiErr := models.NewInternalServerError("Failed to fetch playlist", err)
			c.JSON(apiErr.StatusCode, apiErr.ToErrorResponse())
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

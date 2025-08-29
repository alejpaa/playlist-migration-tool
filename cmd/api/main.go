package main

import (
	"log"

	"github.com/alejpaa/playlist-migration-tool/internal/config"
	"github.com/alejpaa/playlist-migration-tool/internal/handlers"
	"github.com/alejpaa/playlist-migration-tool/internal/middleware"
	"github.com/alejpaa/playlist-migration-tool/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize services
	authService := services.NewAuthService(cfg.GoogleCredentialsFile)
	playlistService := services.NewPlaylistService()
	exportService := services.NewExportService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	playlistHandler := handlers.NewPlaylistHandler(playlistService)
	exportHandler := handlers.NewExportHandler(exportService)
	healthHandler := handlers.NewHealthHandler()

	// Setup Gin router
	router := gin.New()

	// Apply global middleware
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Health endpoint (no auth required)
	router.GET("/health", healthHandler.HealthCheck)

	// Auth endpoints
	auth := router.Group("/auth")
	{
		auth.GET("/youtube/url", authHandler.GetYouTubeAuthURL)
		auth.POST("/youtube/callback", authHandler.CompleteYouTubeAuth)
		auth.POST("/youtube", authHandler.AuthenticateYouTube)
	}

	// Protected API endpoints (require auth)
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Playlist endpoints
		api.GET("/playlists", playlistHandler.GetPlaylists)
		api.GET("/playlists/:id", playlistHandler.GetPlaylistByID)
		api.GET("/playlists/:id/songs", playlistHandler.GetPlaylistSongs)

		// Export endpoints
		api.POST("/export/:id", exportHandler.ExportPlaylist)
	}

	log.Printf("🚀 Server starting on port %s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}

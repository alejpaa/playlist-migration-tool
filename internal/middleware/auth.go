package middleware

import (
	"strings"

	"github.com/alejpaa/playlist-migration-tool/internal/models"
	"github.com/alejpaa/playlist-migration-tool/pkg/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifies the access token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			respondWithError(c, models.NewUnauthorizedError("Authorization header required", nil))
			c.Abort()
			return
		}

		// Check for Bearer token format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			respondWithError(c, models.NewUnauthorizedError("Invalid authorization header format", nil))
			c.Abort()
			return
		}

		accessToken := tokenParts[1]
		if accessToken == "" {
			respondWithError(c, models.NewUnauthorizedError("Access token required", nil))
			c.Abort()
			return
		}

		// Validate token by trying to use it (simple validation)
		if !auth.ValidateToken(accessToken) {
			respondWithError(c, models.NewUnauthorizedError("Invalid or expired token", nil))
			c.Abort()
			return
		}

		// Add token to context
		c.Set("access_token", accessToken)
		c.Next()
	}
}

// respondWithError writes an error response
func respondWithError(c *gin.Context, err *models.APIError) {
	c.JSON(err.StatusCode, err.ToErrorResponse())
}

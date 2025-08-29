package services

import (
	"context"
	"encoding/json"
	"os"

	"github.com/alejpaa/playlist-migration-tool/internal/models"
	"github.com/alejpaa/playlist-migration-tool/pkg/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

// AuthService handles authentication logic
type AuthService struct {
	credentialsFile string
}

// NewAuthService creates a new AuthService
func NewAuthService(credentialsFile string) *AuthService {
	return &AuthService{
		credentialsFile: credentialsFile,
	}
}

// GetYouTubeAuthURL genera la URL de autenticación de YouTube
func (s *AuthService) GetYouTubeAuthURL() (string, error) {
	// Read credentials file
	b, err := os.ReadFile(s.credentialsFile)
	if err != nil {
		return "", models.NewInternalServerError("Unable to read credentials file", err)
	}

	var creds auth.Credentials
	if err := json.Unmarshal(b, &creds); err != nil {
		return "", models.NewInternalServerError("Unable to parse credentials file", err)
	}

	// Create OAuth2 config
	config := &oauth2.Config{
		ClientID:     creds.Installed.ClientID,
		ClientSecret: creds.Installed.ClientSecret,
		RedirectURL:  creds.Installed.RedirectURIs[0],
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL, nil
}

// CompleteYouTubeAuth completa la autenticación con el código de autorización
func (s *AuthService) CompleteYouTubeAuth(authCode string) (*models.AuthResponse, error) {
	// Read credentials file
	b, err := os.ReadFile(s.credentialsFile)
	if err != nil {
		return nil, models.NewInternalServerError("Unable to read credentials file", err)
	}

	var creds auth.Credentials
	if err := json.Unmarshal(b, &creds); err != nil {
		return nil, models.NewInternalServerError("Unable to parse credentials file", err)
	}

	// Create OAuth2 config
	config := &oauth2.Config{
		ClientID:     creds.Installed.ClientID,
		ClientSecret: creds.Installed.ClientSecret,
		RedirectURL:  creds.Installed.RedirectURIs[0],
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	// Exchange authorization code for token
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, models.NewBadRequestError("Invalid authorization code", err)
	}

	// Save token
	f, err := os.OpenFile("token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, models.NewInternalServerError("Unable to save token", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(tok)

	return &models.AuthResponse{
		Success:     true,
		AccessToken: tok.AccessToken,
		Message:     "Successfully authenticated with YouTube",
	}, nil
}

// AuthenticateWithYouTube handles YouTube OAuth authentication
func (s *AuthService) AuthenticateWithYouTube() (*models.AuthResponse, error) {
	accessToken, err := auth.GetAccessToken(s.credentialsFile)
	if err != nil {
		return nil, models.NewInternalServerError("Failed to authenticate with YouTube", err)
	}

	return &models.AuthResponse{
		Success:     true,
		AccessToken: accessToken,
		Message:     "Successfully authenticated with YouTube",
	}, nil
}

// ValidateToken validates an access token
func (s *AuthService) ValidateToken(token string) bool {
	return auth.ValidateToken(token)
}

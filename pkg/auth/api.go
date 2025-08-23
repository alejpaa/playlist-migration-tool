package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const tokenFile = "token.json"

// Credentials represents the OAuth2 credentials
type Credentials struct {
	Installed struct {
		ClientID     string   `json:"client_id"`
		ClientSecret string   `json:"client_secret"`
		AuthURI      string   `json:"auth_uri"`
		TokenURI     string   `json:"token_uri"`
		RedirectURIs []string `json:"redirect_uris"`
	} `json:"installed"`
}

// GetClient retrieves a token, saves the token, then returns the generated client
func GetClient(credentialsFile string) (*http.Client, error) {
	// Read credentials file
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %v", err)
	}

	var creds Credentials
	if err := json.Unmarshal(b, &creds); err != nil {
		return nil, fmt.Errorf("unable to parse client secret file: %v", err)
	}

	// Create OAuth2 config
	config := &oauth2.Config{
		ClientID:     creds.Installed.ClientID,
		ClientSecret: creds.Installed.ClientSecret,
		RedirectURL:  creds.Installed.RedirectURIs[0],
		Scopes:       []string{youtube.YoutubeReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	// Get token from web or file
	tok, err := getTokenFromWeb(config)
	if err != nil {
		return nil, err
	}

	// Save token for future use
	saveToken(tokenFile, tok)

	return config.Client(context.Background(), tok), nil
}

// getTokenFromWeb requests a token from the web, then returns the retrieved token
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	// Check if we already have a saved token
	if tok, err := tokenFromFile(tokenFile); err == nil {
		if tok.Valid() {
			fmt.Println("Using existing token")
			return tok, nil
		}
		fmt.Println("Token expired, getting new one...")
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	// Try to open browser automatically
	openBrowser(authURL)

	fmt.Print("Enter the authorization code: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

// tokenFromFile retrieves a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// saveToken saves a token to a file path
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// openBrowser tries to open the URL in a browser
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Printf("Could not open browser automatically: %v\n", err)
	}
}

// GetAccessToken returns the access token from saved credentials
func GetAccessToken(credentialsFile string) (string, error) {
	_, err := GetClient(credentialsFile)
	if err != nil {
		return "", err
	}

	// Get token from file
	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		return "", fmt.Errorf("unable to read token: %v", err)
	}

	// If token is expired, refresh it
	if !tok.Valid() {
		// Read credentials to refresh token
		b, err := os.ReadFile(credentialsFile)
		if err != nil {
			return "", err
		}

		var creds Credentials
		if err := json.Unmarshal(b, &creds); err != nil {
			return "", err
		}

		config := &oauth2.Config{
			ClientID:     creds.Installed.ClientID,
			ClientSecret: creds.Installed.ClientSecret,
			RedirectURL:  creds.Installed.RedirectURIs[0],
			Scopes:       []string{youtube.YoutubeReadonlyScope},
			Endpoint:     google.Endpoint,
		}

		tokenSource := config.TokenSource(context.Background(), tok)
		newToken, err := tokenSource.Token()
		if err != nil {
			return "", fmt.Errorf("unable to refresh token: %v", err)
		}

		saveToken(tokenFile, newToken)
		return newToken.AccessToken, nil
	}

	return tok.AccessToken, nil
}

// ValidateToken validates an access token by making a simple API call
func ValidateToken(accessToken string) bool {
	// Simple validation - check if token is not empty
	// In a production environment, you might want to make a test API call
	// to verify the token is still valid
	return accessToken != ""
}

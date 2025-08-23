package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client representa un cliente para la API de YouTube
type Client struct {
	accessToken string
	httpClient  *http.Client
}

// NewClient crea una nueva instancia del cliente de YouTube
func NewClient(accessToken string) *Client {
	return &Client{
		accessToken: accessToken,
		httpClient:  &http.Client{},
	}
}

// PlaylistsResponse representa la respuesta de la API de playlists
type PlaylistsResponse struct {
	Kind          string     `json:"kind"`
	Etag          string     `json:"etag"`
	NextPageToken string     `json:"nextPageToken"`
	PrevPageToken string     `json:"prevPageToken"`
	PageInfo      PageInfo   `json:"pageInfo"`
	Items         []Playlist `json:"items"`
}

// PageInfo contiene información sobre la paginación
type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

// Playlist representa una playlist de YouTube
type Playlist struct {
	Kind           string                 `json:"kind"`
	Etag           string                 `json:"etag"`
	ID             string                 `json:"id"`
	Snippet        PlaylistSnippet        `json:"snippet"`
	Status         PlaylistStatus         `json:"status,omitempty"`
	ContentDetails PlaylistContentDetails `json:"contentDetails,omitempty"`
}

// PlaylistSnippet contiene información básica de la playlist
type PlaylistSnippet struct {
	PublishedAt     string                   `json:"publishedAt"`
	ChannelID       string                   `json:"channelId"`
	Title           string                   `json:"title"`
	Description     string                   `json:"description"`
	Thumbnails      map[string]Thumbnail     `json:"thumbnails"`
	ChannelTitle    string                   `json:"channelTitle"`
	DefaultLanguage string                   `json:"defaultLanguage,omitempty"`
	Localized       PlaylistLocalizedSnippet `json:"localized,omitempty"`
}

// PlaylistStatus contiene el estado de la playlist
type PlaylistStatus struct {
	PrivacyStatus string `json:"privacyStatus"`
}

// PlaylistContentDetails contiene detalles del contenido
type PlaylistContentDetails struct {
	ItemCount int `json:"itemCount"`
}

// PlaylistLocalizedSnippet contiene información localizada
type PlaylistLocalizedSnippet struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Thumbnail representa una miniatura
type Thumbnail struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// ListPlaylists obtiene las playlists del usuario autenticado
func (c *Client) ListPlaylists(options *ListPlaylistsOptions) (*PlaylistsResponse, error) {
	if options == nil {
		options = &ListPlaylistsOptions{
			Part:       "snippet,status,contentDetails",
			Mine:       true,
			MaxResults: 25,
		}
	}

	// Asegurar que Part no esté vacío
	if options.Part == "" {
		options.Part = "snippet,status,contentDetails"
	}

	// Asegurar que MaxResults esté en rango válido (1-50)
	if options.MaxResults <= 0 {
		options.MaxResults = 25
	} else if options.MaxResults > 50 {
		options.MaxResults = 50
	}

	// Construir la URL con parámetros
	baseURL := "https://www.googleapis.com/youtube/v3/playlists"
	params := url.Values{}

	params.Add("part", options.Part)
	if options.Mine {
		params.Add("mine", "true")
	}
	if options.ChannelID != "" {
		params.Add("channelId", options.ChannelID)
	}
	if options.MaxResults > 0 {
		params.Add("maxResults", fmt.Sprintf("%d", options.MaxResults))
	}
	if options.PageToken != "" {
		params.Add("pageToken", options.PageToken)
	}

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Crear la petición HTTP
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}

	// Agregar el header de autorización
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Accept", "application/json")

	// Realizar la petición
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error realizando petición: %w", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornó código %d: %s", resp.StatusCode, string(body))
	}

	// Parsear la respuesta JSON
	var playlistsResp PlaylistsResponse
	if err := json.Unmarshal(body, &playlistsResp); err != nil {
		return nil, fmt.Errorf("error parseando JSON: %w", err)
	}

	return &playlistsResp, nil
}

// ListPlaylistsOptions contiene las opciones para listar playlists
type ListPlaylistsOptions struct {
	Part       string // Partes a incluir: snippet, status, contentDetails, etc.
	Mine       bool   // Si true, obtiene las playlists del usuario autenticado
	ChannelID  string // ID del canal (alternativa a Mine)
	MaxResults int    // Número máximo de resultados (1-50, default: 5)
	PageToken  string // Token para paginación
}

// ListMyPlaylists es una función de conveniencia para listar las playlists del usuario autenticado
func (c *Client) ListMyPlaylists() (*PlaylistsResponse, error) {
	return c.ListPlaylists(&ListPlaylistsOptions{
		Part:       "snippet,status,contentDetails",
		Mine:       true,
		MaxResults: 50,
	})
}

// GetPlaylistByID obtiene una playlist específica por su ID
func (c *Client) GetPlaylistByID(playlistID string) (*Playlist, error) {
	baseURL := "https://www.googleapis.com/youtube/v3/playlists"
	params := url.Values{}
	params.Add("part", "snippet,status,contentDetails")
	params.Add("id", playlistID)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error realizando petición: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornó código %d: %s", resp.StatusCode, string(body))
	}

	var playlistsResp PlaylistsResponse
	if err := json.Unmarshal(body, &playlistsResp); err != nil {
		return nil, fmt.Errorf("error parseando JSON: %w", err)
	}

	if len(playlistsResp.Items) == 0 {
		return nil, fmt.Errorf("playlist no encontrada")
	}

	return &playlistsResp.Items[0], nil
}

// PlaylistItem representa un video en una playlist
type PlaylistItem struct {
	Kind    string              `json:"kind"`
	Etag    string              `json:"etag"`
	ID      string              `json:"id"`
	Snippet PlaylistItemSnippet `json:"snippet"`
}

// PlaylistItemSnippet contiene información del video en la playlist
type PlaylistItemSnippet struct {
	PublishedAt  string               `json:"publishedAt"`
	ChannelID    string               `json:"channelId"`
	Title        string               `json:"title"`
	Description  string               `json:"description"`
	Thumbnails   map[string]Thumbnail `json:"thumbnails"`
	ChannelTitle string               `json:"channelTitle"`
	PlaylistID   string               `json:"playlistId"`
	Position     int                  `json:"position"`
	ResourceID   ResourceID           `json:"resourceId"`
}

// ResourceID identifica el recurso (video)
type ResourceID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

// PlaylistItemsResponse representa la respuesta de items de playlist
type PlaylistItemsResponse struct {
	Kind          string         `json:"kind"`
	Etag          string         `json:"etag"`
	NextPageToken string         `json:"nextPageToken"`
	PrevPageToken string         `json:"prevPageToken"`
	PageInfo      PageInfo       `json:"pageInfo"`
	Items         []PlaylistItem `json:"items"`
}

// ListPlaylistItems obtiene los videos de una playlist
func (c *Client) ListPlaylistItems(playlistID string, maxResults int) (*PlaylistItemsResponse, error) {
	if maxResults <= 0 {
		maxResults = 50
	}

	baseURL := "https://www.googleapis.com/youtube/v3/playlistItems"
	params := url.Values{}
	params.Add("part", "snippet")
	params.Add("playlistId", playlistID)
	params.Add("maxResults", fmt.Sprintf("%d", maxResults))

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creando petición: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error realizando petición: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornó código %d: %s", resp.StatusCode, string(body))
	}

	var itemsResp PlaylistItemsResponse
	if err := json.Unmarshal(body, &itemsResp); err != nil {
		return nil, fmt.Errorf("error parseando JSON: %w", err)
	}

	return &itemsResp, nil
}

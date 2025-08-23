package models

// AuthRequest represents a YouTube authentication request
type AuthRequest struct {
	RedirectURI string `json:"redirect_uri,omitempty"`
}

// ExportRequest represents a playlist export request
type ExportRequest struct {
	Format      string            `json:"format"`       // "json", "csv", "m3u"
	Options     map[string]string `json:"options"`      // Format-specific options
	IncludeInfo bool              `json:"include_info"` // Include video metadata
}

// PlaylistsQuery represents query parameters for listing playlists
type PlaylistsQuery struct {
	MaxResults int    `json:"max_results,omitempty"`
	PageToken  string `json:"page_token,omitempty"`
}

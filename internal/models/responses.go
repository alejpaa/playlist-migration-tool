package models

import "time"

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	Success     bool   `json:"success"`
	AccessToken string `json:"access_token,omitempty"`
	Message     string `json:"message"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// PlaylistResponse represents a simplified playlist response
type PlaylistResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	VideoCount    int       `json:"video_count"`
	PrivacyStatus string    `json:"privacy_status"`
	CreatedAt     time.Time `json:"created_at"`
	ChannelTitle  string    `json:"channel_title"`
	ThumbnailURL  string    `json:"thumbnail_url"`
}

// PlaylistsResponse represents a collection of playlists
type PlaylistsResponse struct {
	Playlists     []PlaylistResponse `json:"playlists"`
	TotalCount    int                `json:"total_count"`
	NextPageToken string             `json:"next_page_token,omitempty"`
	PrevPageToken string             `json:"prev_page_token,omitempty"`
}

// PlaylistDetailResponse represents detailed playlist information
type PlaylistDetailResponse struct {
	PlaylistResponse
	Videos []VideoResponse `json:"videos,omitempty"`
}

// VideoResponse represents a video in a playlist
type VideoResponse struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ChannelTitle string    `json:"channel_title"`
	Duration     string    `json:"duration,omitempty"`
	Position     int       `json:"position"`
	AddedAt      time.Time `json:"added_at"`
	ThumbnailURL string    `json:"thumbnail_url"`
}

// ExportResponse represents the response after exporting a playlist
type ExportResponse struct {
	Success     bool   `json:"success"`
	Format      string `json:"format"`
	Data        string `json:"data,omitempty"`         // For small exports
	DownloadURL string `json:"download_url,omitempty"` // For large exports
	Message     string `json:"message"`
}

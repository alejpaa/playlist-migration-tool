package services

import (
	"time"

	"github.com/alejpaa/playlist-migration-tool/internal/models"
	"github.com/alejpaa/playlist-migration-tool/pkg/youtube"
)

// PlaylistService handles playlist-related business logic
type PlaylistService struct{}

// NewPlaylistService creates a new PlaylistService
func NewPlaylistService() *PlaylistService {
	return &PlaylistService{}
}

// GetPlaylists retrieves user's playlists
func (s *PlaylistService) GetPlaylists(accessToken string, maxResults int, pageToken string) (*models.PlaylistsResponse, error) {
	client := youtube.NewClient(accessToken)

	options := &youtube.ListPlaylistsOptions{
		Part:       "snippet,status,contentDetails",
		Mine:       true,
		MaxResults: maxResults,
		PageToken:  pageToken,
	}

	response, err := client.ListPlaylists(options)
	if err != nil {
		return nil, models.NewInternalServerError("Failed to fetch playlists", err)
	}

	// Convert YouTube API response to our internal model
	playlists := make([]models.PlaylistResponse, len(response.Items))
	for i, item := range response.Items {
		createdAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		thumbnailURL := ""
		if item.Snippet.Thumbnails != nil {
			if medium, ok := item.Snippet.Thumbnails["medium"]; ok {
				thumbnailURL = medium.URL
			} else if def, ok := item.Snippet.Thumbnails["default"]; ok {
				thumbnailURL = def.URL
			}
		}

		playlists[i] = models.PlaylistResponse{
			ID:            item.ID,
			Title:         item.Snippet.Title,
			Description:   item.Snippet.Description,
			VideoCount:    item.ContentDetails.ItemCount,
			PrivacyStatus: item.Status.PrivacyStatus,
			CreatedAt:     createdAt,
			ChannelTitle:  item.Snippet.ChannelTitle,
			ThumbnailURL:  thumbnailURL,
		}
	}

	return &models.PlaylistsResponse{
		Playlists:     playlists,
		TotalCount:    response.PageInfo.TotalResults,
		NextPageToken: response.NextPageToken,
		PrevPageToken: response.PrevPageToken,
	}, nil
}

// GetPlaylistByID retrieves a specific playlist with its videos
func (s *PlaylistService) GetPlaylistByID(accessToken, playlistID string) (*models.PlaylistDetailResponse, error) {
	client := youtube.NewClient(accessToken)

	// Get playlist info
	playlist, err := client.GetPlaylistByID(playlistID)
	if err != nil {
		return nil, models.NewNotFoundError("Playlist not found", err)
	}

	// Get playlist items
	items, err := client.ListPlaylistItems(playlistID, 50)
	if err != nil {
		return nil, models.NewInternalServerError("Failed to fetch playlist items", err)
	}

	// Convert playlist to our model
	createdAt, _ := time.Parse(time.RFC3339, playlist.Snippet.PublishedAt)

	thumbnailURL := ""
	if playlist.Snippet.Thumbnails != nil {
		if medium, ok := playlist.Snippet.Thumbnails["medium"]; ok {
			thumbnailURL = medium.URL
		} else if def, ok := playlist.Snippet.Thumbnails["default"]; ok {
			thumbnailURL = def.URL
		}
	}

	// Convert videos to our model
	videos := make([]models.VideoResponse, len(items.Items))
	for i, item := range items.Items {
		addedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		videoThumbnailURL := ""
		if item.Snippet.Thumbnails != nil {
			if medium, ok := item.Snippet.Thumbnails["medium"]; ok {
				videoThumbnailURL = medium.URL
			} else if def, ok := item.Snippet.Thumbnails["default"]; ok {
				videoThumbnailURL = def.URL
			}
		}

		videos[i] = models.VideoResponse{
			ID:           item.Snippet.ResourceID.VideoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			ChannelTitle: item.Snippet.ChannelTitle,
			Position:     item.Snippet.Position,
			AddedAt:      addedAt,
			ThumbnailURL: videoThumbnailURL,
		}
	}

	return &models.PlaylistDetailResponse{
		PlaylistResponse: models.PlaylistResponse{
			ID:            playlist.ID,
			Title:         playlist.Snippet.Title,
			Description:   playlist.Snippet.Description,
			VideoCount:    playlist.ContentDetails.ItemCount,
			PrivacyStatus: playlist.Status.PrivacyStatus,
			CreatedAt:     createdAt,
			ChannelTitle:  playlist.Snippet.ChannelTitle,
			ThumbnailURL:  thumbnailURL,
		},
		Videos: videos,
	}, nil
}

// GetPlaylistSongs obtiene solo las canciones de una playlist
func (s *PlaylistService) GetPlaylistSongs(accessToken, playlistID string, maxResults int) ([]models.VideoResponse, error) {
	client := youtube.NewClient(accessToken)

	// Get playlist items
	items, err := client.ListPlaylistItems(playlistID, maxResults)
	if err != nil {
		return nil, models.NewInternalServerError("Failed to fetch playlist songs", err)
	}

	// Convert videos to our model
	videos := make([]models.VideoResponse, len(items.Items))
	for i, item := range items.Items {
		addedAt, _ := time.Parse(time.RFC3339, item.Snippet.PublishedAt)

		videoThumbnailURL := ""
		if item.Snippet.Thumbnails != nil {
			if medium, ok := item.Snippet.Thumbnails["medium"]; ok {
				videoThumbnailURL = medium.URL
			} else if def, ok := item.Snippet.Thumbnails["default"]; ok {
				videoThumbnailURL = def.URL
			}
		}

		videos[i] = models.VideoResponse{
			ID:           item.Snippet.ResourceID.VideoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			ChannelTitle: item.Snippet.ChannelTitle,
			Position:     item.Snippet.Position,
			AddedAt:      addedAt,
			ThumbnailURL: videoThumbnailURL,
		}
	}

	return videos, nil
}

package services

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alejpaa/playlist-migration-tool/internal/models"
)

// ExportService handles playlist export logic
type ExportService struct {
	playlistService *PlaylistService
}

// NewExportService creates a new ExportService
func NewExportService() *ExportService {
	return &ExportService{
		playlistService: NewPlaylistService(),
	}
}

// ExportPlaylist exports a playlist in the specified format
func (s *ExportService) ExportPlaylist(accessToken, playlistID string, request *models.ExportRequest) (*models.ExportResponse, error) {
	// Get playlist details
	playlist, err := s.playlistService.GetPlaylistByID(accessToken, playlistID)
	if err != nil {
		return nil, err
	}

	var exportData string
	switch request.Format {
	case "json":
		exportData, err = s.exportAsJSON(playlist)
	case "csv":
		exportData, err = s.exportAsCSV(playlist, request.IncludeInfo)
	case "m3u":
		exportData, err = s.exportAsM3U(playlist)
	default:
		return nil, models.NewBadRequestError("Unsupported export format", nil)
	}

	if err != nil {
		return nil, models.NewInternalServerError("Failed to export playlist", err)
	}

	return &models.ExportResponse{
		Success: true,
		Format:  request.Format,
		Data:    exportData,
		Message: fmt.Sprintf("Successfully exported playlist '%s' as %s", playlist.Title, request.Format),
	}, nil
}

// exportAsJSON exports playlist as JSON
func (s *ExportService) exportAsJSON(playlist *models.PlaylistDetailResponse) (string, error) {
	data, err := json.MarshalIndent(playlist, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// exportAsCSV exports playlist as CSV
func (s *ExportService) exportAsCSV(playlist *models.PlaylistDetailResponse, includeInfo bool) (string, error) {
	var buffer strings.Builder
	writer := csv.NewWriter(&buffer)

	// Write header
	if includeInfo {
		header := []string{"Position", "Title", "Channel", "Video ID", "Description", "Added At"}
		if err := writer.Write(header); err != nil {
			return "", err
		}
	} else {
		header := []string{"Position", "Title", "Channel", "Video ID"}
		if err := writer.Write(header); err != nil {
			return "", err
		}
	}

	// Write videos
	for _, video := range playlist.Videos {
		var row []string
		if includeInfo {
			row = []string{
				fmt.Sprintf("%d", video.Position+1),
				video.Title,
				video.ChannelTitle,
				video.ID,
				video.Description,
				video.AddedAt.Format("2006-01-02 15:04:05"),
			}
		} else {
			row = []string{
				fmt.Sprintf("%d", video.Position+1),
				video.Title,
				video.ChannelTitle,
				video.ID,
			}
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

// exportAsM3U exports playlist as M3U format
func (s *ExportService) exportAsM3U(playlist *models.PlaylistDetailResponse) (string, error) {
	var buffer strings.Builder

	// M3U header
	buffer.WriteString("#EXTM3U\n")
	buffer.WriteString(fmt.Sprintf("#PLAYLIST:%s\n", playlist.Title))

	// Add each video
	for _, video := range playlist.Videos {
		buffer.WriteString(fmt.Sprintf("#EXTINF:-1,%s - %s\n", video.ChannelTitle, video.Title))
		buffer.WriteString(fmt.Sprintf("https://www.youtube.com/watch?v=%s\n", video.ID))
	}

	return buffer.String(), nil
}

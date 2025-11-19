package yt

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/api/youtube/v3"
)

// Video represents a video file to be uploaded to YouTube.
type Video struct {
	Title    string // Title is the video title that will appear on YouTube
	FilePath string // FilePath is the full path to the video file on disk
}

// UploadVideo uploads a video to YouTube and adds it to the configured playlist.
// The YouTube channel ID and playlist ID are read from environment variables:
// - MDROID_YOUTUBE_CHANNEL_ID (defaults to hardcoded value if not set)
// - MDROID_YOUTUBE_PLAYLIST_ID (defaults to hardcoded value if not set)
// The video is uploaded as private and marked as not made for kids.
// Returns an error if authentication, upload, or playlist addition fails.
func UploadVideo(video Video) error {
	service, err := getService()
	if err != nil {
		return fmt.Errorf("failed to get YouTube service: %w", err)
	}

	// Get channel ID from environment or use default
	channelID := os.Getenv("MDROID_YOUTUBE_CHANNEL_ID")
	if channelID == "" {
		channelID = "UCGWuMVzuiaJBrE-HwcB7j-w" // Default fallback
	}

	return upload(service, channelID, video)
}

// upload quota is about 1650 (1600 for video.insert and 50 for playlistItems.insert), max quota per day is 10000
// solution is to only upload 6 videos/day, so the backlog will take a few weeks to finish uploading
func upload(service *youtube.Service, channelID string, v Video) error {

	call := service.Videos.Insert([]string{"snippet", "status"}, &youtube.Video{
		Status: &youtube.VideoStatus{
			PrivacyStatus:           "private",
			SelfDeclaredMadeForKids: false,
		},
		Snippet: &youtube.VideoSnippet{
			ChannelId: channelID,
			Title:     v.Title,
		},
	})

	file, err := os.Open(v.FilePath)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", v.FilePath, err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("error in closing file, %v", closeErr)
		}
	}()

	response, err := call.Media(file).Do()
	if err != nil {
		return err
	}

	fmt.Printf("Upload successful! Video ID: %s\n", response.Id)

	// Get playlist ID from environment or use default
	playlistID := os.Getenv("MDROID_YOUTUBE_PLAYLIST_ID")
	if playlistID == "" {
		playlistID = "PL1PP365t1jqj2xSIjh17vIB4tmDfJwJDQ" // Default fallback
	}

	callAddToPlaylist := service.PlaylistItems.Insert([]string{"snippet"}, &youtube.PlaylistItem{Snippet: &youtube.PlaylistItemSnippet{
		PlaylistId: playlistID, Position: 0, ResourceId: &youtube.ResourceId{Kind: "youtube#video", VideoId: response.Id}}})

	respPlaylistInsert, err := callAddToPlaylist.Do()
	if err != nil {
		return err
	}
	fmt.Printf("Video with ID %s added successfully to playlist (%s) Response ID: %s\n", response.Id, playlistID, respPlaylistInsert.Id)
	return nil
}

package yt

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	Title    string
	FilePath string
}

func UploadVideo(video Video) error {
	service, err := getService()
	if err != nil {
		return fmt.Errorf("failed to get YouTube service: %w", err)
	}

	return upload(service, "UCGWuMVzuiaJBrE-HwcB7j-w", video)
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

	backlog := "PL1PP365t1jqj2xSIjh17vIB4tmDfJwJDQ"
	callAddToPlaylist := service.PlaylistItems.Insert([]string{"snippet"}, &youtube.PlaylistItem{Snippet: &youtube.PlaylistItemSnippet{
		PlaylistId: backlog, Position: 0, ResourceId: &youtube.ResourceId{Kind: "youtube#video", VideoId: response.Id}}})

	respPlaylistInsert, err := callAddToPlaylist.Do()
	if err != nil {
		return err
	}
	fmt.Printf("Video with ID %s added successfully to playlist backlog (%s) Response ID: %s\n", response.Id, backlog, respPlaylistInsert.Id)
	return nil
}

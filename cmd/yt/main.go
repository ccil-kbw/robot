package main

import (
	"fmt"
	"log"
	"os"

	"google.golang.org/api/youtube/v3"

	"github.com/ccil-kbw/robot/pkg/helpers"
)

type Video struct {
	Title    string
	FilePath string
}

func main() {
	/*
		uploadVideo(Video{
			Title:    "test",
			FilePath: "/home/seraf/Documents/MEGAKBW/2023-08-10/2023-08-10_05-00-20.mkv",
		})
	*/
	fetchVideos()
}

func fetchVideos() {

	service := getService()

	channel := channelsList(service, []string{"snippet", "contentDetails", "statistics"})
	uploadPlaylistID := channel.ContentDetails.RelatedPlaylists.Uploads
	ids := getVideoIDs(service, channel.Id, uploadPlaylistID)
	getVideos(service, ids)
}

func getVideos(service *youtube.Service, videoIDs []string) {
	// seems like number of ids are limited, do some manual pagination
	offset := 0
	limit := 50

	videos := []*youtube.Video{}
	for {
		call := service.Videos.List([]string{"snippet", "fileDetails"})
		call.Id(videoIDs[offset : offset+limit]...)
		response, err := call.Do()
		helpers.HandleError(err, "could not fetch videos with given ids")

		videos = append(videos, response.Items...)
		offset = offset + limit
		if offset > len(videoIDs) {
			break
		}
	}

	for i, video := range videos {
		fmt.Printf("%d) %s, %s\n", i, video.Snippet.Title, video.FileDetails.FileName)
	}

}

func getVideoIDs(service *youtube.Service, channelID, playlistID string) []string {
	items := []*youtube.PlaylistItem{}
	token := ""
	call := service.PlaylistItems.List([]string{"snippet", "contentDetails", "status"})

	for {
		call = call.MaxResults(50).PlaylistId(playlistID)
		if token != "" {
			call.PageToken(token)
		}
		response, err := call.Do()
		helpers.HandleError(err, "error listing videos")

		items = append(items, response.Items...)

		token = response.NextPageToken
		if token == "" {
			break
		}
	}

	var ids []string
	{
		fmt.Printf("fetching %d video ids from uploaded videos\n", len(items))
		for _, item := range items {
			fmt.Println(item.Snippet.ResourceId.VideoId)
			ids = append(ids, item.Snippet.ResourceId.VideoId)
		}
	}

	return ids
}

func uploadVideo(video Video) {
	service := getService()

	channel := channelsList(service, []string{"snippet", "contentDetails", "statistics"})

	upload(service, channel.Id, video)
}

func channelsList(service *youtube.Service, part []string) *youtube.Channel {
	call := service.Channels.List(part)
	call = call.Mine(true)
	response, err := call.Do()
	helpers.HandleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))

	return response.Items[0]
}

func upload(service *youtube.Service, channelID string, v Video) {

	call := service.Videos.Insert([]string{"snippet", "status"}, &youtube.Video{
		Status: &youtube.VideoStatus{
			PrivacyStatus: "private",
		},
		Snippet: &youtube.VideoSnippet{
			ChannelId: channelID,
			Title:     v.Title,
		},
	})

	file, err := os.Open(v.FilePath)
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Printf("error in closing file, %v", err)
		}
	}()
	if err != nil {
		log.Fatalf("error opening %s: %v", v.FilePath, err)
	}

	response, err := call.Media(file).Do()
	helpers.HandleError(err, "")

	fmt.Printf("Upload successful! Video ID: %s\n", response.Id)
}

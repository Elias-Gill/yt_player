package yt_api

import (
	"context"
	"log"

	"github.com/elias-gill/yt_player/globals"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	TYPE_PLAYLIST = iota
	TYPE_VIDEO
)

const Yt_url = "https://www.youtube.com/watch?v="

type Result struct {
	Title string
	Id    string
	Type  int
}

type Results struct {
	Playlists []Result
	Videos    []Result
}

func RetrieveResults(query string) Results {
	service, err := youtube.NewService(
		context.TODO(),
		option.WithAPIKey(globals.GetApiKey()),
	)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(globals.GetMaxResults())

	response, _ := call.Do()

	// Group video, channel, and playlist results in separate lists.
	videos := []Result{}
	playlists := []Result{}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			videos = append(videos,
				Result{
					Title: item.Snippet.Title,
					Id:    item.Id.VideoId,
					Type:  TYPE_VIDEO,
				})
		}
	}

	return Results{Videos: videos, Playlists: playlists}
}

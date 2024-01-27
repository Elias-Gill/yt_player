package yt_api

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	Title string
	Id    string
}

func RetrieveVideos(query string, maxResults int64, key string) []Video {
	service, err := youtube.NewService(context.TODO(), option.WithAPIKey(key))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(maxResults)

	response, _ := call.Do()

	// Group video, channel, and playlist results in separate lists.
	videos := []Video{}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			videos = append(videos, Video{item.Snippet.Title, item.Id.VideoId})
		}
	}

	return videos
}

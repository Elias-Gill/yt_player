package player

import (
	"context"
	"log"
	"time"

	"github.com/elias-gill/yt_player/settings"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const Yt_url = "https://www.youtube.com/watch?v="

type Player struct {
	settings    *settings.Settings
	mpvInstance *mpvInstance
	ytService   *youtube.Service

	Playlists []Playlist
	Videos    []Video
}

type Playlist struct {
	Title string
	Id    string
}

type Video struct {
	Title string
	Id    string
}

func MustCreatePlayer(settings *settings.Settings) *Player {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	service, err := youtube.NewService(
		ctx, option.WithAPIKey(settings.GetApiKey()),
	)

	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	return &Player{
		mpvInstance: startMpvInstance(),
		settings:    settings,
		ytService:   service,
	}
}

func (p *Player) Search(searchKey string) error {
	// Make the API call to YouTube.
	call := p.ytService.Search.List([]string{"id", "snippet"}).
		Q(searchKey).
		MaxResults(p.settings.GetMaxResults())

	response, err := call.Do()

	if err != nil {
		return err
	}

	// Group video, channel, and playlist results in separate lists.
	videos := []Video{}
	playlists := []Playlist{}

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			videos = append(videos,
				Video{
					Title: item.Snippet.Title,
					Id:    item.Id.VideoId,
				})
		}

		if item.Id.Kind == "youtube#playlist" {
			playlists = append(playlists,
				Playlist{
					Title: item.Snippet.Title,
					Id:    item.Id.VideoId,
				})
		}
	}

	p.Videos = videos
	p.Playlists = playlists

	return nil
}

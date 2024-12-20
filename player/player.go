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

	nextPageToken string
	prevPageToken string
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
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
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

	return p.callApi(call)
}

func (p *Player) GetVideoInfo(index int) (*youtube.VideoSnippet, error) {
	videoCall := p.ytService.Videos.List([]string{"snippet", "contentDetails"}).
		Id(p.Videos[index].Id).MaxResults(1)

	videoResponse, err := videoCall.Do()
	if err != nil {
		return nil, err
	}

	return videoResponse.Items[0].Snippet, nil
}

func (p *Player) NextPage() error {
	// Make the API call to YouTube.
	call := p.ytService.Search.List([]string{"id", "snippet"}).
		PageToken(p.nextPageToken).
		MaxResults(p.settings.GetMaxResults())

	return p.callApi(call)
}

func (p *Player) PrevPage() error {
	// Make the API call to YouTube.
	call := p.ytService.Search.List([]string{"id", "snippet"}).
		PageToken(p.prevPageToken).
		MaxResults(p.settings.GetMaxResults())

	return p.callApi(call)
}

func (p Player) Play(index int) {
	p.mpvInstance.ChangeSong(p.Videos[index].Id)
}

func (p Player) GetStatus() (float64, float64) {
	return p.mpvInstance.GetSongStatus()
}

func (p Player) Deinit() {
	p.mpvInstance.StopPlayer()
}

func (p *Player) callApi(call *youtube.SearchListCall) error {
	response, err := call.Do()
	if err != nil {
		return err
	}

	p.nextPageToken = response.NextPageToken
	p.prevPageToken = response.PrevPageToken

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

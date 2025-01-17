package player

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"time"

	"github.com/blang/mpv"
	"github.com/elias-gill/yt_player/settings"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const (
	maxResults = 80
	Yt_url     = "https://www.youtube.com/watch?v="
)

type Player struct {
	settings    *settings.Settings
	mpvInstance *mpv.Client

	ytService  *youtube.Service
	cancelFunc context.CancelFunc
	serviceCtx context.Context

	Playlists []Playlist
	Videos    []Video

	history []HistoryEntry
}

type Playlist struct {
	Title string
	Id    string
}

type VideoDetails struct {
	Duration    string `json:"duration_string"`
	Author      string `json:"uploader"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Video struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

func MustCreatePlayer(settings *settings.Settings) *Player {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	service, err := youtube.NewService(
		ctx, option.WithAPIKey(settings.GetApiKey()),
	)

	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	var mpvInstance *mpv.Client
	if settings.Tryreattach() {
		mpvInstance = reattachPlayer()
	} else {
		mpvInstance = startMpvInstance()
	}

	player := &Player{
		mpvInstance: mpvInstance,
		settings:    settings,

		ytService:  service,
		serviceCtx: ctx,
		cancelFunc: cancelFunc,
	}

	history := loadHistory()
	if len(history) > 0 {
		lastEntry := history[len(history)-1]
		player.Videos = lastEntry.Videos
		player.Playlists = lastEntry.Playlists
	}
	player.history = history

	return player
}

func (p *Player) Search(searchKey string) error {
	call := p.ytService.Search.List([]string{"id", "snippet"}).
		Q(searchKey).
		MaxResults(maxResults)
	defer p.cancelFunc()

	response, err := call.Do()

	if err != nil {
		return fmt.Errorf("Cannot search for youtube videos")
	}

	// Group video, channel, and playlist results in separate lists.
	videos := []Video{}
	playlists := []Playlist{}

	// To clean emojis and special characters that are cuasing bad rendering on
	// Tmux.

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			videos = append(videos,
				Video{
					Title: removeEmojis(item.Snippet.Title),
					Id:    removeEmojis(item.Id.VideoId),
				})
		}

		if item.Id.Kind == "youtube#playlist" {
			playlists = append(playlists,
				Playlist{
					Title: removeEmojis(item.Snippet.Title),
					Id:    removeEmojis(item.Id.VideoId),
				})
		}
	}

	p.Videos = videos
	p.Playlists = playlists

	p.addHistoryEntry(searchKey, videos, playlists)

	return nil
}

func (p *Player) Play(index int) {
	if index > len(p.Videos)-1 || len(p.Videos) == 0 || index < 0 {
		return
	}

	p.ChangeSong(p.Videos[index].Id)
}

func (p Player) GetStatus() (float64, float64) {
	return p.GetSongStatus()
}

func (p Player) Deinit() {
	p.persistHistory()

	if p.settings.DetachOnQuit() {
		p.DetachPlayer()
	} else {
		p.StopPlayer()
	}
}

// Retrieves the video duration, author, and description
func GetVideoDetails(videoID string) (*VideoDetails, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", "--skip-download", fmt.Sprintf("%s%s", Yt_url, videoID))

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var info VideoDetails
	if err := json.Unmarshal(out.Bytes(), &info); err != nil {
		return nil, err
	}

	info.Description = removeEmojis(info.Description)
	info.Author = removeEmojis(info.Author)
	info.Title = removeEmojis(info.Title)

	return &info, nil
}

func removeEmojis(input string) string {
	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}\x{1F300}-\x{1F5FF}\x{1F680}-\x{1F6FF}\x{1F1E0}-\x{1F1FF}\x{2702}-\x{27B0}\x{24C2}-\x{1F251}]`)
	return emojiRegex.ReplaceAllString(input, "")
}

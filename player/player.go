package player

import (
	"github.com/blang/mpv"
	"github.com/elias-gill/yt_player/history"
	"github.com/elias-gill/yt_player/settings"
	ytservice "github.com/elias-gill/yt_player/yt_service"
)

type Player struct {
	settings    *settings.Settings
	mpvInstance *mpv.Client
	ytService   *ytservice.YtService

	history   history.History
	Playlists []ytservice.Playlist
	Videos    []ytservice.Video
}

func MustCreatePlayer(settings *settings.Settings) *Player {
	var mpvInstance *mpv.Client
	if settings.Tryreattach() {
		mpvInstance = reattachPlayer()
	} else {
		mpvInstance = startMpvInstance()
	}

	player := &Player{
		mpvInstance: mpvInstance,
		settings:    settings,
		ytService:   ytservice.MustCreateYtService(settings),
	}

	history := history.LoadHistory()
	if history.LastSearch != nil {
		lastEntry := history.LastSearch
		player.Videos = lastEntry.Videos
		player.Playlists = lastEntry.Playlists
	}
	player.history = history

	return player
}

func (p *Player) Search(searchKey string) error {
	videos, playlists, err := p.ytService.Search(searchKey)
	if err != nil {
		return err
	}

	p.Videos = videos
	p.Playlists = playlists

	p.history.AddHistoryEntry(searchKey, videos, playlists)

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

func (p Player) GetHistory() history.History {
	return p.history
}

func (p *Player) SelHistoryEntry(index int) {
	item := p.history.SelectEntry(index)
	if item != nil {
		p.Playlists = item.Playlists
		p.Videos = item.Videos
	}
}

func (p Player) Deinit() {
	p.history.PersistHistory()

	if p.settings.DetachOnQuit() {
		p.DetachPlayer()
	} else {
		p.StopPlayer()
	}
}

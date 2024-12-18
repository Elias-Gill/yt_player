package context

import (
	"github.com/elias-gill/yt_player/completition"
	"github.com/elias-gill/yt_player/player"
	"github.com/elias-gill/yt_player/settings"
)

type Context struct {
	Player    *player.Player
	Config    *settings.Settings
	History   *completition.Completition
	CurrMode  Mode
	WinHeight int
	WinWidth  int
	Error     error
}

type Mode int

const (
	SEARCH = iota
	LIST
	HELP
)

func MustLoadContext() *Context {
	config := settings.MustParseConfig()

	return &Context{
		Config:   config,
		Player:   player.MustCreatePlayer(config),
		History:  completition.LoadHistory(),
		CurrMode: SEARCH,
	}
}

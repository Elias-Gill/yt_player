package context

import (
	// "github.com/elias-gill/yt_player/completition"
	"github.com/elias-gill/yt_player/completition"
	"github.com/elias-gill/yt_player/mpv"
	"github.com/elias-gill/yt_player/settings"
)

type Context struct {
	Player    *mpv.MpvPlayer
	Config    *settings.Settings
	History   *completition.Completition
	CurrMode  Mode
	WinHeight int
	WinWidth  int
}

type Mode int

const (
	SEARCH = iota
	LIST
	HELP
)

func NewContext() *Context {
	return &Context{
		Config:   settings.MustParseConfig(),
		Player:   mpv.MustInitPlayer(),
		History:  completition.LoadHistory(),
		CurrMode: SEARCH,
	}
}

package context

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/completition"
	"github.com/elias-gill/yt_player/player"
	"github.com/elias-gill/yt_player/settings"
)

type Mode int

const (
	SEARCH = iota
	LIST
)

type styles struct {
	// Global Colors
	Background lipgloss.Style

	// Foreground focus Colors
	ForegroundRed  lipgloss.Style
	ForegroundAqua lipgloss.Style
	ForegroundGray lipgloss.Style

	BackgroundRed  lipgloss.Style
	BackgroundGray lipgloss.Style
}

type Context struct {
	// Instances
	Player  *player.Player
	Config  *settings.Settings
	History *completition.Completition
	Styles  styles

	// App Status
	CurrMode  Mode
	CurrItem  int
	WinHeight int
	WinWidth  int
	Error     error
}

func (c *Context) NextMode() {
	if c.CurrMode == LIST {
		c.CurrMode = SEARCH
	} else {
		c.CurrMode++
	}
}

func (c *Context) Deinit() {
	c.Player.Deinit()
	c.History.Deinit()
}

func MustLoadContext() *Context {
	config := settings.MustParseConfig()

	return &Context{
		Config:   config,
		Player:   player.MustCreatePlayer(config),
		History:  completition.LoadHistory(),
		CurrMode: SEARCH,
		CurrItem: 0,

		Styles: styles{
			Background: lipgloss.NewStyle().
				Background(lipgloss.Color("#1d2021")).
				Padding(1),
			BackgroundRed: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282828")).
				Background(lipgloss.Color("#f96c5b")),
			BackgroundGray: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282828")).
				Background(lipgloss.Color("#928374")),
			ForegroundRed:  lipgloss.NewStyle().Foreground(lipgloss.Color("#f96c5b")),
			ForegroundAqua: lipgloss.NewStyle().Foreground(lipgloss.Color("#8ec07c")),
			ForegroundGray: lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")),
		},
	}
}

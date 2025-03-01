package context

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/player"
	"github.com/elias-gill/yt_player/settings"
)

type Mode int

const (
	SEARCH = iota
	LIST
	HISTORY
	HELP
)

type styles struct {
	// Global background color
	Background lipgloss.Style

	// Foreground focus Colors
	ForegroundRed  lipgloss.Style
	ForegroundAqua lipgloss.Style
	ForegroundGray lipgloss.Style

	// Background focus Colors
	BackgroundRed  lipgloss.Style
	BackgroundGray lipgloss.Style
}

type Context struct {
	// Instances
	Player *player.Player
	Config *settings.Settings
	Styles styles

	CurrMode Mode
	Error    error

	WinHeight int
	WinWidth  int
}

func (c *Context) Deinit() {
	c.Player.Deinit()
}

func MustLoadContext() *Context {
	config := settings.MustParseConfig()

	return &Context{
		Config:   config,
		Player:   player.MustCreatePlayer(config),
		CurrMode: SEARCH,

		Styles: styles{
			Background: lipgloss.NewStyle().
				Background(lipgloss.Color("#1d2021")),
			BackgroundRed: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282828")).
				Background(lipgloss.Color("#f96c5b")),
			BackgroundGray: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282828")).
				Background(lipgloss.Color("#928374")),
			ForegroundRed: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f96c5b")).
				Background(lipgloss.Color("#1d2021")),
			ForegroundAqua: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8ec07c")).
				Background(lipgloss.Color("#1d2021")),
			ForegroundGray: lipgloss.NewStyle().
				Foreground(lipgloss.Color("#928374")).
				Background(lipgloss.Color("#1d2021")),
		},
	}
}

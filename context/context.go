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
	GruvboxBg     lipgloss.Style
	GruvboxRed    lipgloss.Style
	GruvboxGreen  lipgloss.Style
	GruvboxYellow lipgloss.Style
	GruvboxBlue   lipgloss.Style
	GruvboxPurple lipgloss.Style
	GruvboxAqua   lipgloss.Style
	GruvboxOrange lipgloss.Style
	GruvboxGray   lipgloss.Style
	GruvboxBlack  lipgloss.Style
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
			GruvboxBg:     lipgloss.NewStyle().Background(lipgloss.Color("#1d2021")).Padding(1), // Background
			GruvboxRed:    lipgloss.NewStyle().Foreground(lipgloss.Color("#fb4934")),            // Red
			GruvboxGreen:  lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26")),            // Green
			GruvboxYellow: lipgloss.NewStyle().Foreground(lipgloss.Color("#fabd2f")),            // Yellow
			GruvboxBlue:   lipgloss.NewStyle().Foreground(lipgloss.Color("#83a598")),            // Blue
			GruvboxPurple: lipgloss.NewStyle().Foreground(lipgloss.Color("#d3869b")),            // Purple
			GruvboxAqua:   lipgloss.NewStyle().Foreground(lipgloss.Color("#8ec07c")),            // Aqua
			GruvboxOrange: lipgloss.NewStyle().Foreground(lipgloss.Color("#fe8019")),            // Orange
			GruvboxGray:   lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")),            // Gray
			GruvboxBlack:  lipgloss.NewStyle().Foreground(lipgloss.Color("#282828")),            // Black
		},
	}
}

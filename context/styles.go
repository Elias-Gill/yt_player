package context

import "github.com/charmbracelet/lipgloss"

type styles struct {
	Background     lipgloss.Style
	BackgroundRed  lipgloss.Style
	BackgroundGray lipgloss.Style
	ForegroundRed  lipgloss.Style
	ForegroundAqua lipgloss.Style
	ForegroundGray lipgloss.Style
}

var Styles = newStyles()

func newStyles() *styles {
	return &styles{
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
	}
}

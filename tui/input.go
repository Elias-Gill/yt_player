package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
)

func initInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Your music"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80

	return ti
}

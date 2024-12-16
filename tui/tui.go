package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/context"
	"github.com/elias-gill/yt_player/tui/components"
)

type Tui struct {
	context   *context.Context
	textInput components.Input
}

func NewModel(ctx *context.Context) tea.Model {
	return Tui{
		context:   ctx,
		textInput: components.NewInput(ctx),
	}
}

func (t Tui) View() string {
	return t.textInput.View()
}

func (t Tui) Init() tea.Cmd {
	return nil
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if tea.KeyCtrlC == msg.Type {
			return t, tea.Quit
		}
	case tea.WindowSizeMsg:
		t.context.WinHeight = msg.Height
		t.context.WinWidth = msg.Width
	}

	var cmd tea.Cmd
	t.textInput, cmd = t.textInput.Update(msg)
	return t, cmd
}

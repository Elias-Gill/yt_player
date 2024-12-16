package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/context"
)

type Input struct {
	ctx       *context.Context
	textArea  textinput.Model
	inputHelp help.Model
	prompt    string
}

func NewInput(ctx *context.Context) Input {
	ta := textinput.New()
	ta.Prompt = "Search> "
	ta.Focus()

	return Input{
		ctx:      ctx,
		textArea: ta,
	}
}

func (i Input) View() string {
	// TODO: wrap into some styling with lipgloss
	return i.textArea.View()
}

func (i Input) Init() tea.Cmd {
	return nil
}

func (in Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	var cmd tea.Cmd
	in.textArea, cmd = in.textArea.Update(msg)
	return in, cmd
}

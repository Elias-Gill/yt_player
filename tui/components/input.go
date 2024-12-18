package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
)

type Input struct {
	ctx       *context.Context
	textInput textinput.Model
	inputHelp help.Model
	prompt    string
}

func NewInput(ctx *context.Context) Input {
	ti := textinput.New()
	ti.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(context.GruvboxAqua)).Render(" Search> ")
	ti.Focus()

	return Input{
		ctx:       ctx,
		textInput: ti,
	}
}

func (i Input) View() string {
	return i.textInput.View()
}

func (i Input) Init() tea.Cmd {
	return nil
}

func (in Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		in.textInput.Width = int(float32(msg.Width-len(in.textInput.Prompt)) * 0.7)
	}

	var cmd tea.Cmd
	in.textInput, cmd = in.textInput.Update(msg)

	return in, cmd
}

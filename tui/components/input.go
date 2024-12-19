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
	if i.ctx.CurrMode != context.SEARCH {
		i.textInput.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(context.GruvboxGray)).Render(" Search> ")
	}

	return i.textInput.View()
}

func (i Input) Init() tea.Cmd {
	return nil
}

func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.textInput.Width = int(float32(msg.Width-len(i.textInput.Prompt)) * 0.7)
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			i.ctx.CurrMode = context.LIST
			i.ctx.Player.Search(i.textInput.Value())

			return i, nil
		}
	}

	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)

	return i, cmd
}

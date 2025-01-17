package components

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/context"
)

type Input struct {
	ctx       *context.Context
	textInput textinput.Model
	inputHelp help.Model
	prompt    string
}

const prompt = " Search> "

func NewInput(ctx *context.Context) Input {
	ti := textinput.New()
	ti.Prompt = ctx.Styles.BackgroundRed.Render(prompt) + " "
	ti.Focus()

	return Input{
		ctx:       ctx,
		textInput: ti,
	}
}

func (i Input) View() string {
	if i.ctx.CurrMode != context.SEARCH {
		i.textInput.Prompt = i.ctx.Styles.BackgroundGray.Render(prompt) + " "
	}

	return i.textInput.View()
}

func (i Input) Init() tea.Cmd {
	return nil
}

func (i Input) Update(msg tea.Msg) (Input, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		i.textInput.Width = msg.Width - len(i.textInput.Prompt)
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			i.ctx.CurrMode = context.LIST
			err := i.ctx.Player.Search(i.textInput.Value())
			i.ctx.Error = err
			return i, nil
		}
	}

	var cmd tea.Cmd
	i.textInput, cmd = i.textInput.Update(msg)

	return i, cmd
}

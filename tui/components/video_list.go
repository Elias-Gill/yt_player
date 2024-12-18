package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
)

type VideoList struct {
	ctx *context.Context
}

func NewVideoList(ctx *context.Context) VideoList {
	ti := textinput.New()
	ti.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(context.GruvboxAqua)).Render(" Search> ")
	ti.Focus()

	return VideoList{
		ctx: ctx,
	}
}

func (v VideoList) View() string {
	return ""
}

func (v VideoList) Init() tea.Cmd {
	return nil
}

func (v VideoList) Update(msg tea.Msg) (VideoList, tea.Cmd) {
	return v, nil
}

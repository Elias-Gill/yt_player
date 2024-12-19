package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/elias-gill/yt_player/context"
)

type tickMsg time.Time

type PlayerProgress struct {
	ctx      *context.Context
	percent  float64
	progress progress.Model
}

func NewPlayerInfo(ctx *context.Context) PlayerProgress {
	prg := progress.New()
	prg.ShowPercentage = false

	return PlayerProgress{
		ctx:      ctx,
		progress: prg,
		percent:  0,
	}
}

func (p PlayerProgress) Update(msg tea.Msg) (PlayerProgress, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.progress.Width = msg.Width
		return p, nil

	case tickMsg:
		currPos, duration := p.ctx.Player.GetStatus()
		p.percent = currPos / duration
		return p, tickCmd()
	}

	return p, nil
}

func (p PlayerProgress) View() string {
	currPos, duration := p.ctx.Player.GetStatus()
	currPosDuration := time.Duration(currPos) * time.Second
	durationDuration := time.Duration(duration) * time.Second
	playerTime := fmt.Sprintf("%s / %s", currPosDuration.String(), durationDuration.String())

	return lipgloss.JoinVertical(0, p.progress.ViewAs(p.percent), playerTime)
}

func (p PlayerProgress) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
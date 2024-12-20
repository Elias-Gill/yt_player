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
	prg.Width = 80

	return PlayerProgress{
		ctx:      ctx,
		progress: prg,
		percent:  0,
	}
}

func (p PlayerProgress) Update(msg tea.Msg) (PlayerProgress, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.progress.Width = msg.Width - 2
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
	currPosDuration := time.Duration(currPos).String()
	durationDuration := time.Duration(duration).String()

	hPrompt := "Help: '?'"
	help := p.ctx.Styles.ForegroundGray.Inherit(p.ctx.Styles.Background).Render(hPrompt)

	playerTime := p.ctx.Styles.ForegroundAqua.
		Inherit(p.ctx.Styles.Background).
		Width(p.ctx.WinWidth - len(hPrompt) - 2).
		Render(fmt.Sprintf("%s / %s  %s",
			currPosDuration,
			durationDuration,
			p.ctx.Player.GetCurrentSong()))

	return lipgloss.JoinVertical(
		0,
		p.progress.ViewAs(p.percent),
		lipgloss.JoinHorizontal(0, playerTime, help),
	)
}

func (p PlayerProgress) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

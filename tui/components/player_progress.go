package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/context"
)

type tickMsg time.Time

type PlayerProgress struct {
	ctx      *context.Context
	percent  float64
	progress progress.Model

	songDuration string
	currPosition string
}

func NewPlayer(ctx *context.Context) PlayerProgress {
	prg := progress.New()
	prg.ShowPercentage = false
	prg.Width = 80

	return PlayerProgress{
		ctx:      ctx,
		progress: prg,
		percent:  0,

		currPosition: "0s",
		songDuration: "0s",
	}
}

func (p PlayerProgress) Update(msg tea.Msg) (PlayerProgress, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.progress.Width = msg.Width - 2 // minus padding
		return p, nil

	case tickMsg:
		currPos, duration := p.ctx.Player.GetStatus()

		p.percent = currPos / duration
		p.currPosition = time.Duration(currPos).String()
		p.songDuration = time.Duration(duration).String()

		return p, tickCmd()
	}

	return p, nil
}

func (p PlayerProgress) View() string {
	hPrompt := "Help: '~'"
	help := p.ctx.Styles.ForegroundGray.Inherit(p.ctx.Styles.Background).Render(hPrompt)

	style := p.ctx.Styles.ForegroundAqua
	if p.ctx.Player.IsPaused() {
		style = p.ctx.Styles.ForegroundGray
	}

	playerTime := style.
		Width(p.ctx.WinWidth - len(hPrompt) - 2).
		Render(fmt.Sprintf("%s / %s  %s",
			p.currPosition,
			p.songDuration,
			p.ctx.Player.GetCurrentSong()))

	return fmt.Sprintf("%s\n%s",
		p.progress.ViewAs(p.percent),
		playerTime+help)
}

func (p PlayerProgress) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

package components

import (
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
	appCtx "github.com/elias-gill/yt_player/context"
	player "github.com/elias-gill/yt_player/player"
)

type VideoInfo struct {
	ctx     *appCtx.Context
	details *player.VideoDetails
	errors  error
}

func NewVideoInfo(ctx *appCtx.Context) VideoInfo {
	return VideoInfo{
		ctx: ctx,
	}
}

func (v VideoInfo) Update(id string) VideoInfo {
	details, err := player.GetVideoDetails(id)
	v.errors = err
	v.details = details

	return v
}

func (v VideoInfo) View() string {
	// Set a rounded, yellow-on-purple border to the top and left
	var anotherStyle = v.ctx.Styles.Background.
		PaddingTop(0).
		Width(int(math.Round(float64(v.ctx.WinWidth) * 0.30)) - 1).
		MaxWidth(int(math.Round(float64(v.ctx.WinWidth) * 0.30)) - 1).
		Height(v.ctx.WinHeight - 9).
		MaxHeight(v.ctx.WinHeight - 9).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#999999")).
		BorderLeft(true).
		BorderTop(true)

	if v.details == nil {
		if v.errors != nil {
			return anotherStyle.Render("Cannot retrieve video details\n" + v.errors.Error())
		}

		return ""
	}

	return anotherStyle.Render(
		fmt.Sprintf("Author: %s\nDuration: %s\nTitle: %s\nDescription: %s",
			v.details.Author,
			v.details.Duration,
			v.details.Title,
			v.details.Description))
}

func (v VideoInfo) Init() tea.Cmd {
	return nil
}

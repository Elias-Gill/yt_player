package main

import (
	"fmt"
	"os"

	"github.com/elias-gill/yt_player/completition"
	"github.com/elias-gill/yt_player/settings"
	"github.com/elias-gill/yt_player/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/elias-gill/yt_player/mpv"
)

func main() {
	settings.ParseConfig()
	mpv.StartPlayer()
	completition.LoadHistory()

	if _, err := tea.NewProgram(tui.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

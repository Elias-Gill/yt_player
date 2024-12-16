package main

import (
	"fmt"
	"os"

	"github.com/elias-gill/yt_player/context"
	"github.com/elias-gill/yt_player/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := tui.NewModel(context.NewContext())

	if _, err := tea.NewProgram(
		model,
		tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

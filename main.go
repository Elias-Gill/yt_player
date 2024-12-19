package main

import (
	"fmt"
	"os"

	appCtx "github.com/elias-gill/yt_player/context"
	"github.com/elias-gill/yt_player/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ctx := appCtx.MustLoadContext()
	defer ctx.Deinit()

	model := tui.NewModel(ctx)

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

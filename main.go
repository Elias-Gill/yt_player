package main

import (
	"github.com/elias-gill/yt_player/globals"
	"github.com/elias-gill/yt_player/tui"
)

func main() {
    globals.ParseFlags()
	tui.InitTUI()
}

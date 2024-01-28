package yt_api

import (
	"testing"

	"github.com/elias-gill/yt_player/globals"
)

func TestYTConnection(t *testing.T) {
	globals.ParseFlags()
	res := RetrieveResults("")

	if len(res.Videos) == 0 {
		t.Fatalf("Video list is empty")
	}
}

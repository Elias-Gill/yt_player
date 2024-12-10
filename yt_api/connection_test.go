package yt_api

import (
	"testing"

	"github.com/elias-gill/yt_player/settings"
)

func TestYTConnection(t *testing.T) {
	settings.ParseConfig()
	res := RetrieveResults("")

	if len(res.Videos) == 0 {
		t.Fatalf("Video list is empty")
	}
}

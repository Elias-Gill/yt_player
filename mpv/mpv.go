// Wrapper functions for mpv player
package mpv

import (
	"fmt"
	"os"
	"os/exec"
)

const yt_link = "https://www.youtube.com/watch?v="

func NewPlayer(ytId string) *exec.Cmd {
	url := yt_link + ytId

	cmd := exec.Command("mpv", "--no-video", url)
	if cmd == nil {
		fmt.Println("Cannot find 'mpv' player on your machine.")
		os.Exit(1)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error running mpv: ", err.Error())
		os.Exit(1)
	}

	return cmd
}

func KillInstance(instance *exec.Cmd) {
	// TODO: fix windows kill error
	if instance != nil {
		instance.Process.Kill()
	}
}

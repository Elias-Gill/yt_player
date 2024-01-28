// Wrapper functions for mpv player
package mpv

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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
	// TEST: this should be tested on a windows environment
	if runtime.GOOS == "windows" {
		err := exec.Command("taskkill.exe", "/im", "mpv.exe", "/f").Run()
		if err != nil {
			// TODO: make a log system
		}
		return
	}

	if instance != nil {
		instance.Process.Kill()
	}
}

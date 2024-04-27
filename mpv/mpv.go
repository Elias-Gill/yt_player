// Wrapper functions for mpv player
package mpv

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/blang/mpv"
)

var player *mpv.Client
var cmd *exec.Cmd

func StartPlayer() *exec.Cmd {
	if !commandExists("mpv") {
		fmt.Println("Cannot find mpv player")
		os.Exit(1)
	}

	if !commandExists("youtube-dl") {
		fmt.Println("Cannot find youtube-dl executable")
		os.Exit(1)
	}

	cmd = exec.Command("mpv", "--idle=yes", "--input-ipc-server=/tmp/mpvsocket", "--no-video")

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting mpv player: ", err.Error())
		os.Exit(1)
	}

	// Little hacky but anyways
	time.Sleep(time.Second)

	ipc := mpv.NewIPCClient("/tmp/mpvsocket")
	player = mpv.NewClient(ipc) // Lowlevel client

	return cmd
}

func StopPlayer() {
	cmd.Process.Kill()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)

	return err == nil
}

func ChangeSong(url string) {
	err := player.Loadfile(url, mpv.LoadFileModeReplace)
	if err != nil {
		panic(err)
	}
}

func TogglePause() {
	p, _ := player.Pause()
	player.SetPause(!p)
}

func PlusFiveSecs() {
	curTime, _ := player.GetFloatProperty("time-pos")
	newTime := string(strconv.Itoa(int(curTime + 5)))
	player.SetProperty("time-pos", newTime)
}

func LessFiveSecs() {
	curTime, _ := player.GetFloatProperty("time-pos")
	newTime := string(strconv.Itoa(int(curTime - 5)))
	player.SetProperty("time-pos", newTime)
}

func GetSongLength() int {
	property, _ := player.GetProperty("duration")
	duration, _ := strconv.Atoi(property)

	return duration
}
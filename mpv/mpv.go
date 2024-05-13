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

const mpvSocket = "/tmp/mpvsocket"

func StartPlayer() *exec.Cmd {
	if !commandExists("mpv") {
		fmt.Println("Cannot find mpv player")
		os.Exit(1)
	}

	if !commandExists("youtube-dl") {
		fmt.Println("Cannot find youtube-dl executable")
		os.Exit(1)
	}

	cmd = exec.Command("mpv", "--idle=yes", "--input-ipc-server="+mpvSocket, "--no-video")

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting mpv player: ", err.Error())
		os.Exit(1)
	}

	// wait for the mpv process to start. A little hacky but anyways
	time.Sleep(time.Second)

	ipc := mpv.NewIPCClient("/tmp/mpvsocket")
	player = mpv.NewClient(ipc)

	return cmd
}

func StopPlayer() {
	cmd.Process.Kill()
}

// detach the player but stop idle mode, so the mpv process
// would stop after the song is over
func DetachPlayer() {
	player.SetProperty("idle", "no")
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

func GetSongStatus() string {
	duration, _ := player.GetFloatProperty("duration")
	curPos, _ := player.GetFloatProperty("time-pos")

	return time.Duration(curPos*1e9).String() + " / " + time.Duration(duration*1e9).String()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)

	return err == nil
}

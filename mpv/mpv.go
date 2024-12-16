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

const mpvSocket = "/tmp/mpvsocket"

type MpvPlayer struct {
	player *mpv.Client
	cmd    *exec.Cmd
}

// Generates a new instance of the MpvPlayer cmd. Panics if MPV or youtube-dl executables cannot
// be located in the path, or if the socket connection with MPV fails.
func MustInitPlayer() *MpvPlayer {
	if !commandExists("mpv") {
		fmt.Println("Cannot find mpv player")
		os.Exit(1)
	}

	if !commandExists("youtube-dl") {
		fmt.Println("Cannot find youtube-dl executable")
		os.Exit(1)
	}

	cmd := exec.Command("mpv", "--idle=yes", "--input-ipc-server="+mpvSocket, "--no-video")

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting mpv player: ", err.Error())
		os.Exit(1)
	}

	// wait for the mpv process to start. A little hacky but anyways
	time.Sleep(time.Second)

	ipc := mpv.NewIPCClient("/tmp/mpvsocket")
	player := mpv.NewClient(ipc)

	return &MpvPlayer{
		player: player,
		cmd:    cmd,
	}
}

func (m MpvPlayer) StopPlayer() {
	m.cmd.Process.Kill()
}

// detach the player but stop idle mode, so the mpv process
// would stop after the song is over
func (m MpvPlayer) DetachPlayer() {
	m.player.SetProperty("idle", "no")
}

func (m MpvPlayer) ChangeSong(url string) {
	err := m.player.Loadfile(url, mpv.LoadFileModeReplace)
	if err != nil {
		panic(err)
	}
}

func (m MpvPlayer) TogglePause() {
	p, _ := m.player.Pause()
	m.player.SetPause(!p)
}

func (m MpvPlayer) PlusFiveSecs() {
	curTime, _ := m.player.GetFloatProperty("time-pos")
	newTime := string(strconv.Itoa(int(curTime + 5)))
	m.player.SetProperty("time-pos", newTime)
}

func (m MpvPlayer) LessFiveSecs() {
	curTime, _ := m.player.GetFloatProperty("time-pos")
	newTime := string(strconv.Itoa(int(curTime - 5)))
	m.player.SetProperty("time-pos", newTime)
}

func (m MpvPlayer) GetSongStatus() string {
	duration, _ := m.player.GetFloatProperty("duration")
	curPos, _ := m.player.GetFloatProperty("time-pos")

	return time.Duration(curPos*1e9).String() + " / " + time.Duration(duration*1e9).String()
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// Wrapper functions for mpv player
package player

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/blang/mpv"
	ytservice "github.com/elias-gill/yt_player/yt_service"
)

const mpvSocket = "/tmp/mpvsocket"

// Tries to reattach to a previous detached mpv player instance
func reattachPlayer() *mpv.Client {
	// Safely create the IPC connection and instantiate the client
	safeNewClient := func(socket string) (*mpv.Client, error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Cannot reattach to player instance: \n\t%v\n\nStarting a new MPV instance", r)
			}
		}()
		ipc := mpv.NewIPCClient(socket)
		client := mpv.NewClient(ipc)

		return client, nil
	}

	player, err := safeNewClient("/tmp/mpvsocket")
	if err != nil || player == nil {
		return startMpvInstance()
	}

	player.SetProperty("idle", "yes")

	return player
}

// Generates a new instance of the MpvPlayer cmd. Panics if MPV or youtube-dl executables cannot
// be located in the path, or if the socket connection with MPV fails.
func startMpvInstance() *mpv.Client {
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

	return player
}

func (m Player) StopPlayer() {
	m.mpvInstance.Exec("quit")
}

// Detach the player and stop idle mode so the mpv process
// would stop after the song is over
func (m Player) DetachPlayer() {
	m.mpvInstance.SetProperty("idle", "no")
}

func (m Player) ChangeSong(id string) {
	err := m.mpvInstance.Loadfile(ytservice.Yt_url+id, mpv.LoadFileModeReplace)
	if err != nil {
		panic(err)
	}
}

func (m Player) TogglePause() {
	p, _ := m.mpvInstance.Pause()
	m.mpvInstance.SetPause(!p)
}

func (m Player) PlusFiveSecs() {
	curTime, _ := m.mpvInstance.GetFloatProperty("time-pos")
	newTime := string(strconv.Itoa(int(curTime + 5)))
	m.mpvInstance.SetProperty("time-pos", newTime)
}

func (m Player) LessFiveSecs() {
	curTime, _ := m.mpvInstance.GetFloatProperty("time-pos")
	newTime := string(strconv.Itoa(int(curTime - 5)))
	m.mpvInstance.SetProperty("time-pos", newTime)
}

// Returns de lenght of the current position and the length of the song (position, length)
func (m Player) GetSongStatus() (float64, float64) {
	duration, _ := m.mpvInstance.GetFloatProperty("duration")
	curPos, _ := m.mpvInstance.GetFloatProperty("time-pos")

	return curPos * 1e9, duration * 1e9
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (p Player) IsPaused() bool {
	s, _ := p.mpvInstance.Pause()
	return s
}

func (p Player) GetCurrentSong() string {
	title, err := p.mpvInstance.GetProperty("media-title")
	if err != nil {
		panic("Mpv connection failed: \n" + err.Error())
	}

	return title
}

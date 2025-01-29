// Wrapper functions for mpv player
package player

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/blang/mpv"
	ytservice "github.com/elias-gill/yt_player/yt_service"
)

const mpvSocket = "/tmp/mpvsocket"

// Tries to reattach to a previous detached mpv player instance
func reattachPlayer() *mpv.Client {
	// Safely create the IPC connection and instantiate the client intercepting the mpv library
	// panic.
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
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // Create a new process group to detach tmux
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting mpv: ", err.Error())
		os.Exit(1)
	}

	// Create the socket connection
	ipc := mpv.NewIPCClient("/tmp/mpvsocket")

	// Wait for MPV IPC socket to be available (retry mechanism)
	for i := 0; i < 10; i++ { // Retry up to 10 times
		if _, err := os.Stat(mpvSocket); err == nil {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Connect to the mpv instance
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

func (m *Player) PlusFiveSecs() {
	curTime, _ := m.mpvInstance.GetFloatProperty("time-pos")
	m.mpvInstance.SetProperty("time-pos", fmt.Sprintf("%f", curTime+5))
}

func (m *Player) LessFiveSecs() {
	curTime, _ := m.mpvInstance.GetFloatProperty("time-pos")
	m.mpvInstance.SetProperty("time-pos", fmt.Sprintf("%f", curTime-5))
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

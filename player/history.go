package player

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const maxHistorySize = 50
const historyFile = "yt_history"

type HistoryEntry struct {
	Input     string     `json:"input"`
	Videos    []Video    `json:"videos"`
	Playlists []Playlist `json:"playlists"`
}

// Loads the history file into memory. If the history file is not available, the history
// is still managed in volatile memory.
func loadHistory() []HistoryEntry {
	cache, err := os.UserCacheDir()
	if err != nil {
		return []HistoryEntry{}
	}

	filePath := filepath.Join(cache, historyFile)

	content, err := os.Open(filePath)
	if err != nil {
		return []HistoryEntry{}
	}

	var hist []HistoryEntry
	err = json.NewDecoder(content).Decode(&hist)
	if err != nil {
		return []HistoryEntry{}
	}

	return hist
}

func (p *Player) selectEntry(index int) {
	if index > len(p.history)-1 || index < 0 {
		return
	}

	entry := p.history[index]

	p.Playlists = entry.Playlists
	p.Videos = entry.Videos
}

// Adds a new entry to the search history. If the maxHistorySize is reached, then replaces
// the LRU element in the completition history.
func (p *Player) addHistoryEntry(input string, videos []Video, playlists []Playlist) {
	size := len(p.history)
	entry := HistoryEntry{Input: input, Videos: videos, Playlists: playlists}

	if size == maxHistorySize {
		p.history = p.history[1 : size-1]
		p.history = append(p.history, entry)
	} else {
		p.history = append(p.history, entry)
	}
}

// Tries to persist the current history. If the file cannot be accesed returns an error
func (c Player) persistHistory() error {
	cache, err := os.UserCacheDir()
	if err != nil {
		return err
	}

	filePath := filepath.Join(cache, historyFile)

	// Create the directory if it does not exist
	dir := filepath.Dir(filePath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	// Open the file in truncation mode (create if not exists)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	json.NewEncoder(file).Encode(c.history)
	return nil
}

package history

import (
	"encoding/json"
	"os"
	"path/filepath"

	ytservice "github.com/elias-gill/yt_player/yt_service"
)

const maxHistorySize = 20
const historyFile = "yt_history"

type HistoryEntry struct {
	Input     string               `json:"input"`
	Videos    []ytservice.Video    `json:"videos"`
	Playlists []ytservice.Playlist `json:"playlists"`
}

type History struct {
	History    []HistoryEntry   `json:"history"`
	LastSearch *HistoryEntry    `json:"last_search"`
	LastSong   *ytservice.Video `json:"last_song"`
}

// Loads the history file into memory. If the history file is not available, the history
// is still managed in volatile memory.
func LoadHistory() History {
	cache, err := os.UserCacheDir()
	if err != nil {
		return History{}
	}

	filePath := filepath.Join(cache, historyFile)

	content, err := os.Open(filePath)
	if err != nil {
		return History{}
	}

	var hist History
	err = json.NewDecoder(content).Decode(&hist)
	if err != nil {
		return History{}
	}

	return hist
}

func (h *History) SelectEntry(index int) *HistoryEntry {
	if index > len(h.History)-1 || index < 0 {
		return nil
	}

	// Invert the selection index because the `getHistory` function returns the entry list
	// in reverse order. Calculate the inverted index as `history.length - 1 - index`.
	entry := h.History[len(h.History)-1-index]
	h.LastSearch = &entry

	return &entry
}

// Adds a new entry to the search history. If the maxHistorySize is reached, then replaces
// the LRU element in the completition history.
func (h *History) AddHistoryEntry(input string, videos []ytservice.Video, playlists []ytservice.Playlist) {
	size := len(h.History)
	entry := HistoryEntry{Input: input, Videos: videos, Playlists: playlists}
	h.LastSearch = &entry

	if size == maxHistorySize {
		h.History = h.History[1 : size-1]
		h.History = append(h.History, entry)
	} else {
		h.History = append(h.History, entry)
	}
}

func (p History) GetHistoryList() []HistoryEntry {
	if len(p.History) == 0 {
		return p.History
	}

	// reverse history
	var aux []HistoryEntry
	for i := len(p.History) - 1; i >= 0; i-- {
		aux = append(aux, p.History[i])
	}

	return aux
}

// Tries to persist the current history. If the file cannot be accesed returns an error
func (h History) PersistHistory() error {
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

	json.NewEncoder(file).Encode(h)
	return nil
}

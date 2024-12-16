package completition

import (
	"bufio"
	"fmt"
	"os"
	"path"
)

const maxHistorySize = 50

type Completition struct {
	history     []string
	historyFd   *os.File
	historyPath string

	currItem int
}

// Loads the history file into memory. If the history file is not available, the history
// is still managed in volatile memory.
func LoadHistory() *Completition {
	completition := &Completition{
		history:     []string{},
		currItem:    0,
		historyFd:   nil,
		historyPath: "",
	}

	// If the history file cannot be loaded, then return a volatile completition
	cachePath, err := os.UserCacheDir()
	if err != nil {
		return completition
	}

	hp := path.Join(cachePath, "yt_player.history")
	historyFile, err := os.Open(hp)

	// Creates the file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(hp)
		if err != nil {
			return completition
		}
		historyFile = file
	}

	completition.historyFd = historyFile
	completition.historyPath = hp

	var history []string
	scanner := bufio.NewScanner(historyFile)
	for i := 0; scanner.Scan() && i < maxHistorySize; i++ {
		history = append(history, scanner.Text())
	}

	completition.currItem = len(history) - 1

	return completition
}

// Adds a new entry to the search history. If the maxHistorySize is reached, then replaces
// the LRU element in the completition history.
func (c *Completition) AddHistoryEntry(entry string) {
	size := len(c.history)
	// LRU history replacement
	if size == maxHistorySize {
		c.history = c.history[1 : size-1]
		c.history = append(c.history, entry)
		c.currItem = len(c.history) - 1
	} else {
		c.history = append(c.history, entry)
	}
}

// Retrieves and older entry from the history list
func (c *Completition) NextEntry() (string, error) {
	if len(c.history) == 0 {
		return "", fmt.Errorf("The history list is empty")
	}

	if c.currItem > 0 {
		c.currItem--
	}

	return c.history[c.currItem], nil
}

// Retrieves an newer entry on the history list
func (c *Completition) PrevEntry() (string, error) {
	if len(c.history) == 0 {
		return "", fmt.Errorf("The history list is empty")
	}

	if c.currItem+1 < len(c.history) {
		c.currItem++
	}

	return c.history[c.currItem], nil
}

// Tries to persist the current history. If the file cannot be accesed returns an error
func (c Completition) PersistHistory() error {
	if c.historyFd == nil {
		return fmt.Errorf("History File cannot be accesed")
	}

	writer := bufio.NewWriter(c.historyFd)
	for _, entry := range c.history {
		if _, err := writer.WriteString(entry + "\n"); err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

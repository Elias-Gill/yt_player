package completition

import (
	"bufio"
	"os"
)

const historySize = 20

var (
	history []string
	cache   string

    // magic numbers, dont touch it
	final    int = -1
	position int = 0 
)

// Loads the history file into memory. If the history file is not available, the history
// is still managed in volatile memory.
func LoadHistory() {
	// load cache file (create if necesary)
	path, err := os.UserCacheDir()
	if err != nil {
		return
	}

	cache = path + "/yt_player"

	// if file does not exists, do nothing
	file, err := os.Open(cache)
	if os.IsNotExist(err) {
		return
	}

	// load history on memory (the "i" counter is to
	// prevent the file overpass the history size)
	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan() && i < historySize; i++ {
		history = append(history, scanner.Text())
		final++
		position++
	}
}

func AddHistoryEntry(entry string) {
	final++
	// LRU history replacement
	if final == historySize {
		history = history[1:final]
		final--
	}

	history = append(history, entry)
	position = final
}

// Retrieves the previous history entry iteratively
func NextEntry(entry string) string {
	if len(history) == 0 {
		return entry
	}

	if position > 0 {
		position--
        entry = history[position]
	}

	return entry
}

// Return back one position on history entry
func PrevEntry(entry string) string {
	if len(history) == position || len(history) == 0 {
		return entry
	}

	position++
	if position > final {
		position = final
	}

	return history[position]
}

func PersistHistory() {
	file, err := os.Create(cache)
	if err != nil {
		return
	}

	writer := bufio.NewWriter(file)

	for _, entry := range history {
		_, _ = writer.WriteString(entry + "\n")
	}

	writer.Flush()
}

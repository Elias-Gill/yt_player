package settings

import (
	"flag"
	"log"
	"os"
)

type Settings struct {
	// flags
	apiKey       string
	reattach     bool
	detachOnQuit bool
}

func MustParseConfig() *Settings {
	keyFlag := flag.String("key", "", "Youtube developer key")
	reattach := flag.Bool("reattach", false, "Try to reattach to a previous mpv instance")
	detachOnQuit := flag.Bool("detach-on-quit", false, "Detaches the MPV instance instead of stoping the player.")

	flag.Parse()

	var apiKey string
	if *keyFlag != "" {
		apiKey = *keyFlag
	} else {
		var exists bool
		apiKey, exists = os.LookupEnv("YT_PLAYER_KEY")

		if !exists {
			log.Fatal("Cannot retrieve youtube API key." +
				"\n\nPlease submit the API key using the '--key=<key>' flag" +
				"\nor setting the 'YT_PLAYER_KEY' env variable")
		}
	}

	return &Settings{
		apiKey:       apiKey,
		reattach:     *reattach,
		detachOnQuit: *detachOnQuit,
	}
}

func (s Settings) GetApiKey() string {
	return s.apiKey
}

func (s Settings) Tryreattach() bool {
	return s.reattach
}

func (s Settings) DetachOnQuit() bool {
	return s.detachOnQuit
}

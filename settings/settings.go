package settings

import (
	"flag"
	"log"
	"os"
)

type Settings struct {
	// flags
	apiKey string
}

func MustParseConfig() *Settings {
	keyFlag := *flag.String("key", "", "Youtube developer key")
	flag.Parse()

	var apiKey string
	if keyFlag != "" {
		apiKey = keyFlag
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
		apiKey: apiKey,
	}
}

func (s Settings) GetApiKey() string {
	return s.apiKey
}

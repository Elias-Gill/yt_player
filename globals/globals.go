package globals

import (
	"flag"
	"fmt"
	"os"
)

var (
    // flags 
	maxResults = flag.Int64("max-results", 50, "Max YouTube results")
	keyFlag    = flag.String("key", "", "Youtube developer key")

	apiKey        string
)

func GetMaxResults() int64 {
    return *maxResults
}

func GetApiKey() string {
    return apiKey
}

func ParseFlags() {
	flag.Parse()

	if *keyFlag != "" {
		apiKey = *keyFlag
	} else {
		var exists bool
		apiKey, exists = os.LookupEnv("YT_PLAYER_KEY")

		if !exists {
			fmt.Println("Cannot retrieve api key. \n\nPlease submit the key using the '--key=<key>' flag \nor setting the 'YT_PLAYER_KEY' env variable")
			os.Exit(1)
		}
	}

}

# Youtube music player
This is a simple wrapper for the youtube API that search and list youtube videos.


### Dependencies
For playing music this program uses [MPV player](https://mpv.io/).

# Building

Clone the repo and run 
```bash
go install .
```

# Usage
You have to go to [Google console](https://console.cloud.google.com/). Then create a new project and
enable the for `Youtube API v3` for that project. 

Finally go to the [API pannel](https://console.cloud.google.com/apis/api/youtube.googleapis.com/credentials) and copy the api key.

You can provide the api key using the `--key=<key>` flag or setting the `YT_PLAYER_KEY` enviroment variable

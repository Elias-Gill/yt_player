# YouTube Music Player 

Welcome to the YouTube Music Player, a simple and usefull wrapper for the YouTube API and the
MPV player.
This tool allows you to effortlessly search for and list YouTube videos, and play them directly
using MPV.

**Important Note**:
To use this application, you will need a valid YouTube API key.

# Dependencies 

This program relies on the following dependencies:
- MPV Player:
  A versatile media player that supports a wide range of formats.
- youtube-dl:
  A command-line program to download videos from YouTube and other sites.

On `Windows` you can easily install both MPV and youtube-dl using Chocolatey:
```bash
choco install mpv youtube-dl
```

On `Linux` and `Mac`, all major distributions provide the necessary dependencies for this
project, which can be easily installed using your package manager.

# Installation

As we (currently) not distribute pre-compiled executables, you can install the project using
Go's package manager:

```bash
go install github.com/elias-gill/yt_player@latest 
```

Alternativelly you can build this project from source with:

```bash
go install .
```

# Usage Instructions

1. **Obtain a YouTube API Key**:
   - Navigate to the [Google Cloud Console](https://console.cloud.google.com/).
   - Create a new project and enable the **YouTube Data API v3** for that project.
   - Go to the
     [API Credentials page](https://console.cloud.google.com/apis/api/youtube.googleapis.com/credentials)
     and copy your API key.

2. **Provide the API Key**:
   You can provide your API key in one of two ways:
   - Using the command line flag:
     `--key=<your_api_key>`
   - Setting the environment variable:
     `YT_PLAYER_KEY=<your_api_key>`

# Features

- **Search and Play**:
  Quickly search for YouTube videos and play them using MPV.
- **User-Friendly**:
  Simple command-line interface for easy navigation and playback.

# Contributing

We welcome contributions!
If you have suggestions for improvements or new features, feel free to open an issue or submit
a pull request.

# License 

This project is licensed under the MIT License.
See the [LICENSE](LICENSE.md) file for more details.

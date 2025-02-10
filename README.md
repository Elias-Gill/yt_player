# YT Player

A simple wrapper for the YouTube API and MPV player.
This tool allows you to effortlessly search for and list YouTube videos, and play them directly
using MPV.

**NOTE**: Windows NOT SUPPORTED

![Demo](https://github.com/user-attachments/assets/f18fdf9e-2a61-4277-83c4-8542069eb923)

## **Important Note**

To use this application, you will need a valid **YouTube API key**.
Follow the instructions below to obtain one.

### 1. Obtain a YouTube API Key

1. Navigate to the [Google Cloud Console](https://console.cloud.google.com/).
2. Create a new project and enable the **YouTube Data API v3** for that project.
3. Go to the
   [API Credentials page](https://console.cloud.google.com/apis/api/youtube.googleapis.com/credentials)
   and generate an API key.

### 2. Provide the API Key

You can provide your API key in one of two ways:
- **With a Command Line Flag**:
```bash
yt_player --key=<your_api_key>
```
- **As an Environment Variable**:
```bash
# Linux and macOS
export YT_PLAYER_KEY=<your_api_key>
```


## **Dependencies**

This program relies on the following dependencies:

- **MPV Player**:
  A versatile media player that supports a wide range of formats.
- **yt-dlp**:
  A command-line program to download videos from YouTube and other sites.

### **Installation Instructions**

#### **Linux and macOS** 

Most distributions provide the necessary dependencies, which can be installed using your
package manager.

- **Debian/Ubuntu**:
  ```bash
  sudo apt install mpv yt-dlp
  ```
- **macOS (Homebrew)**:
  ```bash
  brew install mpv yt-dlp
  ```

---

## **Installation**

### **Pre-compiled Binaries** 

We now provide pre-compiled binaries for Windows, macOS, and Linux on the
[Releases page](https://github.com/elias-gill/yt_player/releases).
Simply download the appropriate binary for your platform and run it.

### **Building from Source** 

If you prefer to build the project from source, you can install it using Go's package manager:

```bash
go install github.com/elias-gill/yt_player@latest
```

Alternatively, you can clone the repository and build it locally:

```bash
git clone https://github.com/elias-gill/yt_player.git
cd yt_player
go install .
```

---

## **Features**

- **Search and Play**:
  Quickly search for YouTube videos and play them using MPV.
- **User-Friendly Interface**:
  Simple command-line interface for easy navigation and playback.
- **Detach mode**:
  You can detach and reattach to the player.
  Very usefull when you want to quickly open and close tmux floating windows or pannels to
  control your music.

To see available command line flags use the "--help" flag.

## **Future Plans**

- Support for listing and playing YouTube playlists.
- A playback queue for continuous video playback.
- Custom playlists for saving and managing favorite videos.

## **Contributing**

We welcome contributions!
If you have suggestions for improvements or new features, feel free to:
- Open an issue on [GitHub](https://github.com/elias-gill/yt_player/issues).
- Submit a pull request with your changes.

## **License**

This project is licensed under the **MIT License**.
See the [LICENSE](LICENSE) file for more details.

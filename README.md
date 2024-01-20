# YouTube To Spotify Playlist Converter

Convert YouTube playlists to Spotify playlists via different programming languages


## Prerequisites

* <span style="color:#b2071d">YouTube API Key</span> - Follow the steps in the official
  documentation [YouTube API][YouTubeAPILink]

---

* <span style="color:#b2071d">YouTube Playlist ID</span> - Go to a YouTube playlist, and the ID will appear in the
  url (.../playlist?list=`ID`)

---

* <span style="color:#1db954">Spotify Client ID & Secret</span> - Follow the steps in the
  official documentation [Spotify API][SpotifyAPILink]

---

* <span style="color:#1db954">Spotify Access Token</span> (**Lasts 1 Hour**) - Choose one of the methods below:
    * Follow the steps in the official documentation [Spotify Access Token][SpotifyTokenLink]
    * First, make sure the token configuration file [config/token_config.json](config/token_config.json) is ready. Then, 
      run one of the following scripts:
        1. [receive_token.py](python/src/spotify/receive_token.py)
        2. [receiveToken.ts](typescript/spotify/receiveToken.ts)

---

* <span style="color:#1db954">Spotify Playlist ID</span> - Go to a Spotify playlist, and the ID will
  appear in the url (.../playlist/`ID`)


## Config Files

* [config/api_config.json](config/api_config.json): This file contains information about the relevant APIs. It should
  not be touched.


* [config/user_config.json](config/user_config.json): This file contains the user's data & preferences:
  ```shell
  {
    "logging": {
      # OPTIONAL: Change the logging level
      "level": "INFO"
    },
    "youtube": {
      # Set the YouTube API Key as described in the 'Prerequisites' section
      "api_key": null, 
  
      # Set the YouTube Playlist IDs (source) as described in the 'Prerequisites' section
      "playlist_ids": [] 
    },
    "spotify": {
      # Set the Spotify Access Token as described in the 'Prerequisites' section
      "access_token": null, 
      
      # NOTE: You can choose to add tracks to new and/or existing playlists
  
      # Provide the information for the new Spotify playlists (destination)
      # It should be structured as follows: [{"name": "<name>", "description": "<description>", "public": <true|false>}, ...]
      "new_playlists": [], 
        
      # Provide the Spotify Playlist IDs (destination) as described in the 'Prerequisites' section
      "existing_playlist_ids": [] 
    }
  }
  ```
  **__NOTE:__** It should be filled before running one of the main programs.


* [config/token_config.json](config/token_config.json): This file contains token related data:
  ```shell
  {
    "spotify": {
      # Set the Spotify Client ID as described in the 'Prerequisites' section
      "client_id": null,
  
      # Set the Spotify Client Secret as described in the 'Prerequisites' section
      "client_secret": null,
  
      # Set the Spotify redirect URI (It is included in the Client ID/Secret creation process) 
      "redirect_uri": null
    }
  }
  ```
  **__NOTE:__** It should be filled before running one of the token receival scripts.


## Samples

The samples directory [samples](samples) contains examples of how the config files should look.

* For user configuration, look at [samples/user_config.json](samples/user_config.json)

* For token configuration, look at [samples/token_config.json](samples/token_config.json)


## Running Code

First, make sure the user configuration file [config/user_config.json](config/user_config.json) is ready. Then, 
run one of the following programs:

1. **Python 3^6**: [main.py](python/main.py)
    ```shell
    python python/main.py
    ```
2. **TypeScript**: [main.ts](typescript/main.ts)
    ```shell
    ts-node typescript/main.ts
    ```
   

[YouTubeAPILink]:https://developers.google.com/youtube/v3/getting-started
[SpotifyAPILink]:https://developer.spotify.com/documentation/web-api/concepts/apps
[SpotifyTokenLink]:https://developer.spotify.com/documentation/web-api/concepts/access-token
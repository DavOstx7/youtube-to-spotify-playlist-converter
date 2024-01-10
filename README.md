# YouTubeToSpotifyPlaylistConverter

Convert YouTube playlists to Spotify via different programming languages


## Prerequisites

* <span style="color:#b2071d">YouTube API Key</span> - Follow the steps in the
  documentation [YouTube API][YouTubeAPILink]



* <span style="color:#b2071d">YouTube Playlist ID</span> - Go to your YouTube playlist, and the ID will appear in the
  url (.../playlist?list=`ID`)



* <span style="color:#1db954">Spotify Client ID & Secret</span> - Follow the steps in the
  documentation [Spotify API][SpotifyAPILink]



* <span style="color:#1db954">Spotify Access Token</span> (**Lasts 1 Hour**) - Choose one of the methods below:
    * Follow the steps in the documentation [Spotify Access Token][SpotifyTokenLink]
    * Fill the values inside [config/token_config.json](config/token_config.json). Then run one of the following
      scripts:
        1. [receive_token.py](python/src/spotify/receive_token.py)
        2. [receiveToken.ts](typescript/spotify/receiveToken.ts)



* <span style="color:#1db954">Spotify Playlist ID</span> (**Optional**) - Go to your Spotify playlist, and the ID will
  appear in the url (.../playlist/`ID`)


## Config Files

* [config/api_config.json](config/api_config.json): This file contains relevant information about the APIs. It should
  remain untouched.


* [config/user_config.json](config/user_config.json): This file contains the relevant access information/preferences:
  ```shell
  {
    "logging": {
      # Set the logging level of the main program
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
      
      # OPTIONAL: Provide the information for the new Spotify playlists (destination). It should be structured as follows: [{"name": "<name>", "description": "<description>", "public": <true|false>}, ...]
      "new_playlists": [], 
        
      # OPTIONAL: Provide the Spotify Playlist IDs (destination) as described in the 'Prerequisites' section
      "existing_playlist_ids": [] 
    }
  }
  ```
  It should be filled by the user before running one of the main programs.


* [config/token_config.json](config/token_config.json): This file contains the relevant information for receiving an
  access token:
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
  It should be filled by the user before running one of the 'receive-token' scripts.


## Samples

You can look at the samples directory [samples](samples) to get an idea for how you should fill the needed config files:

* For user configuration, look at [samples/user_config.json](samples/user_config.json)

* For token configuration, look at [samples/token_config.json](samples/token_config.json)

[YouTubeAPILink]:https://developers.google.com/youtube/v3/getting-started

[SpotifyAPILink]:https://developer.spotify.com/documentation/web-api/concepts/apps

[SpotifyTokenLink]:https://developer.spotify.com/documentation/web-api/concepts/access-token


## Running Code

Fill the values inside [config/user_config.json](config/user_config.json). Then run one of the following main programs:

1. **Python 3^6**: [main.py](python/main.py)
    ```shell
    python python/main.py
    ```
2. **TypeScript**: [main.ts](typescript/main.ts)
    ```shell
    ts-node typescript/main.ts
    ```
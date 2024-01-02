# YouTubeToSpotifyPlaylistConverter

Convert YouTube playlists to Spotify via different programming languages

## Prerequisites

* <span style="color:#b2071d">YouTube API Key</span> - Follow the steps in the
  documentation [YouTube API][YouTubeAPILink]

---

* <span style="color:#b2071d">YouTube Playlist ID</span> - Go to your YouTube playlist, and the ID will appear in the
  url (.../playlist?list=`ID`)

---

* <span style="color:#1db954">Spotify Client ID & Secret</span> - Follow the steps in the
  documentation [Spotify API][SpotifyAPILink]

---

* <span style="color:#1db954">Spotify Access Token</span> (**Lasts 1 Hour**) - Choose one of the methods below:
  * Follow the steps in the documentation [Spotify Access Token][SpotifyTokenLink]
  * Fill the values inside [config/token_config.json](config/token_config.json). Then run one of the following scripts:
     * [receive_token.py](python/src/spotify/receive_token.py)
     * [receiveToken.ts](typescript/spotify/receiveToken.ts)

---

* <span style="color:#1db954">Spotify Playlist ID</span> (**Optional**) - Go to your Spotify playlist, and the ID will
  appear in the url (.../playlist/`ID`)

## Running Code

Fill the values inside [config/user_config.json](config/user_config.json). Then run one of the following scripts:

* **Python 3^6** - [main.py](python/main.py)
* **TypeScript** - [main.ts](typescript/main.ts)

## Samples

You can look at the samples directory [samples](samples) to get an idea for how you should fill the config files:

* For token configuration, look at [samples/token_config.json](samples/token_config.json)
* For user configuration, look at [samples/user_config.json](samples/user_config.json)

[YouTubeAPILink]:https://developers.google.com/youtube/v3/getting-started

[SpotifyAPILink]:https://developer.spotify.com/documentation/web-api/concepts/apps

[SpotifyTokenLink]:https://developer.spotify.com/documentation/web-api/concepts/access-token 
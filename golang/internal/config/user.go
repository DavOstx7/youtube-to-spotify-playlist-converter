package config

type SpotifyPlaylistInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type UserConfig struct {
	Logging struct {
		Level string `json:"level"`
	} `json:"logging"`

	YouTube struct {
		APIKey      string   `json:"api_key"`
		PlaylistIDs []string `json:"playlist_ids"`
	} `json:"youtube"`

	Spotify struct {
		AccessToken         string                `json:"access_token"`
		NewPlaylists        []SpotifyPlaylistInfo `json:"new_playlists"`
		ExistingPlaylistIDs []string              `json:"existing_playlist_ids"`
	} `json:"spotify"`
}

var userConfig *UserConfig

func GetUserConfig() *UserConfig {
	if userConfig != nil {
		return userConfig
	}

	userConfig = new(UserConfig)
	if err := loadConfig("user_config.json", userConfig); err != nil {
		panic(err)
	}

	return userConfig
}

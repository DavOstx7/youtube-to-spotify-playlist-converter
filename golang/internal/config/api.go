package config

type APIConfig struct {
	YouTube struct {
		URLs struct {
			PlaylistItems string `json:"playlist_items"`
		} `json:"urls"`

		MaxItemsPerRequest int `json:"max_items_per_request"`
	} `json:"youtube"`

	Spotify struct {
		URLs struct {
			UserProfile   string `json:"user_profile"`
			Search        string `json:"search"`
			Token         string `json:"token"`
			Authorization string `json:"authorization"`
			Playlists     string `json:"playlists"`
			Tracks        string `json:"tracks"`
		} `json:"urls"`

		AuthorizationScopes string `json:"authorization_scopes"`
		MaxTracksPerRequest int    `json:"max_tracks_per_request"`
	} `json:"spotify"`
}

var apiConfig *APIConfig

func GetAPIConfig() *APIConfig {
	if apiConfig != nil {
		return apiConfig
	}

	apiConfig = new(APIConfig)
	if err := loadConfig("api_config.json", apiConfig); err != nil {
		panic(err)
	}

	return apiConfig
}

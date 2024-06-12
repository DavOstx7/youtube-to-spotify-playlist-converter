package config

type SpotifyTokenConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type TokenConfig struct {
	Spotify SpotifyTokenConfig `json:"spotify"`
}

var tokenConfig *TokenConfig

func GetTokenConfig() *TokenConfig {
	if tokenConfig != nil {
		return tokenConfig
	}

	tokenConfig = new(TokenConfig)
	if err := loadConfig("token_config.json", tokenConfig); err != nil {
		panic(err)
	}

	return tokenConfig
}

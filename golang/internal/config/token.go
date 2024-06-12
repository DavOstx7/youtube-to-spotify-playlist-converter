package config

type SpotifyTokenConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type TokenConfig struct {
	Spotify SpotifyTokenConfig `json:"spotify"`
}

var tokenConfigRef *TokenConfig

func GetTokenConfig() *TokenConfig {
	if tokenConfigRef != nil {
		return tokenConfigRef
	}

	tokenConfigRef = &TokenConfig{}
	if err := loadConfig("token_config.json", tokenConfigRef); err != nil {
		panic(err)
	}

	return tokenConfigRef
}

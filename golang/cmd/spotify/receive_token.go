package main

import (
	"github.com/gin-gonic/gin"
	"golang/internal/config"
	"golang/internal/requests"
	"golang/internal/spotify"
	"net/http"
	"net/url"
)

var TokenConfig = &config.GetTokenConfig().Spotify

func authorize(c *gin.Context) {
	queryParams := spotify.NewAuthorizationQueryParams(TokenConfig.ClientID, TokenConfig.RedirectURI)
	url, err := requests.BuildURLWithQuery(spotify.APIConfig.URLs.Authorization, queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, url)
}

func accessTokenCallback(c *gin.Context) {
	response := spotify.GetAccessToken(TokenConfig, c.Query("code"))
	c.JSON(http.StatusTemporaryRedirect, gin.H{"access_token": response.AccessToken})
}

func main() {
	parsedRedirectURI, err := url.Parse(TokenConfig.RedirectURI)
	if err != nil {
		panic(err)
	}

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/", authorize)
	router.GET(parsedRedirectURI.Path, accessTokenCallback)
	router.Run(parsedRedirectURI.Host)
}

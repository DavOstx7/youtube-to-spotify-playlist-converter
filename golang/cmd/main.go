package main

import (
	"golang/internal"
	"golang/internal/config"
)

func main() {
	converter := internal.NewPlaylistsConverter(config.GetUserConfig())
	converter.Setup()
	converter.Run()
}

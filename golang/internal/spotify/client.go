package spotify

import (
	"fmt"
	"golang/internal/config"
	"log/slog"
	"sync"
)

const SingleTrack = 1

type SpotifyClient struct {
	accessToken string
	userID      *string
}

func NewSpotifyClient(accessToken string) *SpotifyClient {
	return &SpotifyClient{accessToken: accessToken}
}

func (sc *SpotifyClient) SetUserID() {
	userProfile := GetUserProfile(sc.accessToken)
	sc.userID = &userProfile.ID
}

func (sc *SpotifyClient) CreatePlaylists(playlistsInfo []config.SpotifyPlaylistInfo) <-chan string {
	playlistIDChan := make(chan string, len(playlistsInfo))
	var wg sync.WaitGroup
	wg.Add(len(playlistsInfo))

	for _, playlistInfo := range playlistsInfo {
		playlistInfo := playlistInfo
		go func() {
			defer wg.Done()
			sc.createPlaylist(playlistIDChan, playlistInfo)
		}()
	}
	go func() {
		wg.Wait()
		close(playlistIDChan)
	}()

	return playlistIDChan
}

func (sc *SpotifyClient) SearchForTrackURIs(trackNames []string) <-chan string {
	trackURIChan := make(chan string, len(trackNames))
	var wg sync.WaitGroup
	wg.Add(len(trackNames))
	
	slog.Info(fmt.Sprintf("Starting to search Spotify track uris for %d track names...", len(trackNames)))
	for _, trackName := range trackNames {
		trackName := trackName
		go func() {
			defer wg.Done()
			sc.searchForTrackURI(trackURIChan, trackName)
		}()
	}
	go func() {
		wg.Wait()
		close(trackURIChan)
	}()

	return trackURIChan
}

func (sc *SpotifyClient) AddTracksToPlaylist(snapshotIDChan chan<- string, playlistID string, trackURIs []string, position int) {
	slog.Info(fmt.Sprintf("Adding %d track uris to Spotify playlist '%s'", len(trackURIs), playlistID))
	snapshotID := AddTracksToPlaylist(sc.accessToken, playlistID, trackURIs, position)
	slog.Debug(fmt.Sprintf("Snapshot id of the new added tracks for Spotify playlist '%s' is '%s'", playlistID, snapshotID))
	snapshotIDChan <- snapshotID
}

func (sc *SpotifyClient) createPlaylist(playlistIDChan chan<- string, playlistInfo config.SpotifyPlaylistInfo) {
	if sc.userID == nil {
		panic("user ID must be set in order to create Spotify playlist")
	}

	if playlistInfo.Public {
		slog.Info(fmt.Sprintf("Creating new public Spotify playlist '%s'", playlistInfo.Name))
	} else {
		slog.Info(fmt.Sprintf("Creating new private Spotify playlist '%s'", playlistInfo.Name))
	}
	response := CreatePlaylist(sc.accessToken, *sc.userID, playlistInfo)
	playlistIDChan <- response.ID
}

func (sc *SpotifyClient) searchForTrackURI(trackURIChan chan<- string, trackName string) {
	tracks := SearchForTracks(sc.accessToken, trackName, SingleTrack)
	if tracks == nil || len(tracks.Items) == 0 {
		slog.Warn(fmt.Sprintf("Failed to find Spotify track uri for '%s'", trackName))
		return
	}
	
	slog.Debug(fmt.Sprintf("Found Spotify track uri '%s' for '%s'", tracks.Items[0].URI, trackName))
	trackURIChan <- tracks.Items[0].URI
}

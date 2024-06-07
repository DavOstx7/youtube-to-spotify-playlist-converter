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
	var wg sync.WaitGroup
	playlistIDChan := make(chan string, len(playlistsInfo))
	for _, playlistInfo := range playlistsInfo {
		playlistInfo := playlistInfo
		wg.Add(1)
		go func() {
			defer wg.Done()
			if playlistID, err := sc.createPlaylist(playlistInfo); err == nil {
				playlistIDChan <- *playlistID
			}
		}()
	}
	go func() {
		wg.Wait()
		close(playlistIDChan)
	}()
	return playlistIDChan
}

func (sc *SpotifyClient) SearchForTrackURIs(trackNames []string) <-chan string {
	slog.Info(fmt.Sprintf("Starting to search Spotify track uris for %d track names...", len(trackNames)))

	var wg sync.WaitGroup
	trackURIChan := make(chan string, len(trackNames))
	for _, trackName := range trackNames {
		trackName := trackName
		wg.Add(1)
		go func() {
			defer wg.Done()
			trackURI := sc.searchForTrackURI(trackName)
			if trackURI != nil {
				trackURIChan <- *trackURI
			}
		}()
	}
	go func() {
		wg.Wait()
		close(trackURIChan)
	}()
	return trackURIChan
}

func (sc *SpotifyClient) AddTracksToPlaylist(playlistID string, trackURIs []string, position int) string {
	slog.Info(fmt.Sprintf("Adding %d track uris to Spotify playlist '%s'", len(trackURIs), playlistID))
	snapshotID := AddTracksToPlaylist(sc.accessToken, playlistID, trackURIs, position)
	slog.Debug(fmt.Sprintf("Snapshot id of the new added tracks for Spotify playlist '%s' is '%s'", playlistID, snapshotID))
	return snapshotID
}

func (sc *SpotifyClient) createPlaylist(playlistInfo config.SpotifyPlaylistInfo) (*string, error) {
	if sc.userID == nil {
		return nil, fmt.Errorf("user ID must be set in order to create Spotify playlist")
	}

	if playlistInfo.Public {
		slog.Info(fmt.Sprintf("Creating new public Spotify playlist '%s'", playlistInfo.Name))
	} else {
		slog.Info(fmt.Sprintf("Creating new private Spotify playlist '%s'", playlistInfo.Name))
	}
	response := CreatePlaylist(sc.accessToken, *sc.userID, playlistInfo)
	return &response.ID, nil
}

func (sc *SpotifyClient) searchForTrackURI(trackName string) *string {
	tracks := SearchForTracks(sc.accessToken, trackName, SingleTrack)
	if tracks == nil || len(tracks.Items) == 0 {
		slog.Warn(fmt.Sprintf("Failed to find Spotify track uri for '%s'", trackName))
		return nil
	}

	slog.Debug(fmt.Sprintf("Found Spotify track uri '%s' for '%s'", tracks.Items[0].URI, trackName))
	return &tracks.Items[0].URI
}

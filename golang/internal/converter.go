package internal

import (
	"golang/internal/config"
	"golang/internal/spotify"
	"golang/internal/utils"
	"golang/internal/youtube"
	"log/slog"
	"os"
	"sync"
)

const MaxTitlesBatchSize = 50
const TracksInsertionPosition = 0

type PlaylistsConverter struct {
	userConfig         *config.UserConfig
	youtubeClient      *youtube.YouTubeClient
	spotifyClient      *spotify.SpotifyClient
	spotifyPlaylistIDs []string
}

func NewPlaylistsConverter(userConfig *config.UserConfig) *PlaylistsConverter {
	return &PlaylistsConverter{
		userConfig:         userConfig,
		youtubeClient:      youtube.NewYouTubeClient(userConfig.YouTube.APIKey),
		spotifyClient:      spotify.NewSpotifyClient(userConfig.Spotify.AccessToken),
		spotifyPlaylistIDs: nil,
	}
}

func (pc *PlaylistsConverter) Run() [][]string {
	if pc.spotifyPlaylistIDs == nil {
		panic("playlists converter must be set up before running")
	}

	var wg sync.WaitGroup
	snapshotIDsChan := make(chan []string)
	titlesBatchChan := pc.youtubeClient.FetchPlaylistsTitles(pc.userConfig.YouTube.PlaylistIDs, MaxTitlesBatchSize)
	for titlesBatch := range titlesBatchChan {
		titlesBatch := titlesBatch
		wg.Add(1)
		go func() {
			defer wg.Done()
			snapshotIDs := utils.CollectFromChannel(pc.addYouTubeTitlesToSpotifyPlaylists(titlesBatch))
			if snapshotIDs != nil {
				snapshotIDsChan <- snapshotIDs
			}
		}()
	}
	go func() {
		wg.Wait()
		close(snapshotIDsChan)
	}()
	return utils.CollectFromChannel(snapshotIDsChan)
}

func (pc *PlaylistsConverter) Setup() {
	pc.setupLogs()
	pc.setupSpotify()
}

func (pc *PlaylistsConverter) addYouTubeTitlesToSpotifyPlaylists(titles []string) <-chan string {
	trackURIs := utils.CollectFromChannel(pc.spotifyClient.SearchForTrackURIs(titles))
	if len(trackURIs) == 0 {
		slog.Warn("Could not find a single Spotify track uri for the given track names")
		return nil
	}

	var wg sync.WaitGroup
	snapshotIDChan := make(chan string, len(trackURIs))
	for _, spotifyPlaylistID := range pc.spotifyPlaylistIDs {
		spotifyPlaylistID := spotifyPlaylistID
		wg.Add(1)
		go func() {
			defer wg.Done()
			snapshotIDChan <- pc.spotifyClient.AddTracksToPlaylist(spotifyPlaylistID, trackURIs, TracksInsertionPosition)
		}()
	}
	go func() {
		wg.Wait()
		close(snapshotIDChan)
	}()
	return snapshotIDChan
}

func (pc *PlaylistsConverter) setupLogs() {
	slogLevel, err := utils.ConvertLevelToSlogLevel(pc.userConfig.Logging.Level)
	if err != nil {
		panic(err)
	}

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: false,
	})
	slog.SetDefault(slog.New(logHandler))
}

func (pc *PlaylistsConverter) setupSpotify() {
	pc.spotifyClient.SetUserID()
	pc.spotifyPlaylistIDs = append(pc.spotifyPlaylistIDs, pc.userConfig.Spotify.ExistingPlaylistIDs...)
	if len(pc.userConfig.Spotify.NewPlaylists) > 0 {
		newPlaylistIDs := utils.CollectFromChannel(pc.spotifyClient.CreatePlaylists(pc.userConfig.Spotify.NewPlaylists))
		pc.spotifyPlaylistIDs = append(pc.spotifyPlaylistIDs, newPlaylistIDs...)
	}
}

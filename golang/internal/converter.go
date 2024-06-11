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
const AddTitlesToPlaylistsWorkers = 5

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

func (pc *PlaylistsConverter) Setup() {
	pc.setupLogs()
	pc.setupSpotify()
}

func (pc *PlaylistsConverter) Run() <-chan string {
	if pc.spotifyPlaylistIDs == nil {
		panic("playlists converter must be set up before running")
	}

	titleBatchChan := pc.youtubeClient.FetchPlaylistsTitles(pc.userConfig.YouTube.PlaylistIDs, MaxTitlesBatchSize)
	snapshotIDChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(AddTitlesToPlaylistsWorkers)
	
	for i := 0; i < AddTitlesToPlaylistsWorkers; i++ {
		go func() {
			defer wg.Done()
			pc.addTitlesToPlaylists(titleBatchChan, snapshotIDChan)
		}()
	}
	go func() {
		wg.Wait()
		close(snapshotIDChan)
	}()

	return snapshotIDChan
}

func (pc *PlaylistsConverter) addTitlesToPlaylists(titleBatchChan <-chan []string, snapshotIDChan chan<- string) {
	for titleBatch := range titleBatchChan {
		trackURIs := utils.CollectFromChannel(pc.spotifyClient.SearchForTrackURIs(titleBatch))
		if len(trackURIs) == 0 {
			slog.Warn("Could not find a single Spotify track uri for the given track names")
			continue
		}

		for _, playlistID := range pc.spotifyPlaylistIDs {
			playlistID := playlistID
			pc.spotifyClient.AddTracksToPlaylist(snapshotIDChan, playlistID, trackURIs, TracksInsertionPosition)
		}
	}
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

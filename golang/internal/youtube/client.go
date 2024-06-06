package youtube

import (
	"fmt"
	"log/slog"
)

const MaxResults = 50

type YouTubeClient struct {
	apiKey string
}

func NewYouTubeClient(apiKey string) *YouTubeClient {
	return &YouTubeClient{apiKey: apiKey}
}

func (yc *YouTubeClient) FetchPlaylistsTitles(playlistIDs []string, maxBatchSize int) <-chan []string {
	titlesBatchChan := make(chan []string)
	go func() {
		defer close(titlesBatchChan)
		
		var titlesBatch []string
		for _, playlistID := range playlistIDs {
			slog.Info(fmt.Sprintf("Starting to search for YouTube video titles inside playlist '%s'...", playlistID))
						
			playlist := NewYouTubePlaylist(yc.apiKey, playlistID, MaxResults)
			for page := range playlist.FetchPlaylistItemsPages() {
				for _, item := range page.Items {
					slog.Debug(fmt.Sprintf("Found YouTube video title '%s'", item.Snippet.Title))
					titlesBatch = append(titlesBatch, item.Snippet.Title)

					if len(titlesBatch) >= maxBatchSize {
						titlesBatchChan <- titlesBatch
						titlesBatch = nil
					}
				}
			}
		}
		if len(titlesBatch) > 0 {
			titlesBatchChan <- titlesBatch
		}
	}()
	return titlesBatchChan
}

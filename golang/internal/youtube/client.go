package youtube

import (
	"fmt"
	"log/slog"
	"sync"
)

const MaxResults = 50

type YouTubeClient struct {
	apiKey string
}

func NewYouTubeClient(apiKey string) *YouTubeClient {
	return &YouTubeClient{apiKey: apiKey}
}

func (yc *YouTubeClient) FetchPlaylistsTitles(playlistIDs []string, maxBatchSize int) <-chan []string {
	itemsChan := yc.FetchPlaylistsItems(playlistIDs)
	titleChan := yc.ExtractPlaylistsTitles(itemsChan)
	titleBatchChan := yc.SplitPlaylistsTitles(titleChan, maxBatchSize)
	return titleBatchChan
}

func (yc *YouTubeClient) FetchPlaylistsItems(playlistIDs []string) <-chan []PlaylistItem {
	itemsChan := make(chan []PlaylistItem)
	var wg sync.WaitGroup
	wg.Add(len(playlistIDs))

	for _, playlistID := range playlistIDs {
		playlistID := playlistID
		go func() {
			defer wg.Done()

			playlist := NewYouTubePlaylist(yc.apiKey, playlistID, MaxResults)
			slog.Info(fmt.Sprintf("Starting to search for YouTube video titles inside playlist '%s'...", playlistID))
			playlist.FetchAllPlaylistItems(itemsChan)
		}()
	}
	go func() {
		wg.Wait()
		close(itemsChan)
	}()

	return itemsChan
}

func (yc *YouTubeClient) ExtractPlaylistsTitles(itemsChan <-chan []PlaylistItem) <-chan string {
	titleChan := make(chan string)

	go func() {
		defer close(titleChan)

		for items := range itemsChan {
			for _, item := range items {
				slog.Debug(fmt.Sprintf("Found YouTube video title '%s'", item.Snippet.Title))
				titleChan <- item.Snippet.Title
			}
		}
	}()

	return titleChan
}

func (yc *YouTubeClient) SplitPlaylistsTitles(titleChan <-chan string, maxBatchSize int) <-chan []string {
	titleBatchChan := make(chan []string)

	go func() {
		defer close(titleBatchChan)

		var titleBatch []string
		for title := range titleChan {
			titleBatch = append(titleBatch, title)
			if len(titleBatch) >= maxBatchSize {
				titleBatchChan <- titleBatch
				titleBatch = nil
			}
		}
		if len(titleBatch) > 0 {
			titleBatchChan <- titleBatch
		}
	}()

	return titleBatchChan
}

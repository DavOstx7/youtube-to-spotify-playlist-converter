package youtube

import (
	"fmt"
	"log/slog"
)

type YouTubePlaylist struct {
	queryParams  *PlaylistItemsQueryParams
	lastResponse *GetPlaylistItemsResponse
}

func NewYouTubePlaylist(apiKey string, playlistID string, maxResults int) *YouTubePlaylist {
	return &YouTubePlaylist{
		queryParams:  NewPlaylistItemsQueryParams(apiKey, playlistID, maxResults),
		lastResponse: nil,
	}
}

func (yp *YouTubePlaylist) HasResponse() bool {
	return yp.lastResponse != nil
}

func (yp *YouTubePlaylist) Reset() {
	yp.queryParams.PageToken = nil
	yp.lastResponse = nil
}

func (yp *YouTubePlaylist) HasPrevPage() bool {
	if !yp.HasResponse() {
		return false
	}
	return yp.lastResponse.PrevPageToken != nil
}

func (yp *YouTubePlaylist) HasNextPage() bool {
	if !yp.HasResponse() {
		return false
	}
	return yp.lastResponse.NextPageToken != nil
}

func (yp *YouTubePlaylist) SetPrevPage() bool {
	if !yp.HasPrevPage() {
		return false
	}
	yp.queryParams.PageToken = yp.lastResponse.PrevPageToken
	return true
}

func (yp *YouTubePlaylist) SetNextPage() bool {
	if !yp.HasNextPage() {
		return false
	}
	yp.queryParams.PageToken = yp.lastResponse.NextPageToken
	return true
}

func (yp *YouTubePlaylist) FetchPlaylistItemsFromPage(playlistItemsChan chan<- []PlaylistItem) {
	if !yp.HasResponse() {
		slog.Debug(fmt.Sprintf("Searching initial YouTube playlist page for playlist '%s'", yp.queryParams.PlaylistID))
	} else {
		slog.Debug(
			fmt.Sprintf("Searching YouTube playlist page for playlist '%s' with page token '%s'",
				yp.queryParams.PlaylistID, *yp.queryParams.PageToken),
		)
	}

	response := GetPlaylistItems(yp.queryParams)
	playlistItemsChan <- response.Items
	yp.lastResponse = response
}

func (yp *YouTubePlaylist) FetchAllPlaylistItems(playlistItemsChan chan<- []PlaylistItem) {
	yp.Reset()

	yp.FetchPlaylistItemsFromPage(playlistItemsChan)
	for yp.SetNextPage() {
		yp.FetchPlaylistItemsFromPage(playlistItemsChan)
	}
}

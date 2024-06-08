package youtube

import (
	"fmt"
	"log/slog"
)

type YouTubePlaylist struct {
	queryParams     *PlaylistItemsQueryParams
	currentResponse *GetPlaylistItemsResponse
}

func NewYouTubePlaylist(apiKey string, playlistID string, maxResults int) *YouTubePlaylist {
	return &YouTubePlaylist{
		queryParams:     NewPlaylistItemsQueryParams(apiKey, playlistID, maxResults),
		currentResponse: nil,
	}
}

func (yp *YouTubePlaylist) IsInitialState() bool {
	return yp.currentResponse == nil
}

func (yp *YouTubePlaylist) HasPrevPage() bool {
	if yp.IsInitialState() {
		return false
	}
	return yp.currentResponse.PrevPageToken != nil
}

func (yp *YouTubePlaylist) HasNextPage() bool {
	if yp.IsInitialState() {
		return false
	}
	return yp.currentResponse.NextPageToken != nil
}

func (yp *YouTubePlaylist) SetPrevPage() {
	yp.queryParams.PageToken = yp.currentResponse.PrevPageToken
}

func (yp *YouTubePlaylist) SetNextPage() {
	yp.queryParams.PageToken = yp.currentResponse.NextPageToken
}

func (yp *YouTubePlaylist) FetchPlaylistItemsPage() *GetPlaylistItemsResponse {
	if yp.IsInitialState() {
		slog.Debug(fmt.Sprintf("Searching initial YouTube playlist page for playlist '%s'", yp.queryParams.PlaylistID))
	} else {
		slog.Debug(
			fmt.Sprintf("Searching YouTube playlist page for playlist '%s' with page token '%s'",
				yp.queryParams.PlaylistID, *yp.queryParams.PageToken),
		)
	}

	yp.currentResponse = GetPlaylistItems(yp.queryParams)
	return yp.currentResponse
}

func (yp *YouTubePlaylist) FetchPlaylistItemsPages() <-chan *GetPlaylistItemsResponse {
	pageChan := make(chan *GetPlaylistItemsResponse)
	go func() {
		defer close(pageChan)

		pageChan <- yp.FetchPlaylistItemsPage()
		for yp.HasNextPage() {
			yp.SetNextPage()
			pageChan <- yp.FetchPlaylistItemsPage()
		}
	}()
	return pageChan
}

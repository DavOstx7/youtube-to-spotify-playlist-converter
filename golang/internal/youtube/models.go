package youtube

import "golang/internal/common/types"

type PlaylistItemsQueryParams struct {
	Key        string  `url:"key"`
	Part       string  `url:"part"`
	PlaylistID string  `url:"playlistId"`
	MaxResults int     `url:"maxResults"`
	PageToken  *string `url:"pageToken,omitempty"`
}

type PlaylistItemSnippet struct {
	PublishedAt            string     `json:"publishedAt"`
	ChannelID              string     `json:"channelId"`
	Title                  string     `json:"title"`
	Description            string     `json:"description"`
	Thumbnails             types.JSON `json:"thumbnails"`
	ChannelTitle           string     `json:"channelTitle"`
	VideoOwnerChannelTitle string     `json:"videoOwnerChannelTitle"`
	VideoOwnerChannelID    string     `json:"videoOwnerChannelId"`
	PlaylistID             string     `json:"playlistId"`
	Position               int        `json:"position"`
	ResourceID             types.JSON `json:"resourceId"`
	ContentDetails         types.JSON `json:"contentDetails"`
	Status                 types.JSON `json:"status"`
}

type PlaylistItem struct {
	ID      string              `json:"id"`
	Snippet PlaylistItemSnippet `json:"snippet"`
}

type GetPlaylistItemsResponse struct {
	ID            string         `json:"id"`
	NextPageToken *string        `json:"nextPageToken,omitempty"`
	PrevPageToken *string        `json:"prevPageToken,omitempty"`
	PageInfo      types.JSON     `json:"pageInfo"`
	Items         []PlaylistItem `json:"items"`
}

package spotify

type ItemType string

const (
	TrackItem ItemType = "track"
	// More items...
)

type AuthorizationQueryParams struct {
	ClientID     string `url:"client_id"`
	ResponseType string `url:"response_type"`
	RedirectURI  string `url:"redirect_uri"`
	Scope        string `url:"scope"`
}

type GetAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	// More fields...
}

type GetUserProfileResponse struct {
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Href        string `json:"href"`
	URI         string `json:"uri"`
	// More fields...
}

type CreatePlaylistResponse struct {
	ID         string `json:"id"`
	SnapshotID string `json:"snapshot_id"`
	URI        string `json:"uri"`
	// More fields...
}

type Tracks struct {
	Total int `json:"total"`
	Items []struct {
		URI string `json:"uri"`
	}
}

type SearchForItemsResponse struct {
	Tracks *Tracks `json:"tracks"`
	// More fields...
}

type AddItemsToPlaylistResponse struct {
	SnapshotID string `json:"snapshot_id"`
}

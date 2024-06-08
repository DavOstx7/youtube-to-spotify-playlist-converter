package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang/internal/common/errors"
	"golang/internal/config"
	"golang/internal/requests"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const MinValidTrackURIs = 1

var APIConfig = &config.GetAPIConfig().Spotify

func GetAccessToken(tokenConfig *config.SpotifyTokenConfig, code string) *GetAccessTokenResponse {
	respBody := &GetAccessTokenResponse{}
	requests.Post(APIConfig.URLs.Token,
		requests.Headers(map[string]string{
			"Authorization": "Basic " + encodeBasicAuthCredentials(tokenConfig.ClientID, tokenConfig.ClientSecret),
			"Content-Type":  "application/x-www-form-urlencoded",
		}),
		requests.Body(bytes.NewBufferString(url.Values{
			"grant_type":   {"authorization_code"},
			"code":         {code},
			"redirect_uri": {tokenConfig.RedirectURI},
		}.Encode())),
		requests.ExpectedStatusCodes([]int{http.StatusOK}),
		requests.MaxAttempts(1),
		requests.SaveResponseBody(&respBody),
	)
	return respBody
}

func GetUserProfile(accessToken string) *GetUserProfileResponse {
	respBody := &GetUserProfileResponse{}
	requests.Get(APIConfig.URLs.UserProfile,
		requests.Headers(map[string]string{
			"Authorization": "Bearer " + accessToken,
		}),
		requests.ExpectedStatusCodes([]int{http.StatusOK}),
		requests.MaxAttempts(1),
		requests.SaveResponseBody(respBody),
	)
	return respBody
}

func CreatePlaylist(accessToken string, userID string, playlistInfo config.SpotifyPlaylistInfo) *CreatePlaylistResponse {
	body, err := json.Marshal(playlistInfo)
	if err != nil {
		panic(err)
	}

	respBody := &CreatePlaylistResponse{}
	requests.Post(strings.Replace(APIConfig.URLs.Playlists, "{user_id}", userID, 1),
		requests.Headers(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + accessToken,
		}),
		requests.Body(bytes.NewReader(body)),
		requests.ExpectedStatusCodes([]int{http.StatusCreated}),
		requests.SaveResponseBody(respBody),
	)
	return respBody
}

func SearchForTracks(accessToken string, query string, limit int) *Tracks {
	return searchForItems(accessToken, query, TrackItem, limit).Tracks
}

func AddTracksToPlaylist(accessToken string, playlistID string, trackURIs []string, position int) string {
	if err := validateTrackURIsSize(trackURIs); err != nil {
		panic(err)
	}
	return addItemsToPlaylist(accessToken, playlistID, trackURIs, position).SnapshotID
}

func NewAuthorizationQueryParams(clientID string, redirectURI string) *AuthorizationQueryParams {
	return &AuthorizationQueryParams{
		ClientID:     clientID,
		ResponseType: "code",
		RedirectURI:  redirectURI,
		Scope:        APIConfig.AuthorizationScopes,
	}
}

func encodeBasicAuthCredentials(clientID string, clientSecret string) string {
	return base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
}

func validateTrackURIsSize(trackURIs []string) error {
	trackURIsSize := len(trackURIs)
	if trackURIsSize >= MinValidTrackURIs && trackURIsSize <= APIConfig.MaxTracksPerRequest {
		return nil
	}

	validRange := fmt.Sprintf("%d-%d", MinValidTrackURIs, APIConfig.MaxTracksPerRequest)
	return &errors.ValidationError{
		Message: fmt.Sprintf("size of track uris ({%d}}) is not in the valid range of %s", trackURIsSize, validRange),
	}
}

func searchForItems(accessToken string, query string, itemType ItemType, limit int) *SearchForItemsResponse {
	respBody := &SearchForItemsResponse{}
	requests.Get(APIConfig.URLs.Search,
		requests.Headers(map[string]string{
			"Authorization": "Bearer " + accessToken,
		}),
		requests.QueryParams(map[string]string{"q": query, "type": string(itemType), "limit": strconv.Itoa(limit)}),
		requests.ExpectedStatusCodes([]int{http.StatusOK}),
		requests.SaveResponseBody(respBody),
	)
	return respBody
}

func addItemsToPlaylist(accessToken string, playlistID string, itemURIs []string, position int) *AddItemsToPlaylistResponse {
	body, err := json.Marshal(map[string]interface{}{"uris": itemURIs, "position": position})
	if err != nil {
		panic(err)
	}

	respBody := &AddItemsToPlaylistResponse{}
	requests.Post(strings.Replace(APIConfig.URLs.Tracks, "{playlist_id}", playlistID, 1),
		requests.Headers(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + accessToken,
		}),
		requests.Body(bytes.NewReader(body)),
		requests.ExpectedStatusCodes([]int{http.StatusOK}),
		requests.SaveResponseBody(respBody),
	)
	return respBody
}

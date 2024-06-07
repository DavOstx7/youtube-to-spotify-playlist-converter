package youtube

import (
	"fmt"
	"golang/internal/common/errors"
	"golang/internal/config"
	"golang/internal/requests"
	"net/http"
)

const MinValidMaxResults = 1

var APIConfig = &config.GetAPIConfig().YouTube

func RequestPlaylistItemsPage(queryParams *PlaylistItemsQueryParams) *GetPlaylistItemsResponse {
	respBody := &GetPlaylistItemsResponse{}
	requests.Get(APIConfig.URLs.PlaylistItems,
		requests.QueryParams(queryParams),
		requests.ExpectedStatusCodes([]int{http.StatusOK}),
		requests.SaveResponseBody(respBody),
	)
	return respBody
}

func GetPlaylistQueryParams(apiKey string, playlistID string, maxResults int) *PlaylistItemsQueryParams {
	if err := validateMaxResultsValue(maxResults); err != nil {
		panic(err)
	}

	return &PlaylistItemsQueryParams{Key: apiKey, Part: "snippet", PlaylistID: playlistID, MaxResults: maxResults}
}

func validateMaxResultsValue(maxResults int) error {
	if maxResults >= MinValidMaxResults && maxResults <= APIConfig.MaxItemsPerRequest {
		return nil
	}

	validRange := fmt.Sprintf("%d-%d", MinValidMaxResults, APIConfig.MaxItemsPerRequest)
	return &errors.ValidationError{
		Message: fmt.Sprintf("value of max results ({%d}}) is not in the valid range of %s", maxResults, validRange),
	}
}

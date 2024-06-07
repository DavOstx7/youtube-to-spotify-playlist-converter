package requests

import (
	"fmt"
	"github.com/google/go-querystring/query"
	"golang/internal/common/errors"
	"golang/internal/common/types"
	"golang/internal/utils"
	"net/http"
	"net/url"
	"reflect"
)

func BuildURLWithQuery(url string, queryParams types.Any) (string, error) {
	rawQuery, err := BuildRawQuery(queryParams)
	if err != nil {
		return "", err
	}

	return url + "?" + rawQuery, nil
}

func BuildRawQuery(queryParams types.Any) (string, error) {
	kind := utils.ReflectUnderlyingKind(queryParams)
	
	switch kind {
	case reflect.Struct:
		return buildRawQueryFromStruct(queryParams)
	case reflect.Map:
		return buildRawQueryFromMap(queryParams)
	default:
		return "", fmt.Errorf("input type '%s' is not supported in raw query conversion", kind)
	}
}

func ValidateResponseStatusCode(resp *http.Response, expectedStatusCodes []int) error {
	if utils.SliceContains(expectedStatusCodes, resp.StatusCode) {
		return nil
	}

	return errors.ValidationError{
		Message: fmt.Sprintf("Invalid response status '%s' from url '%s'", resp.Status, resp.Request.URL),
	}
}

func buildRawQueryFromStruct(queryParams types.Any) (string, error) {
	values, err := query.Values(queryParams)
	if err != nil {
		return "", err
	}
	return values.Encode(), nil
}

func buildRawQueryFromMap(queryParams types.Any) (string, error) {
	queryMap, ok := queryParams.(map[string]string)
	if !ok {
		return "", fmt.Errorf("only map of type 'map[string]string' is supported in raw query conversion")
	}

	values := url.Values{}
	for key, value := range queryMap {
		values.Add(key, value)
	}
	return values.Encode(), nil
}

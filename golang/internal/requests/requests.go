package requests

import (
	"encoding/json"
	"fmt"
	"github.com/avast/retry-go"
	"io"
	"log/slog"
	"net/http"
)

var client = &http.Client{Timeout: DefaultRequestTimeout}

func Get(url string, options ...Option) {
	config := NewConfig(options...)
	Do("GET", url, config)
}

func Post(url string, options ...Option) {
	config := NewConfig(options...)
	Do("POST", url, config)
}

func Do(method string, url string, config *Config) error {
	retryOptions := []retry.Option{
		retry.DelayType(retry.BackOffDelay),
		retry.Attempts(config.MaxAttempts),
		retry.Delay(DefaultInitialRetryDelay),
		retry.MaxDelay(DefaultMaxRetryDelay),
		retry.OnRetry(func(n uint, err error) {
			slog.Warn(fmt.Sprintf("Retry #%d: %s\n", n, err))
		}),
	}

	err := retry.Do(func() error {
		req, err := http.NewRequest(method, url, config.Body)
		if err != nil {
			return err
		}

		if config.QueryParams != nil {
			rawQuery, err := BuildRawQuery(config.QueryParams)
			if err != nil {
				return err
			}
			req.URL.RawQuery = rawQuery
		}

		if config.Headers != nil {
			for key, value := range config.Headers {
				req.Header.Set(key, value)
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if config.LogResponse {
			slog.Debug(fmt.Sprintf("%s %s [%d]: %s", req.Method, req.URL.String(), resp.StatusCode, respBody))
		}

		if config.ExpectedStatusCodes != nil {
			if err := ValidateResponseStatusCode(resp, config.ExpectedStatusCodes); err != nil {
				return err
			}
		}

		if config.SaveResponseBody != nil {
			if err := json.Unmarshal(respBody, config.SaveResponseBody); err != nil {
				return err
			}
		}

		return nil
	}, retryOptions...)

	if err != nil && config.PanicOnError {
		panic(err)
	}

	return err
}
package requests

import (
	"bytes"
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
	err := do(method, url, config)

	if err != nil && config.PanicOnError {
		panic(err)
	}

	return err
}

func do(method string, url string, config *Config) error {
	var rawBody []byte

	if config.Body != nil {
		var err error
		rawBody, err = io.ReadAll(config.Body)
		if err != nil {
			return err
		}
	}

	retryOptions := []retry.Option{
		retry.DelayType(retry.CombineDelay(
			retry.BackOffDelay,
			retry.RandomDelay,
		)),
		retry.Attempts(config.MaxAttempts),
		retry.Delay(DefaultRetryDelay),
		retry.MaxDelay(DefaultMaxRetryDelay),
		retry.MaxJitter(DefaultMaxRetryJitter),
		retry.OnRetry(func(n uint, err error) {
			slog.Warn(fmt.Sprintf("Retry #%d: %s", n, err))
		}),
	}

	return retry.Do(func() error {
		var body io.Reader
		if rawBody != nil {
			body = bytes.NewReader(rawBody)
		}

		req, err := http.NewRequest(method, url, body)
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
			if err := ValidateResponseStatusCode(resp, respBody, config.ExpectedStatusCodes); err != nil {
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
}

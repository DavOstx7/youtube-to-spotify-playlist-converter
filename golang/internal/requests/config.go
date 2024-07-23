package requests

import (
	"golang/internal/common/types"
	"io"
)

type Option func(*Config)

type Config struct {
	Body                io.Reader
	QueryParams         types.Any
	Headers             types.StringMap
	ExpectedStatusCodes []int
	MaxAttempts         uint
	LogResponse         bool
	SaveResponseBody    types.Any
	PanicOnError        bool
}

func NewConfig(options ...Option) *Config {
	config := &Config{
		Body:                DefaultBody,
		QueryParams:         DefaultQueryParams,
		Headers:             DefaultHeaders,
		ExpectedStatusCodes: DefaultExpectedStatusCodes,
		MaxAttempts:         DefaultMaxRequestAttempts,
		LogResponse:         DefaultLogResponse,
		SaveResponseBody:    DefaultSaveResponseBody,
		PanicOnError:        DefaultPanicOnError,
	}

	for _, option := range options {
		option(config)
	}

	return config
}

func Body(body io.Reader) Option {
	return func(c *Config) {
		c.Body = body
	}
}

func QueryParams(queryParams types.Any) Option {
	return func(c *Config) {
		c.QueryParams = queryParams
	}
}

func Headers(headers types.StringMap) Option {
	return func(c *Config) {
		c.Headers = headers
	}
}

func ExpectedStatusCodes(expectedStatusCodes []int) Option {
	return func(c *Config) {
		c.ExpectedStatusCodes = expectedStatusCodes
	}
}

func MaxAttempts(maxAttempts uint) Option {
	return func(c *Config) {
		c.MaxAttempts = maxAttempts
	}
}

func LogResponse(logResponse bool) Option {
	return func(c *Config) {
		c.LogResponse = logResponse
	}
}

func SaveResponseBody(saveResponseBody types.Any) Option {
	return func(c *Config) {
		c.SaveResponseBody = saveResponseBody
	}
}

func PanicOnError(panicOnError bool) Option {
	return func(c *Config) {
		c.PanicOnError = panicOnError
	}
}

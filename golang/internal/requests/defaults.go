package requests

import (
	"golang/internal/common/types"
	"io"
	"time"
)

var (
	DefaultRequestTimeout      time.Duration   = 2000 * time.Millisecond
	DefaultBody                io.Reader       = nil
	DefaultQueryParams         types.Any       = nil
	DefaultHeaders             types.StringMap = nil
	DefaultExpectedStatusCodes []int           = nil
	DefaultLogResponse         bool            = true
	DefaultMaxRequestAttempts  uint            = 3
	DefaultSaveResponseBody    types.Any       = nil
	DefaultPanicOnError        bool            = true
	DefaultInitialRetryDelay   time.Duration   = 1000 * time.Millisecond
	DefaultMaxRetryDelay       time.Duration   = 15000 * time.Millisecond
)

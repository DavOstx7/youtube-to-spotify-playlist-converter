package utils

import (
	"fmt"
	"log/slog"
	"strings"
)

func ConvertLevelToSlogLevel(level string) (slog.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("given level '%s' is invalid", level)
	}
}

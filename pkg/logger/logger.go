package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

func init() {
	// Create or open app.log file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	// Create structured logger that writes to file only
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewJSONHandler(file, opts)
	log = slog.New(handler)
}

// GetLogger returns the global logger instance
func GetLogger() *slog.Logger {
	return log
}

// Close closes the logger (cleanup if needed)
func Close() error {
	// slog doesn't require explicit closing, but we can add cleanup here if needed
	return nil
}

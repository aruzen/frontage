package log

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func Init(debug bool) {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo}
	if debug {
		opts.Level = slog.LevelDebug
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, opts))

	slog.SetDefault(logger)
}

package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"

	"github.com/joho/godotenv"
)

func run(ctx context.Context, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("load config file: %w", err)
	}

	setLogger(stderr, os.Getenv("LOG_LEVEL"))

	return nil
}

func main() {
	if err := run(context.Background(), os.Stderr); err != nil {
		slog.
			With("err", err).
			Error("run main")
		os.Exit(1)
	}
}

func setLogger(stderr io.Writer, levelStr string) {
	levelStr = strings.ToLower(levelStr)

	lvl := slog.LevelDebug

	switch levelStr {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "error":
		lvl = slog.LevelError
	default:
	}

	logger := slog.New(
		slog.NewTextHandler(stderr, &slog.HandlerOptions{Level: lvl}),
	)

	slog.SetDefault(logger)
}

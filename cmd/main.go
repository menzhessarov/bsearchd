package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"

	"github.com/menzhessarov/bsearchd"
)

func run(ctx context.Context, stderr io.Writer) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("load config file: %w", err)
	}

	setLogger(stderr, os.Getenv("LOG_LEVEL"))

	conformation, err := strconv.Atoi(os.Getenv("CONFORMATION"))
	if err != nil {
		return fmt.Errorf("parse conformation: %w", err)
	}

	store := bsearchd.NewStore(os.Getenv("INPUT_FILE"), conformation)

	err = store.Load()
	if err != nil {
		return fmt.Errorf("load store: %w", err)
	}

	port := os.Getenv("HTTP_PORT")

	s := bsearchd.NewHTTPServer(port, store)

	s.RegisterRoutes()

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		slog.InfoContext(ctx, "server started", "port", port)

		if err := s.Server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-gctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return s.Server.Shutdown(ctx)
	})

	err = g.Wait()
	if err != nil {
		return fmt.Errorf("goroutines finished: %w", err)
	}

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

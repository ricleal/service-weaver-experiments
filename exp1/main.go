package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/lmittmann/tint"

	"exp/api"
	"exp/count"
	"exp/reverse"
)

//go:generate ../../cmd/weaver/weaver generate

func Logger(w io.Writer, levelAsString string) *slog.Logger {
	var level slog.Level

	switch strings.ToLower(levelAsString) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "Error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      level,
			TimeFormat: time.TimeOnly,
		}),
	)

	return logger
}

type app struct {
	weaver.Implements[weaver.Main]
	c weaver.Ref[count.Counter]
	r weaver.Ref[reverse.Reverser]
}

func serve(ctx context.Context, app *app) error {
	log := Logger(os.Stderr, os.Getenv("LOG_LEVEL"))

	rev := api.Reverser{
		Reverser: app.r,
	}
	count := api.Counter{
		Counter: app.c,
	}

	mux := http.NewServeMux()
	mux.Handle("GET /reverser", rev)
	mux.Handle("GET /counter", count)

	log.Info("listening on :8080")
	return http.ListenAndServe(":8080", mux)
}

func main() {
	ctx := context.Background()

	log := Logger(os.Stderr, os.Getenv("LOG_LEVEL"))

	if err := weaver.Run(ctx, serve); err != nil {
		log.Error("failed to run service", tint.Err(err))
	}
}

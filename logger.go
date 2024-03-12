package main

import (
	"context"
	"log/slog"
	"os"
)

type LogHandler struct {
	slog.Handler
}

func (lh *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestId, ok := ctx.Value(RequestIdKey).(string); ok {
		r.AddAttrs(slog.String("requestId", requestId))
	}

	return lh.Handler.Handle(ctx, r)
}
func customizeLogging() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	wrapper := LogHandler{handler}
	logger := slog.New(&wrapper)
	slog.SetDefault(logger)

	slog.Info("logging was customized, requestId was added")
}

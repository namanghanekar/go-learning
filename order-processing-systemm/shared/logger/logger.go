package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger() {

	Log = slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			nil,
		),
	)
}

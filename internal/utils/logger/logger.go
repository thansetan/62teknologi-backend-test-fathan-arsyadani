package logger

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func New() *slog.Logger {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		AddSource:  true,
		TimeFormat: "02-Jan-2006 15:04 -0700",
	}))

	return logger
}

package bootstrap

import (
	"os"

	"golang.org/x/exp/slog"
)

var Log *slog.Logger

func InitLog() {
	textHandler := slog.NewTextHandler(os.Stdout)
	Log = slog.New(textHandler)
}

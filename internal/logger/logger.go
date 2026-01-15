package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/thdxg/llog/internal/config"
	"github.com/spf13/cobra"
)

type Logger struct {
	file *os.File
}

func Load(cfg *config.Config, lg *Logger) error {
	file, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file; %w", err)
	}

	handler := slog.NewTextHandler(file, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

	slog.SetDefault(slog.New(handler))

	lg.file = file

	return nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func LogCmdStart(cmd *cobra.Command) {
	slog.Info(
		"command started",
		cmdAttr(cmd),
	)
}

func LogCmdComplete(cmd *cobra.Command) {
	slog.Info(
		"command completed",
		cmdAttr(cmd),
	)
}

func cmdAttr(cmd *cobra.Command) slog.Attr {
	return slog.Group(
		"command",
		slog.String("name", cmd.Name()),
		slog.Any("args", cmd.Flags().Args()),
	)
}

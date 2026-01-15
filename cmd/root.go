/*
Copyright © 2025 Ethan Lee <ethantlee21@gmail.com>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/thdxg/llog/internal/config"
	_db "github.com/thdxg/llog/internal/db"
	"github.com/thdxg/llog/internal/handler"
	"github.com/thdxg/llog/internal/logger"
	"github.com/spf13/cobra"
)

var (
	cfg = &config.Config{}
	db  = &_db.DB{}
	lg  = &logger.Logger{}
)

var rootCmd = &cobra.Command{
	Use:               "llog",
	Short:             "Life log",
	Long:              `Record your fleeting moments with llog.`,
	PersistentPreRunE: handler.Init(cfg, db, lg),
	SilenceUsage:      true,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		slog.Error("failed to execute command", slog.Any("error", err))
		os.Exit(1)
	}
	_ = lg.Close()
}

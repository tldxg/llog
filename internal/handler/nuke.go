package handler

import (
	"github.com/thdxg/llog/internal/config"
	_db "github.com/thdxg/llog/internal/db"
	"github.com/thdxg/llog/internal/logger"
	"github.com/spf13/cobra"
)

func Nuke(cfg *config.Config, db *_db.DB, opts *NukeOpts) HandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		logger.LogCmdStart(cmd)
		defer logger.LogCmdComplete(cmd)

		if err := db.Nuke(); err != nil {
			return err
		}

		return nil
	}
}

type NukeOpts struct{}

func (o *NukeOpts) applyFlags(cmd *cobra.Command) {}

func (o *NukeOpts) validate(cfg *config.Config, args []string, flags []string) error { return nil }

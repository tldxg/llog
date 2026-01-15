package handler

import (
	"fmt"

	"github.com/thdxg/llog/internal/config"
	_db "github.com/thdxg/llog/internal/db"
	"github.com/thdxg/llog/internal/logger"
	"github.com/spf13/cobra"
)

type HandlerFunc func(cmd *cobra.Command, args []string) error

func Init(cfg *config.Config, db *_db.DB, lg *logger.Logger) HandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		if err := config.Load(cfg); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if err := logger.Load(cfg, lg); err != nil {
			return fmt.Errorf("failed to load logger: %w", err)
		}

		if err := _db.Load(cfg, ctx, db); err != nil {
			return fmt.Errorf("failed to load db: %w", err)
		}

		return nil
	}
}

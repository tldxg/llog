/*
Copyright © 2025 Ethan Lee <ethantlee21@gmail.com>
*/
package cmd

import (
	"github.com/thdxg/llog/internal/handler"
	"github.com/spf13/cobra"
)

var searchOpts = &handler.SearchOpts{}

var searchCmd = &cobra.Command{
	Use:          "search",
	Short:        "Fuzzy find log entries",
	Long:         `Fuzzy find log entries.`,
	Args:         cobra.ArbitraryArgs,
	PreRunE:      handler.ValidateOptions(cfg, searchOpts),
	RunE:         handler.Search(cfg, db, searchOpts),
	SilenceUsage: true,
}

func init() {
	handler.ApplyFlags(searchCmd, searchOpts)
	rootCmd.AddCommand(searchCmd)
}

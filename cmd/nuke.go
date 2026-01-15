/*
Copyright © 2025 Ethan Lee <ethantlee21@gmail.com>
*/
package cmd

import (
	"github.com/thdxg/llog/internal/handler"
	"github.com/spf13/cobra"
)

var nukeOpts = &handler.NukeOpts{}

var nukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Delete all data",
	Long: `Delete all data.

This will nuke the entire database and clear all data.
This action cannot be undone.`,
	Args:         cobra.ArbitraryArgs,
	PreRunE:      handler.ValidateOptions(cfg, nukeOpts),
	RunE:         handler.Nuke(cfg, db, nukeOpts),
	SilenceUsage: true,
}

func init() {
	handler.ApplyFlags(nukeCmd, nukeOpts)
	rootCmd.AddCommand(nukeCmd)
}

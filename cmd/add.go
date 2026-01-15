/*
Copyright © 2025 Ethan Lee <ethantlee21@gmail.com>
*/
package cmd

import (
	"github.com/thdxg/llog/internal/handler"
	"github.com/spf13/cobra"
)

var addOpts = &handler.AddOpts{}

var addCmd = &cobra.Command{
	Use:   "add [body]",
	Short: "Add log entries",
	Long: `Add log entries.

You can create multiple entries by providing as multiple arguments.
`,
	Args:         cobra.MinimumNArgs(1),
	PreRunE:      handler.ValidateOptions(cfg, addOpts),
	RunE:         handler.Add(cfg, db, addOpts),
	SilenceUsage: true,
}

func init() {
	handler.ApplyFlags(addCmd, addOpts)
	rootCmd.AddCommand(addCmd)
}

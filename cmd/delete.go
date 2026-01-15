/*
Copyright © 2025 Ethan Lee <ethantlee21@gmail.com>
*/
package cmd

import (
	"github.com/thdxg/llog/internal/handler"
	"github.com/spf13/cobra"
)

var deleteOpts = &handler.DeleteOpts{}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete log entries",
	Long: `Delete log entries.

You can delete specific entries by providing IDs as arguments.
Alternatively, you can delete entries by using flags to filter by date range or limit the count.

Note: Providing an ID is mutually exclusive with using any flags.`,
	Args:         cobra.ArbitraryArgs,
	PreRunE:      handler.ValidateOptions(cfg, deleteOpts),
	RunE:         handler.Delete(cfg, db, deleteOpts),
	SilenceUsage: true,
}

func init() {
	handler.ApplyFlags(deleteCmd, deleteOpts)
	rootCmd.AddCommand(deleteCmd)
}

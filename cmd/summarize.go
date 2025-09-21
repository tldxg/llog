/*
Copyright © 2025 Ethan Lee <ethantlee21@gmail.com>
*/
package cmd

import (
	"github.com/ethn1ee/llog/internal/handler"
	"github.com/spf13/cobra"
)

var summarizeOpts = &handler.SummarizeOpts{}

var summarizeCmd = &cobra.Command{
	Use:          "summarize",
	Short:        "Summarize log entries",
	Long:         `Summarize log entries.`,
	Args:         cobra.NoArgs,
	PreRunE:      handler.ValidateOptions(cfg, summarizeOpts),
	RunE:         handler.Summarize(cfg, db, summarizeOpts),
	SilenceUsage: true,
}

func init() {
	handler.ApplyFlags(summarizeCmd, summarizeOpts)
	rootCmd.AddCommand(summarizeCmd)
}

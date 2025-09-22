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
	Use:   "summarize",
	Short: "Summarize log entries with Gemini",
	Long: `Summarize log entries with Gemini.

You can use flags to filter entries to summarize, similarly to the 'get' command.

Note: 'GEMINI_API_KEY' environment variable must be set to use this feature.`,
	Args:         cobra.NoArgs,
	PreRunE:      handler.ValidateOptions(cfg, summarizeOpts),
	RunE:         handler.Summarize(cfg, db, summarizeOpts),
	SilenceUsage: true,
}

func init() {
	handler.ApplyFlags(summarizeCmd, summarizeOpts)
	rootCmd.AddCommand(summarizeCmd)
}

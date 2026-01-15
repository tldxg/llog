package handler

import (
	"github.com/thdxg/llog/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Opts interface {
	applyFlags(cmd *cobra.Command)
	validate(cfg *config.Config, args []string, flags []string) error
}

func ValidateOptions(cfg *config.Config, opts Opts) HandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		flags := make([]string, 0)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			flags = append(flags, f.Name)
		})
		return opts.validate(cfg, args, flags)
	}
}

func ApplyFlags(cmd *cobra.Command, opts Opts) {
	opts.applyFlags(cmd)
}

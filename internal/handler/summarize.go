package handler

import (
	"errors"
	"fmt"
	"slices"

	"github.com/ethn1ee/llog/internal/config"
	_db "github.com/ethn1ee/llog/internal/db"
	"github.com/ethn1ee/llog/internal/logger"
	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

func Summarize(cfg *config.Config, db *_db.DB, opts *SummarizeOpts) HandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		logger.LogCmdStart(cmd)
		defer logger.LogCmdComplete(cmd)

		ctx := cmd.Context()

		entries, err := getWithOpts(cfg, db, cmd, &GetOpts{
			Time:  opts.Time,
			Limit: opts.Limit,
			All:   false,
		})
		if err != nil {
			return err
		}

		client, err := genai.NewClient(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to initialize Gemini client: %w", err)
		}

		thinkingBudget := int32(0)
		prompt := fmt.Sprintf("You are a summarizer who concisely summarizes a list of timestamped activities that the user provides. You should make any other output except the summary itself, and address in first-person in the summary. Summarize the entries provided: %+v", entries)

		res, err := client.Models.GenerateContent(
			ctx,
			"gemini-2.5-flash",
			genai.Text(prompt),
			&genai.GenerateContentConfig{
				ThinkingConfig: &genai.ThinkingConfig{
					ThinkingBudget: &thinkingBudget,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to generate content: %w", err)
		}

		fmt.Println(res.Text())

		return nil
	}
}

type SummarizeOpts struct {
	Time  timeOpts
	Limit int
}

func (o *SummarizeOpts) applyFlags(cmd *cobra.Command) {
	o.Time.applyFlags(cmd)
	cmd.Flags().IntVarP(&(o.Limit), "limit", "n", 10, "number of entries to summarize")
}

func (o *SummarizeOpts) validate(cfg *config.Config, args []string, flags []string) error {
	if len(args) > 0 && len(flags) > 0 {
		return errors.New(flagIdMutexError)
	}

	if err := o.Time.validate(cfg, args, flags); err != nil {
		return err
	}

	if !o.Time.fromTime.IsZero() {
		if slices.Contains(flags, "limit") {
			return fmt.Errorf(flagMutexError, "from", "limit")
		}
		o.Limit = -1
	}

	if !o.Time.toTime.IsZero() {
		if slices.Contains(flags, "limit") {
			return fmt.Errorf(flagMutexError, "to", "limit")
		}
		o.Limit = -1
	}

	return nil
}

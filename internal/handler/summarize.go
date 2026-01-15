package handler

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/thdxg/llog/internal/config"
	_db "github.com/thdxg/llog/internal/db"
	"github.com/thdxg/llog/internal/logger"
	"github.com/thdxg/llog/internal/model"
	"github.com/thdxg/llog/internal/view"
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
			All:   opts.All,
		})
		if err != nil {
			return err
		}

		spinner := view.StartSpinner(fmt.Sprintf("summarizing %d entries...", len(entries)))
		summaries, err := summarizeEntries(cfg, ctx, entries)
		if err != nil {
			return fmt.Errorf("failed to summarize entries: %w", err)
		}
		view.StopSpinner(spinner)

		view.PrintSummaries(cfg, summaries)
		view.PrintSummarize(len(entries))

		return nil
	}
}

func summarizeEntries(cfg *config.Config, ctx context.Context, entries []model.Entry) (string, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Gemini client: %w", err)
	}

	thinkingBudget := int32(0)
	prompt := fmt.Sprintf(`
			<prompt>
			You are a summarizer who concisely summarizes a list of timestamped activities that the user provides.
			Each entry contains the ID, timestamp, and body.
			You should summarize each day, and return as an array of day summaries containing the date and summary.
			The format of the date must be %s.
			The summary for each day must be one sentence.
			The sentences must be impersonal and agentless, meaning that it should start with a verb in past tense, not a pronoun or agent of the action (e.g. "went shopping").
			You may omit some minor details in summary in order to fit into one sentence concisely.
			</prompt>

			<entries>
			%+v
			</entries>`,
		cfg.DateLayout, entries,
	)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"date":    {Type: genai.TypeString},
					"summary": {Type: genai.TypeString},
				},
				PropertyOrdering: []string{"date", "summary"},
			},
		},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudget,
		},
	}
	res, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return res.Text(), nil
}

type SummarizeOpts struct {
	Time  timeOpts
	Limit int
	All   bool
}

func (o *SummarizeOpts) applyFlags(cmd *cobra.Command) {
	o.Time.applyFlags(cmd)
	cmd.Flags().IntVarP(&(o.Limit), "limit", "n", 10, "number of entries to summarize")
	cmd.Flags().BoolVarP(&(o.All), "all", "a", false, "summarize all entries")
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
		if slices.Contains(flags, "all") {
			return fmt.Errorf(flagMutexError, "from", "all")
		}
		o.Limit = -1
	}

	if !o.Time.toTime.IsZero() {
		if slices.Contains(flags, "limit") {
			return fmt.Errorf(flagMutexError, "to", "limit")
		}
		if slices.Contains(flags, "all") {
			return fmt.Errorf(flagMutexError, "to", "all")
		}
		o.Limit = -1
	}

	if o.All {
		if slices.Contains(flags, "limit") {
			return fmt.Errorf(flagMutexError, "all", "limit")
		}
		o.Limit = -1
	}

	return nil
}

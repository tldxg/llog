package handler

import (
	"fmt"
	"time"

	"github.com/thdxg/llog/internal/config"
	"github.com/spf13/cobra"
)

type timeOpts struct {
	Today     bool
	Yesterday bool
	From      string
	To        string

	fromTime time.Time
	toTime   time.Time
}

func (o *timeOpts) applyFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(
		&(o.Today), "today", "t", false, "scope to today's entries",
	)
	cmd.Flags().BoolVarP(
		&(o.Yesterday), "yesterday", "y", false, "scope to yesterday's entries",
	)
	cmd.Flags().StringVar(
		&(o.From), "from", "", "scope start date in YYYY-MM-DD format (inclusive)",
	)
	cmd.Flags().StringVar(
		&(o.To), "to", "", "scope end date in YYYY-MM-DD format (exclusive)",
	)
}

func (o *timeOpts) validate(cfg *config.Config, args []string, flags []string) error {
	// mutual exclusion checks
	if o.Today && o.Yesterday {
		return fmt.Errorf(flagMutexError, "today", "yesterday")
	}
	if o.Today && o.From != "" {
		return fmt.Errorf(flagMutexError, "today", "from")
	}
	if o.Today && o.To != "" {
		return fmt.Errorf(flagMutexError, "today", "to")
	}
	if o.Yesterday && o.From != "" {
		return fmt.Errorf(flagMutexError, "yesterday", "from")
	}
	if o.Yesterday && o.To != "" {
		return fmt.Errorf(flagMutexError, "yesterday", "to")
	}

	// set fromTime and toTime
	if o.Today {
		now := time.Now()
		o.fromTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		o.toTime = o.fromTime.Add(24 * time.Hour)
		return nil
	}

	if o.Yesterday {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		o.fromTime = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
		o.toTime = o.fromTime.Add(24 * time.Hour)
		return nil
	}

	if o.From != "" {
		from, err := time.Parse(cfg.DateLayout, o.From)
		if err != nil {
			return fmt.Errorf(dateParseError, err)
		}
		o.fromTime = from
	}

	if o.To != "" {
		to, err := time.Parse(cfg.DateLayout, o.To)
		if err != nil {
			return fmt.Errorf(dateParseError, err)
		}
		o.toTime = to
	}

	if !o.fromTime.IsZero() && !o.toTime.IsZero() && o.fromTime.After(o.toTime) {
		return fmt.Errorf(dateRangeError, "'from' cannot be after 'to'")
	}

	return nil
}

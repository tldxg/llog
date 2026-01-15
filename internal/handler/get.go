package handler

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/thdxg/llog/internal/config"
	_db "github.com/thdxg/llog/internal/db"
	"github.com/thdxg/llog/internal/logger"
	"github.com/thdxg/llog/internal/model"
	"github.com/thdxg/llog/internal/view"
	"github.com/spf13/cobra"
)

func Get(cfg *config.Config, db *_db.DB, opts *GetOpts) HandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		logger.LogCmdStart(cmd)
		defer logger.LogCmdComplete(cmd)

		var entries []model.Entry
		var err error

		if len(args) > 0 {
			entries, err = getWithArgs(cfg, db, cmd, args)
		} else {
			entries, err = getWithOpts(cfg, db, cmd, opts)
		}

		if err != nil {
			return err
		}

		view.PrintEntries(cfg, entries)
		view.PrintGet(len(entries))

		return nil
	}
}

func getWithArgs(cfg *config.Config, db *_db.DB, cmd *cobra.Command, args []string) ([]model.Entry, error) {
	ctx := cmd.Context()
	ids := make([]uint64, len(args))

	for i, arg := range args {
		id, err := strconv.ParseUint(arg, 10, 32)
		if err != nil {
			return nil, fmt.Errorf(idParseError, err)
		}

		if id < 1 || id > cfg.Internal.MaxEntryId {
			return nil, fmt.Errorf(idVoidError, id)
		}

		ids[i] = id
	}

	entries, err := db.Entry.Get(ctx, db.Entry.WithIds(ids), -1)
	if err != nil {
		return nil, fmt.Errorf(dbGetEntryError, err)
	}

	return entries, nil
}

func getWithOpts(_ *config.Config, db *_db.DB, cmd *cobra.Command, opts *GetOpts) ([]model.Entry, error) {
	ctx := cmd.Context()

	var entries []model.Entry
	var err error

	from, to := opts.Time.fromTime, opts.Time.toTime

	if !from.IsZero() || !to.IsZero() {
		entries, err = db.Entry.Get(ctx, db.Entry.WithRange(from, to), opts.Limit)
	} else {
		entries, err = db.Entry.Get(ctx, nil, opts.Limit)
	}

	if err != nil {
		return nil, fmt.Errorf(dbGetEntryError, err)
	}

	return entries, nil
}

type GetOpts struct {
	Time  timeOpts
	Limit int
	All   bool
}

func (o *GetOpts) applyFlags(cmd *cobra.Command) {
	o.Time.applyFlags(cmd)
	cmd.Flags().IntVarP(&(o.Limit), "limit", "n", 10, "number of entries to return")
	cmd.Flags().BoolVarP(&(o.All), "all", "a", false, "retrurn all entries")
}

func (o *GetOpts) validate(cfg *config.Config, args []string, flags []string) error {
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

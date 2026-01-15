package handler

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/thdxg/llog/internal/config"
	_db "github.com/thdxg/llog/internal/db"
	"github.com/thdxg/llog/internal/logger"
	"github.com/thdxg/llog/internal/view"
	"github.com/spf13/cobra"
)

func Delete(cfg *config.Config, db *_db.DB, opts *DeleteOpts) HandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		logger.LogCmdStart(cmd)
		defer logger.LogCmdComplete(cmd)

		var count int
		var err error

		if len(args) > 0 {
			count, err = deleteWithArgs(cfg, db, cmd, args)
		} else {
			count, err = deleteWithOpts(cfg, db, cmd, opts)
		}

		if err != nil {
			return err
		}

		view.PrintDelete(count)

		return nil
	}
}

func deleteWithArgs(cfg *config.Config, db *_db.DB, cmd *cobra.Command, args []string) (int, error) {
	ctx := cmd.Context()
	ids := make([]uint64, len(args))

	for i, arg := range args {
		id, err := strconv.ParseUint(arg, 10, 32)
		if err != nil {
			return 0, fmt.Errorf(idParseError, err)
		}

		if id < 1 || id > cfg.Internal.MaxEntryId {
			return 0, fmt.Errorf(idVoidError, id)
		}

		ids[i] = id
	}

	count, err := db.Entry.Delete(ctx, db.Entry.WithIds(ids), -1)
	if err != nil {
		return 0, fmt.Errorf(dbDeleteEntryError, err)
	}

	return count, nil
}

func deleteWithOpts(_ *config.Config, db *_db.DB, cmd *cobra.Command, opts *DeleteOpts) (int, error) {
	ctx := cmd.Context()

	var count int
	var err error

	from, to := opts.Time.fromTime, opts.Time.toTime

	if !from.IsZero() || !to.IsZero() {
		count, err = db.Entry.Delete(ctx, db.Entry.WithRange(from, to), opts.Limit)
	} else {
		count, err = db.Entry.Delete(ctx, nil, opts.Limit)
	}

	if err != nil {
		return 0, fmt.Errorf(dbDeleteEntryError, err)
	}

	return count, nil
}

type DeleteOpts struct {
	Time        timeOpts
	Limit       int
	All         bool
	Interactive bool
}

func (o *DeleteOpts) applyFlags(cmd *cobra.Command) {
	o.Time.applyFlags(cmd)
	cmd.Flags().IntVarP(&(o.Limit), "limit", "n", 10, "number of entries to delete")
	cmd.Flags().BoolVarP(&(o.All), "all", "a", false, "delete all entries")
	// TODO: interactive flag
}

func (o *DeleteOpts) validate(cfg *config.Config, args []string, flags []string) error {
	if len(args) > 0 && len(flags) > 0 {
		return errors.New(flagIdMutexError)
	}

	if len(args) == 0 && len(flags) == 0 {
		return errors.New(noArgsOrFlagsError)
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

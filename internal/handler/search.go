package handler

import (
	"fmt"
	"strings"

	"github.com/thdxg/llog/internal/config"
	_db "github.com/thdxg/llog/internal/db"
	"github.com/thdxg/llog/internal/logger"
	"github.com/thdxg/llog/internal/model"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
)

func Search(cfg *config.Config, db *_db.DB, opts *SearchOpts) HandlerFunc {
	return func(cmd *cobra.Command, args []string) error {
		logger.LogCmdStart(cmd)
		defer logger.LogCmdComplete(cmd)

		ctx := cmd.Context()

		input := strings.Join(args, " ")

		entries, err := db.Entry.Get(ctx, nil, -1)
		if err != nil {
			return fmt.Errorf(dbGetEntryError, err)
		}

		res := fuzzy.FindFrom(input, entrySlice(entries))

		for _, r := range res {

			fmt.Println(entries[r.Index])
		}

		return nil
	}
}

type entrySlice []model.Entry

// implemented specifically to use the fuzzy package FindFrom
func (e entrySlice) String(i int) string {
	return e[i].Body
}

// implemented specifically to use the fuzzy package FindFrom
func (e entrySlice) Len() int {
	return len(e)
}

type SearchOpts struct{}

func (o *SearchOpts) applyFlags(cmd *cobra.Command) {}

func (o *SearchOpts) validate(cfg *config.Config, args []string, flags []string) error {
	return nil
}

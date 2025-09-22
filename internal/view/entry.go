package view

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"github.com/ethn1ee/llog/internal/config"
	"github.com/ethn1ee/llog/internal/model"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

func PrintEntries(cfg *config.Config, entries []model.Entry) {
	data := make([][]string, len(entries))
	for i, e := range entries {
		data[i] = []string{
			color.HiCyanString(strconv.FormatUint(e.ID, 10)),
			color.HiBlackString(e.CreatedAt.Format(cfg.TimeLayout)),
			e.Body,
		}
	}

	printTable(data)
}

func PrintSummaries(_ *config.Config, summaryRaw string) {
	var summaries []model.Summary
	if err := json.Unmarshal([]byte(summaryRaw), &summaries); err != nil {
		fmt.Printf("failed to unmarshal summary: %s", summaryRaw)
	}

	data := make([][]string, len(summaries))
	for i, s := range summaries {
		data[i] = []string{
			color.HiBlackString(s.Date),
			s.Summary,
		}
	}

	printTable(data)
}

func StartSpinner(suffix string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "\n"
	s.Suffix = color.HiBlackString(" " + suffix)
	s.Start()
	return s
}

func StopSpinner(s *spinner.Spinner) {
	s.Stop()
}

func printTable(data [][]string) {
	symbols := tw.NewSymbolCustom("minimal").
		WithRow("").
		WithColumn("")

	table := tablewriter.NewTable(
		os.Stdout,
		tablewriter.WithRenderer(
			renderer.NewBlueprint(tw.Rendition{
				Symbols: symbols,
			}),
		),
	)

	_ = table.Bulk(data)
	_ = table.Render()
}

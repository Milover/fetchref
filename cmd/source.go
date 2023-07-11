package cmd

import (
	"github.com/Milover/fetchpaper/internal/fetch"
	"github.com/spf13/cobra"
)

var sourceCmd = &cobra.Command{
	Use:   "source <DOI...>",
	Short: "Fetch paper(s) from Sci-Hub from supplied DOI(s).",
	Long:  "Fetch paper(s) from Sci-Hub from supplied DOI(s).",
	Args:  cobra.MinimumNArgs(1),
	RunE:  source,
}

func source(cmd *cobra.Command, args []string) error {
	return fetch.Fetch(fetch.SourceMode, args)
}

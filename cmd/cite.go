package cmd

import (
	"github.com/Milover/fetchref/internal/fetch"
	"github.com/spf13/cobra"
)

var citeCmd = &cobra.Command{
	Use:           "cite <DOI...>",
	Short:         "Fetch citation(s) from Crossref from supplied DOI(s).",
	Long:          "Fetch citation(s) from Crossref from supplied DOI(s).",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          cite,
}

func cite(cmd *cobra.Command, args []string) error {
	return fetch.Fetch(fetch.CiteMode, args)
}

package cmd

import (
	"github.com/Milover/fetchref/internal/fetch"
	"github.com/spf13/cobra"
)

var sourceCmd = &cobra.Command{
	Use:           "source <DOI...>",
	Short:         "Fetch reference(s) from Sci-Hub/Libgen from supplied DOI(s)/ISBN(s).",
	Long:          "Fetch reference(s) from Sci-Hub/Libgen from supplied DOI(s)/ISBN(s).",
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          source,
}

func source(cmd *cobra.Command, args []string) error {
	return fetch.Fetch(fetch.SourceMode, args)
}

package cmd

import (
	"github.com/Milover/fetchref/internal/fetch"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {
	return fetch.Fetch(fetch.DefaultMode, args)
}

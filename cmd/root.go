package cmd

import (
	"fmt"
	"os"

	"github.com/Milover/fetchpaper/internal/fetch"
	"github.com/Milover/fetchpaper/internal/metainfo"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "fetchpaper [options...] <DOI...>",
	Short:   "Fetch paper(s) from Sci-Hub from supplied DOI(s).",
	Long:    `Fetch paper(s) from Sci-Hub from supplied DOI(s).`,
	Version: metainfo.Version,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := fetch.Fetch(args); err != nil {
			os.Exit(1)
		}
	},
	DisableFlagsInUseLine: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().DurationVar(
		&fetch.GlobalReqTimeout,
		"timeout",
		fetch.GlobalReqTimeout,
		"HTTP request timeout",
	)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

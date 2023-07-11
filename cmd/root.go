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
	Use:     "fetchpaper <DOI...>",
	Short:   "Fetch paper(s) and citations from supplied DOI(s).",
	Long:    "Fetch paper(s) and citations from supplied DOI(s).",
	Version: metainfo.Version,
	Args:    cobra.MinimumNArgs(1),
	RunE:    run,
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
	rootCmd.AddCommand(sourceCmd, citeCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().DurationVar(
		&fetch.GlobalReqTimeout,
		"timeout",
		fetch.GlobalReqTimeout,
		"HTTP request timeout",
	)
	rootCmd.PersistentFlags().BoolVar(
		&fetch.NoUserAgent,
		"no-user-agent",
		false,
		"omit User-Agent header from HTTP requests",
	)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(
		&fetch.CitationFileName,
		"cite-file",
		"o",
		fetch.CitationFileName,
		"citation output file name, w/o extension",
	)
	citeCmd.Flags().StringVarP(
		&fetch.CitationFileName,
		"cite-file",
		"o",
		fetch.CitationFileName,
		"citation output file name, w/o extension",
	)
	rootCmd.LocalFlags().Var(
		&fetch.CitationFormat,
		"cite-format",
		"article citation output format",
	)
	citeCmd.LocalFlags().Var(
		&fetch.CitationFormat,
		"cite-format",
		"article citation output format",
	)
}

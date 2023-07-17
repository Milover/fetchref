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
	Use:           "fetchpaper <DOI...>",
	Short:         "Fetch paper(s) and citations from supplied DOI(s).",
	Long:          "Fetch paper(s) and citations from supplied DOI(s).",
	Version:       metainfo.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          cobra.MinimumNArgs(1),
	RunE:          run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
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
		&fetch.CiteFileName,
		"cite-file",
		"o",
		fetch.CiteFileName,
		"citation output file name, w/o extension",
	)
	citeCmd.Flags().StringVarP(
		&fetch.CiteFileName,
		"cite-file",
		"o",
		fetch.CiteFileName,
		"citation output file name, w/o extension",
	)
	rootCmd.Flags().Var(
		&fetch.CiteFormat,
		"cite-format",
		"article citation output format",
	)
	citeCmd.Flags().Var(
		&fetch.CiteFormat,
		"cite-format",
		"article citation output format",
	)
	rootCmd.Flags().BoolVar(
		&fetch.CiteAppend,
		"cite-append",
		false,
		"append citations to file instead of overwriting",
	)
	citeCmd.Flags().BoolVar(
		&fetch.CiteAppend,
		"cite-append",
		false,
		"append citations to file instead of overwriting",
	)
	rootCmd.Flags().BoolVar(
		&fetch.CiteSeparate,
		"cite-separate",
		false,
		"write each citation to a different file",
	)
	citeCmd.Flags().BoolVar(
		&fetch.CiteSeparate,
		"cite-separate",
		false,
		"write each citation to a different file",
	)
	rootCmd.MarkFlagsMutuallyExclusive("cite-file", "cite-separate")
	citeCmd.MarkFlagsMutuallyExclusive("cite-file", "cite-separate")
}

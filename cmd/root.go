package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/Milover/fetchpaper/internal/fetch"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fetchpaper [options...] <DOI...>",
	Short: "Fetch paper(s) from Sci-Hub from supplied DOI(s).",
	Long:  `Fetch paper(s) from Sci-Hub from supplied DOI(s).`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if err := fetch.Fetch(args); err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Println("success")
	},
	DisableFlagsInUseLine: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fetchpaper.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

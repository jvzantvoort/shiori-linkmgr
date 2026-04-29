package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("linkmgr %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built:  %s\n", buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

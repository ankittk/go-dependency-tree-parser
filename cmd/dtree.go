package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dtree",
	Short: "Dependency Tree Builder for Go GitHub Repositories",
	Long: `
dtree is a command-line tool that analyzes and caches the Go dependencies.
The tool can be used to analyze the dependencies of a repository before making changes to it, or to cache the dependencies for faster builds.
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose debug logging")
}

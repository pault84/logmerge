package cmd

import (

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "logmerge",
		Short: "A Logmerging tool for portworx",
		Long: `Logmerge is a tool that allows us to merge multiple log files based on date.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(mergeCommand())
}
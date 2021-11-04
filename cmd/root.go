package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "logmerge",
		Short: "A Logmerging tool that merges logfiles with parseable datetime formats in their lines.",
		Long:  `Logmerge is a tool that allows us to merge multiple log files based on date.`,
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Pass files as args. Write to Stdout.
			mergeLogs(args)
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

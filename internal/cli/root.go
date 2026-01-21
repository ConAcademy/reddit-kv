package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "reddit-kv",
	Short: "A key-value store backed by Reddit",
	Long: `reddit-kv uses Reddit as a key-value store backend.

Subreddits become databases, post titles become keys,
and comment trees become values.

This is a proof-of-concept. Please don't use it for anything serious.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(appendCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(keysCmd)
}

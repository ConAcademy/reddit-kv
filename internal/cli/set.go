package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sprite/reddit-kv/pkg/redditkv"
)

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a key to a value",
	Long: `Set a key to a value. If the key already exists, it will be overwritten.

The key becomes a Reddit post title, and the value becomes a comment.`,
	Args: cobra.ExactArgs(2),
	RunE: runSet,
}

func runSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	cfg, err := redditkv.LoadConfig()
	if err != nil {
		return err
	}

	client, err := redditkv.New(*cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	if err := client.Set(key, value); err != nil {
		return fmt.Errorf("failed to set key: %w", err)
	}

	fmt.Printf("OK\n")
	return nil
}

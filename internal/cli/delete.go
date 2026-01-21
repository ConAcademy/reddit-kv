package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sprite/reddit-kv/pkg/redditkv"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <key>",
	Short: "Delete a key",
	Long:  `Delete a key and all its associated values (the post and all comments).`,
	Args:  cobra.ExactArgs(1),
	RunE:  runDelete,
}

func runDelete(cmd *cobra.Command, args []string) error {
	key := args[0]

	cfg, err := redditkv.LoadConfig()
	if err != nil {
		return err
	}

	client, err := redditkv.New(*cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	if err := client.Delete(key); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	fmt.Printf("OK\n")
	return nil
}

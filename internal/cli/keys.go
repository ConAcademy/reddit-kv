package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sprite/reddit-kv/pkg/redditkv"
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "List all keys",
	Long:  `List all keys in the store (all post titles in the subreddit).`,
	Args:  cobra.NoArgs,
	RunE:  runKeys,
}

var flagJSON bool

func init() {
	keysCmd.Flags().BoolVar(&flagJSON, "json", false, "Output as JSON array")
}

func runKeys(cmd *cobra.Command, args []string) error {
	cfg, err := redditkv.LoadConfig()
	if err != nil {
		return err
	}

	client, err := redditkv.New(*cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	keys, err := client.Keys()
	if err != nil {
		return fmt.Errorf("failed to list keys: %w", err)
	}

	if flagJSON {
		output, err := json.MarshalIndent(keys, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal keys: %w", err)
		}
		fmt.Println(string(output))
		return nil
	}

	if len(keys) == 0 {
		fmt.Println("(no keys)")
		return nil
	}

	for _, key := range keys {
		fmt.Println(key)
	}
	return nil
}

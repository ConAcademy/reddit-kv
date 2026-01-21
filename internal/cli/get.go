package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sprite/reddit-kv/pkg/redditkv"
)

var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get the value for a key",
	Long: `Get the value tree for a key.

Returns the value as a JSON tree structure representing the comment hierarchy.`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

var flagRaw bool

func init() {
	getCmd.Flags().BoolVar(&flagRaw, "raw", false, "Output raw value (only works for scalar values)")
}

func runGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	cfg, err := redditkv.LoadConfig()
	if err != nil {
		return err
	}

	client, err := redditkv.New(*cfg)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	value, err := client.Get(key)
	if err != nil {
		return fmt.Errorf("failed to get key: %w", err)
	}

	if flagRaw {
		// Output just the root value
		fmt.Println(value.Value)
		return nil
	}

	// Output as JSON
	output, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

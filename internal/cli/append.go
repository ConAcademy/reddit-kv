package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sprite/reddit-kv/pkg/redditkv"
)

var appendCmd = &cobra.Command{
	Use:   "append <key> <value>",
	Short: "Append a value to an existing key",
	Long: `Append a value to an existing key's tree.

By default, appends as a new top-level comment (sibling to root).
Use --parent to specify a path to append as a child of a specific node.

Path format: comma-separated indices (e.g., "0,1" means second child of first child)`,
	Args: cobra.ExactArgs(2),
	RunE: runAppend,
}

var flagParent string

func init() {
	appendCmd.Flags().StringVar(&flagParent, "parent", "", "Path to parent node (e.g., '0,1')")
}

func runAppend(cmd *cobra.Command, args []string) error {
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

	var parentPath []int
	if flagParent != "" {
		parts := strings.Split(flagParent, ",")
		parentPath = make([]int, len(parts))
		for i, part := range parts {
			idx, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return fmt.Errorf("invalid path: %s", flagParent)
			}
			parentPath[i] = idx
		}
	}

	if err := client.Append(key, value, parentPath); err != nil {
		return fmt.Errorf("failed to append: %w", err)
	}

	fmt.Printf("OK\n")
	return nil
}

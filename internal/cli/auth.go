package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/sprite/reddit-kv/pkg/redditkv"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Configure Reddit API credentials",
	Long: `Configure your Reddit API credentials for reddit-kv.

You need to create a Reddit app first:
1. Go to https://www.reddit.com/prefs/apps
2. Click "create another app..."
3. Choose "script" type
4. Note your client ID (under the app name) and client secret

You also need to create a subreddit to use as your database.`,
	RunE: runAuth,
}

var (
	flagClientID     string
	flagClientSecret string
	flagUsername     string
	flagPassword     string
	flagSubreddit    string
)

func init() {
	authCmd.Flags().StringVar(&flagClientID, "client-id", "", "Reddit app client ID")
	authCmd.Flags().StringVar(&flagClientSecret, "client-secret", "", "Reddit app client secret")
	authCmd.Flags().StringVar(&flagUsername, "username", "", "Reddit username")
	authCmd.Flags().StringVar(&flagPassword, "password", "", "Reddit password")
	authCmd.Flags().StringVar(&flagSubreddit, "subreddit", "", "Subreddit to use as database")
}

func runAuth(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Prompt for any missing values
	if flagClientID == "" {
		fmt.Print("Client ID: ")
		flagClientID, _ = reader.ReadString('\n')
		flagClientID = strings.TrimSpace(flagClientID)
	}

	if flagClientSecret == "" {
		fmt.Print("Client Secret: ")
		flagClientSecret, _ = reader.ReadString('\n')
		flagClientSecret = strings.TrimSpace(flagClientSecret)
	}

	if flagUsername == "" {
		fmt.Print("Reddit Username: ")
		flagUsername, _ = reader.ReadString('\n')
		flagUsername = strings.TrimSpace(flagUsername)
	}

	if flagPassword == "" {
		fmt.Print("Reddit Password: ")
		flagPassword, _ = reader.ReadString('\n')
		flagPassword = strings.TrimSpace(flagPassword)
	}

	if flagSubreddit == "" {
		fmt.Print("Subreddit (without r/): ")
		flagSubreddit, _ = reader.ReadString('\n')
		flagSubreddit = strings.TrimSpace(flagSubreddit)
	}

	// Validate
	if flagClientID == "" || flagClientSecret == "" || flagUsername == "" || flagPassword == "" || flagSubreddit == "" {
		return fmt.Errorf("all fields are required")
	}

	// Remove r/ prefix if present
	flagSubreddit = strings.TrimPrefix(flagSubreddit, "r/")

	cfg := &redditkv.Config{
		ClientID:     flagClientID,
		ClientSecret: flagClientSecret,
		Username:     flagUsername,
		Password:     flagPassword,
		Subreddit:    flagSubreddit,
	}

	// Test the credentials by creating a client
	_, err := redditkv.New(*cfg)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Reddit: %w", err)
	}

	// Save config
	if err := redditkv.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	configPath, _ := redditkv.ConfigPath()
	fmt.Printf("Configuration saved to %s\n", configPath)
	fmt.Println("You can now use reddit-kv commands!")

	return nil
}

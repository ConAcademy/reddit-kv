// Demo program that exercises reddit-kv using the mock backend
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sprite/reddit-kv/pkg/redditkv"
)

func main() {
	fmt.Println("=== reddit-kv Demo (using mock backend) ===\n")

	// Create a mock Reddit API and client
	mock := redditkv.NewMockRedditAPI()
	client := redditkv.NewWithAPI(mock, "demo_subreddit")

	// Demo 1: Basic SET and GET
	fmt.Println("1. Basic SET and GET")
	fmt.Println("   SET greeting \"hello world\"")
	if err := client.Set("greeting", "hello world"); err != nil {
		fatal("Set failed:", err)
	}
	fmt.Println("   OK")

	fmt.Println("   GET greeting")
	val, err := client.Get("greeting")
	if err != nil {
		fatal("Get failed:", err)
	}
	printJSON(val)

	// Demo 2: SET overwrites
	fmt.Println("\n2. SET overwrites existing key")
	fmt.Println("   SET greeting \"hello universe\"")
	if err := client.Set("greeting", "hello universe"); err != nil {
		fatal("Set failed:", err)
	}
	fmt.Println("   OK")

	fmt.Println("   GET greeting")
	val, _ = client.Get("greeting")
	printJSON(val)

	// Demo 3: Multiple keys
	fmt.Println("\n3. Multiple keys")
	fmt.Println("   SET name \"reddit-kv\"")
	client.Set("name", "reddit-kv")
	fmt.Println("   SET version \"0.1.0\"")
	client.Set("version", "0.1.0")
	fmt.Println("   OK")

	fmt.Println("   KEYS")
	keys, _ := client.Keys()
	for _, k := range keys {
		fmt.Printf("   - %s\n", k)
	}

	// Demo 4: APPEND to build a tree
	fmt.Println("\n4. APPEND to build a tree")
	fmt.Println("   SET config \"root\"")
	client.Set("config", "root")

	fmt.Println("   APPEND config \"child1\" (to root)")
	client.Append("config", "child1", []int{0})

	fmt.Println("   APPEND config \"child2\" (to root)")
	client.Append("config", "child2", []int{0})

	fmt.Println("   GET config")
	val, _ = client.Get("config")
	printJSON(val)

	// Demo 5: DELETE
	fmt.Println("\n5. DELETE")
	fmt.Println("   DELETE version")
	if err := client.Delete("version"); err != nil {
		fatal("Delete failed:", err)
	}
	fmt.Println("   OK")

	fmt.Println("   KEYS")
	keys, _ = client.Keys()
	for _, k := range keys {
		fmt.Printf("   - %s\n", k)
	}

	// Demo 6: Error handling
	fmt.Println("\n6. Error handling")
	fmt.Println("   GET nonexistent")
	_, err = client.Get("nonexistent")
	if err != nil {
		fmt.Printf("   Error (expected): %v\n", err)
	}

	fmt.Println("   EXISTS greeting")
	exists, _ := client.Exists("greeting")
	fmt.Printf("   %v\n", exists)

	fmt.Println("   EXISTS nonexistent")
	exists, _ = client.Exists("nonexistent")
	fmt.Printf("   %v\n", exists)

	fmt.Println("\n=== Demo complete ===")
}

func printJSON(v interface{}) {
	out, _ := json.MarshalIndent(v, "   ", "  ")
	fmt.Println("  ", string(out))
}

func fatal(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s %v\n", msg, err)
	os.Exit(1)
}

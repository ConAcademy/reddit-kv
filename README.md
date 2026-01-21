# reddit-kv

A key-value store that uses Reddit as its backend. Because why not.

## Overview

**reddit-kv** treats Reddit as a database:
- **Subreddit** = Database
- **Post title** = Key
- **Comments** = Values (with tree structure support)

This is a proof-of-concept inspired by a [HackerNews typo](https://news.ycombinator.com/item?id=42702344) about "Reddit as a Key-Value store" in a thread about [vibe coding](https://news.ycombinator.com/item?id=42691243).

## Installation

```bash
go install github.com/yourusername/reddit-kv@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/reddit-kv
cd reddit-kv
go build -o reddit-kv ./cmd/reddit-kv
```

## Setup

1. Create a Reddit app at https://www.reddit.com/prefs/apps
   - Choose "script" type for personal use
   - Note your client ID and client secret

2. Create a subreddit to use as your database

3. Configure reddit-kv:

```bash
reddit-kv auth --client-id=YOUR_CLIENT_ID --client-secret=YOUR_CLIENT_SECRET
```

## Usage

### Basic Operations

```bash
# Set a key (creates post with comment)
reddit-kv set mykey "hello world"

# Get a key (returns value tree)
reddit-kv get mykey

# Append to a key (adds sibling comment)
reddit-kv append mykey "another value"

# Append as child of specific node
reddit-kv append mykey "child value" --parent=0,1

# Delete a key (deletes the post)
reddit-kv delete mykey

# List all keys
reddit-kv keys
```

### Value Structure

Values are stored as Reddit comment trees. The structure you get back reflects the comment hierarchy:

**Scalar** (single comment):
```json
{
  "value": "hello world",
  "children": []
}
```

**Array** (linear thread):
```json
{
  "value": "first",
  "children": [
    {
      "value": "second",
      "children": [
        {
          "value": "third",
          "children": []
        }
      ]
    }
  ]
}
```

**Tree** (branching comments):
```json
{
  "value": "root",
  "children": [
    {"value": "branch1", "children": []},
    {"value": "branch2", "children": [...]}
  ]
}
```

### Library Usage

```go
package main

import (
    "fmt"
    "github.com/yourusername/reddit-kv/pkg/redditkv"
)

func main() {
    client, err := redditkv.New(redditkv.Config{
        Subreddit:    "mykvstore",
        ClientID:     "your-client-id",
        ClientSecret: "your-client-secret",
        Username:     "your-username",
        Password:     "your-password",
    })
    if err != nil {
        panic(err)
    }

    // Set a value
    err = client.Set("mykey", "hello world")

    // Get a value
    tree, err := client.Get("mykey")
    fmt.Println(tree.Value) // "hello world"

    // Append to the tree
    err = client.Append("mykey", "new value", nil) // nil = append as sibling

    // Append as child
    path := []int{0, 1}
    err = client.Append("mykey", "child value", path)

    // Delete
    err = client.Delete("mykey")

    // List keys
    keys, err := client.Keys()
}
```

## Limitations

- **Speed**: This is Reddit, not Redis. Expect API latency.
- **Rate limits**: Reddit API has rate limits (~60 requests/minute)
- **Storage**: Subject to Reddit's post/comment limits
- **Terms of Service**: This almost certainly violates Reddit's ToS. Use for educational purposes only.

## License

MIT

## Disclaimer

This is a joke project. Please don't actually use Reddit as your production database. The badgers will come for you.

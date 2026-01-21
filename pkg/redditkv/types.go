package redditkv

// ValueNode represents a node in the value tree.
// A single comment becomes a scalar (no children).
// A linear thread becomes an array (each node has one child).
// A branching thread becomes a tree (nodes can have multiple children).
type ValueNode struct {
	Value    string      `json:"value"`
	Children []ValueNode `json:"children"`
}

// Config holds the configuration for the reddit-kv client.
type Config struct {
	// Reddit OAuth2 credentials
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Username     string `json:"username"`
	Password     string `json:"password"`

	// Target subreddit (the "database")
	Subreddit string `json:"subreddit"`

	// OAuth tokens (managed internally)
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenExpiry  string `json:"token_expiry,omitempty"`
}

// Client is the interface for reddit-kv operations.
// This interface allows for easy mocking in tests.
type Client interface {
	// Set creates or overwrites a key with a scalar value.
	// If the key exists, it will be deleted and recreated.
	Set(key, value string) error

	// Get retrieves the value tree for a key.
	// Returns nil if the key does not exist.
	Get(key string) (*ValueNode, error)

	// Append adds a value to an existing key's tree.
	// If parentPath is nil, appends as a sibling to the root.
	// If parentPath is provided, appends as a child of the specified node.
	Append(key, value string, parentPath []int) error

	// Delete removes a key and all its values.
	Delete(key string) error

	// Keys returns all keys in the store.
	Keys() ([]string, error)

	// Exists checks if a key exists.
	Exists(key string) (bool, error)
}

// KeyNotFoundError is returned when a key does not exist.
type KeyNotFoundError struct {
	Key string
}

func (e *KeyNotFoundError) Error() string {
	return "key not found: " + e.Key
}

// InvalidPathError is returned when a parent path is invalid.
type InvalidPathError struct {
	Path []int
}

func (e *InvalidPathError) Error() string {
	return "invalid path"
}

package redditkv

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// KVClient implements the Client interface using Reddit as a backend.
type KVClient struct {
	api       RedditAPI
	subreddit string
	ctx       context.Context
}

// New creates a new reddit-kv client with the given configuration.
func New(cfg Config) (*KVClient, error) {
	api, err := NewRedditAPI(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Reddit API client: %w", err)
	}

	return &KVClient{
		api:       api,
		subreddit: cfg.Subreddit,
		ctx:       context.Background(),
	}, nil
}

// NewWithAPI creates a new reddit-kv client with a custom RedditAPI implementation.
// This is useful for testing with a mock.
func NewWithAPI(api RedditAPI, subreddit string) *KVClient {
	return &KVClient{
		api:       api,
		subreddit: subreddit,
		ctx:       context.Background(),
	}
}

// Set creates or overwrites a key with a scalar value.
func (c *KVClient) Set(key, value string) error {
	// Check if key exists
	existingPost, err := c.findPostByTitle(key)
	if err != nil {
		return fmt.Errorf("failed to check existing key: %w", err)
	}

	// Delete existing post if found (overwrite behavior)
	if existingPost != nil {
		if err := c.api.DeletePost(c.ctx, existingPost.ID); err != nil {
			return fmt.Errorf("failed to delete existing key: %w", err)
		}
	}

	// Create new post with empty body (title is the key)
	submitted, err := c.api.SubmitPost(c.ctx, c.subreddit, key, "")
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	// Add the value as a comment
	_, err = c.api.SubmitComment(c.ctx, submitted.FullID, value)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// Get retrieves the value tree for a key.
func (c *KVClient) Get(key string) (*ValueNode, error) {
	post, err := c.findPostByTitle(key)
	if err != nil {
		return nil, fmt.Errorf("failed to find key: %w", err)
	}
	if post == nil {
		return nil, &KeyNotFoundError{Key: key}
	}

	// Get post with comments
	postAndComments, err := c.api.GetPost(c.ctx, post.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	// Convert comments to ValueNode tree
	if len(postAndComments.Comments) == 0 {
		return nil, &KeyNotFoundError{Key: key}
	}

	// The root of our value tree is the first top-level comment
	// If there are multiple top-level comments, we need to handle that
	return commentsToValueTree(postAndComments.Comments), nil
}

// Append adds a value to an existing key's tree.
func (c *KVClient) Append(key, value string, parentPath []int) error {
	post, err := c.findPostByTitle(key)
	if err != nil {
		return fmt.Errorf("failed to find key: %w", err)
	}
	if post == nil {
		return &KeyNotFoundError{Key: key}
	}

	// Get post with comments to find the parent
	postAndComments, err := c.api.GetPost(c.ctx, post.ID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	var parentID string

	if parentPath == nil || len(parentPath) == 0 {
		// Append as new top-level comment (sibling to root)
		parentID = post.FullID
	} else {
		// Navigate to the parent comment
		comment, err := navigateToComment(postAndComments.Comments, parentPath)
		if err != nil {
			return &InvalidPathError{Path: parentPath}
		}
		parentID = comment.FullID
	}

	_, err = c.api.SubmitComment(c.ctx, parentID, value)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}

	return nil
}

// Delete removes a key and all its values.
func (c *KVClient) Delete(key string) error {
	post, err := c.findPostByTitle(key)
	if err != nil {
		return fmt.Errorf("failed to find key: %w", err)
	}
	if post == nil {
		return &KeyNotFoundError{Key: key}
	}

	if err := c.api.DeletePost(c.ctx, post.ID); err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// Keys returns all keys in the store.
func (c *KVClient) Keys() ([]string, error) {
	posts, err := c.api.ListNewPosts(c.ctx, c.subreddit, &reddit.ListOptions{
		Limit: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	keys := make([]string, len(posts))
	for i, post := range posts {
		keys[i] = post.Title
	}

	return keys, nil
}

// Exists checks if a key exists.
func (c *KVClient) Exists(key string) (bool, error) {
	post, err := c.findPostByTitle(key)
	if err != nil {
		return false, err
	}
	return post != nil, nil
}

// findPostByTitle searches for a post with the exact title (key).
func (c *KVClient) findPostByTitle(title string) (*reddit.Post, error) {
	posts, err := c.api.SearchPosts(c.ctx, c.subreddit, title)
	if err != nil {
		return nil, err
	}

	// Find exact match
	for _, post := range posts {
		if post.Title == title {
			return post, nil
		}
	}

	return nil, nil
}

// commentsToValueTree converts Reddit comments to our ValueNode tree structure.
func commentsToValueTree(comments []*reddit.Comment) *ValueNode {
	if len(comments) == 0 {
		return nil
	}

	// If there's only one top-level comment, it's the root
	if len(comments) == 1 {
		return commentToValueNode(comments[0])
	}

	// Multiple top-level comments: create a synthetic root
	// This represents the case where SET was called, then APPEND was called
	// multiple times at the root level
	root := &ValueNode{
		Value:    comments[0].Body,
		Children: make([]ValueNode, 0, len(comments)-1),
	}

	// Add remaining top-level comments as children of the first
	for i := 1; i < len(comments); i++ {
		root.Children = append(root.Children, *commentToValueNode(comments[i]))
	}

	// Also add the first comment's replies as children
	if len(comments[0].Replies.Comments) > 0 {
		for _, reply := range comments[0].Replies.Comments {
			root.Children = append(root.Children, *commentToValueNode(reply))
		}
	}

	return root
}

// commentToValueNode converts a single Reddit comment (with replies) to a ValueNode.
func commentToValueNode(comment *reddit.Comment) *ValueNode {
	node := &ValueNode{
		Value:    comment.Body,
		Children: make([]ValueNode, 0, len(comment.Replies.Comments)),
	}

	for _, reply := range comment.Replies.Comments {
		node.Children = append(node.Children, *commentToValueNode(reply))
	}

	return node
}

// navigateToComment follows a path through the comment tree.
func navigateToComment(comments []*reddit.Comment, path []int) (*reddit.Comment, error) {
	if len(path) == 0 || len(comments) == 0 {
		return nil, fmt.Errorf("invalid path")
	}

	idx := path[0]
	if idx < 0 || idx >= len(comments) {
		return nil, fmt.Errorf("path index out of range")
	}

	comment := comments[idx]

	if len(path) == 1 {
		return comment, nil
	}

	// Navigate deeper
	return navigateToComment(comment.Replies.Comments, path[1:])
}

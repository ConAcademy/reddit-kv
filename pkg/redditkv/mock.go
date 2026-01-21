package redditkv

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// MockRedditAPI is a mock implementation of RedditAPI for testing.
type MockRedditAPI struct {
	mu       sync.RWMutex
	posts    map[string]*mockPost // postID -> post
	comments map[string]*reddit.Comment // commentID -> comment
	idCounter int
}

type mockPost struct {
	post     *reddit.Post
	comments []*reddit.Comment // top-level comments
}

// NewMockRedditAPI creates a new mock Reddit API for testing.
func NewMockRedditAPI() *MockRedditAPI {
	return &MockRedditAPI{
		posts:    make(map[string]*mockPost),
		comments: make(map[string]*reddit.Comment),
	}
}

func (m *MockRedditAPI) nextID() string {
	m.idCounter++
	return fmt.Sprintf("%d", m.idCounter)
}

func (m *MockRedditAPI) SubmitPost(ctx context.Context, subreddit, title, text string) (*reddit.Submitted, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.nextID()
	fullID := "t3_" + id
	now := reddit.Timestamp{Time: time.Now()}

	post := &reddit.Post{
		ID:            id,
		FullID:        fullID,
		Title:         title,
		Body:          text,
		SubredditName: subreddit,
		Created:       &now,
	}

	m.posts[id] = &mockPost{
		post:     post,
		comments: []*reddit.Comment{},
	}

	return &reddit.Submitted{
		ID:     id,
		FullID: fullID,
		URL:    fmt.Sprintf("https://reddit.com/r/%s/comments/%s", subreddit, id),
	}, nil
}

func (m *MockRedditAPI) GetPost(ctx context.Context, postID string) (*reddit.PostAndComments, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	mp, ok := m.posts[postID]
	if !ok {
		return nil, fmt.Errorf("post not found: %s", postID)
	}

	return &reddit.PostAndComments{
		Post:     mp.post,
		Comments: mp.comments,
	}, nil
}

func (m *MockRedditAPI) DeletePost(ctx context.Context, postID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.posts[postID]; !ok {
		return fmt.Errorf("post not found: %s", postID)
	}

	delete(m.posts, postID)
	return nil
}

func (m *MockRedditAPI) SubmitComment(ctx context.Context, parentID, text string) (*reddit.Comment, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := m.nextID()
	fullID := "t1_" + id
	now := reddit.Timestamp{Time: time.Now()}

	comment := &reddit.Comment{
		ID:       id,
		FullID:   fullID,
		Body:     text,
		ParentID: parentID,
		Created:  &now,
		Replies:  reddit.Replies{Comments: []*reddit.Comment{}},
	}

	m.comments[id] = comment

	// Add to parent (either post or comment)
	if len(parentID) > 3 && parentID[:3] == "t3_" {
		// Parent is a post
		postID := parentID[3:]
		if mp, ok := m.posts[postID]; ok {
			mp.comments = append(mp.comments, comment)
		}
	} else if len(parentID) > 3 && parentID[:3] == "t1_" {
		// Parent is a comment
		parentCommentID := parentID[3:]
		if parentComment, ok := m.comments[parentCommentID]; ok {
			parentComment.Replies.Comments = append(parentComment.Replies.Comments, comment)
		}
	}

	return comment, nil
}

func (m *MockRedditAPI) ListNewPosts(ctx context.Context, subreddit string, opts *reddit.ListOptions) ([]*reddit.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var posts []*reddit.Post
	for _, mp := range m.posts {
		if mp.post.SubredditName == subreddit {
			posts = append(posts, mp.post)
		}
	}

	// Apply limit if specified
	if opts != nil && opts.Limit > 0 && len(posts) > opts.Limit {
		posts = posts[:opts.Limit]
	}

	return posts, nil
}

func (m *MockRedditAPI) SearchPosts(ctx context.Context, subreddit, query string) ([]*reddit.Post, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var posts []*reddit.Post
	for _, mp := range m.posts {
		if mp.post.SubredditName == subreddit && mp.post.Title == query {
			posts = append(posts, mp.post)
		}
	}

	return posts, nil
}

// Helper methods for testing

// GetPostCount returns the number of posts in the mock.
func (m *MockRedditAPI) GetPostCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.posts)
}

// GetCommentCount returns the number of comments in the mock.
func (m *MockRedditAPI) GetCommentCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.comments)
}

// Reset clears all data in the mock.
func (m *MockRedditAPI) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.posts = make(map[string]*mockPost)
	m.comments = make(map[string]*reddit.Comment)
	m.idCounter = 0
}

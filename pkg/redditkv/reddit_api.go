package redditkv

import (
	"context"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// RedditAPI is an interface wrapping the Reddit API operations we need.
// This allows for easy mocking in tests.
type RedditAPI interface {
	// Post operations
	SubmitPost(ctx context.Context, subreddit, title, text string) (*reddit.Submitted, error)
	GetPost(ctx context.Context, postID string) (*reddit.PostAndComments, error)
	DeletePost(ctx context.Context, postID string) error

	// Comment operations
	SubmitComment(ctx context.Context, parentID, text string) (*reddit.Comment, error)

	// Subreddit operations
	ListNewPosts(ctx context.Context, subreddit string, opts *reddit.ListOptions) ([]*reddit.Post, error)
	SearchPosts(ctx context.Context, subreddit, query string) ([]*reddit.Post, error)
}

// redditAPIClient wraps the go-reddit client to implement RedditAPI.
type redditAPIClient struct {
	client *reddit.Client
}

// NewRedditAPI creates a new Reddit API client with the given credentials.
func NewRedditAPI(cfg Config) (RedditAPI, error) {
	credentials := reddit.Credentials{
		ID:       cfg.ClientID,
		Secret:   cfg.ClientSecret,
		Username: cfg.Username,
		Password: cfg.Password,
	}

	client, err := reddit.NewClient(credentials, reddit.WithUserAgent("reddit-kv/0.1.0"))
	if err != nil {
		return nil, err
	}

	return &redditAPIClient{client: client}, nil
}

func (r *redditAPIClient) SubmitPost(ctx context.Context, subreddit, title, text string) (*reddit.Submitted, error) {
	submitted, _, err := r.client.Post.SubmitText(ctx, reddit.SubmitTextRequest{
		Subreddit: subreddit,
		Title:     title,
		Text:      text,
	})
	return submitted, err
}

func (r *redditAPIClient) GetPost(ctx context.Context, postID string) (*reddit.PostAndComments, error) {
	post, _, err := r.client.Post.Get(ctx, postID)
	return post, err
}

func (r *redditAPIClient) DeletePost(ctx context.Context, postID string) error {
	_, err := r.client.Post.Delete(ctx, postID)
	return err
}

func (r *redditAPIClient) SubmitComment(ctx context.Context, parentID, text string) (*reddit.Comment, error) {
	comment, _, err := r.client.Comment.Submit(ctx, parentID, text)
	return comment, err
}

func (r *redditAPIClient) ListNewPosts(ctx context.Context, subreddit string, opts *reddit.ListOptions) ([]*reddit.Post, error) {
	posts, _, err := r.client.Subreddit.NewPosts(ctx, subreddit, opts)
	return posts, err
}

func (r *redditAPIClient) SearchPosts(ctx context.Context, subreddit, query string) ([]*reddit.Post, error) {
	// Search for posts with exact title match
	posts, _, err := r.client.Subreddit.SearchPosts(ctx, query, subreddit, &reddit.ListPostSearchOptions{
		ListPostOptions: reddit.ListPostOptions{
			ListOptions: reddit.ListOptions{
				Limit: 100,
			},
		},
		Sort: "new",
	})
	return posts, err
}

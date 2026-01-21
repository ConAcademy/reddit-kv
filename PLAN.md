# PLAN.md - Implementation Plan

## Project Structure

```
reddit-kv/
├── cmd/
│   └── reddit-kv/
│       └── main.go           # CLI entry point
├── pkg/
│   └── redditkv/
│       ├── client.go         # Main client implementation
│       ├── auth.go           # OAuth2 authentication
│       ├── types.go          # Data types (ValueNode, etc.)
│       ├── reddit_api.go     # Low-level Reddit API calls
│       └── config.go         # Configuration management
├── internal/
│   └── cli/
│       ├── root.go           # Root command
│       ├── auth.go           # auth command
│       ├── set.go            # set command
│       ├── get.go            # get command
│       ├── append.go         # append command
│       ├── delete.go         # delete command
│       └── keys.go           # keys command
├── go.mod
├── go.sum
├── README.md
├── AGENTS.md
├── MEMORY.md
└── PLAN.md
```

## Implementation Phases

### Phase 1: Project Setup and Core Types

- [ ] Initialize Go module
- [ ] Define core types in `pkg/redditkv/types.go`
  - `ValueNode` struct
  - `Config` struct
  - `Client` interface
- [ ] Set up CLI skeleton with cobra

**Deliverable**: Compiles, `reddit-kv --help` works

### Phase 2: OAuth2 Authentication

- [ ] Implement OAuth2 password grant flow in `pkg/redditkv/auth.go`
- [ ] Token storage/refresh in `~/.config/reddit-kv/config.json`
- [ ] Implement `auth` CLI command
- [ ] Test authentication with Reddit API

**Deliverable**: Can authenticate and store tokens

### Phase 3: Reddit API Layer

- [ ] Implement low-level Reddit API in `pkg/redditkv/reddit_api.go`
  - `createPost(subreddit, title string) (postID string, error)`
  - `createComment(postID, parentID, body string) (commentID string, error)`
  - `getComments(postID string) (CommentTree, error)`
  - `listPosts(subreddit string) ([]Post, error)`
  - `deletePost(postID string) error`
- [ ] Proper User-Agent header
- [ ] Error handling and rate limit awareness

**Deliverable**: Can make raw Reddit API calls

### Phase 4: KV Operations - SET and GET

- [ ] Implement `Set(key, value)` in `pkg/redditkv/client.go`
  - Search for existing post with title = key
  - If exists: decide behavior (error or overwrite)
  - Create post with title = key
  - Create comment with body = value
- [ ] Implement `Get(key)`
  - Search for post with title = key
  - Fetch comment tree
  - Transform to `ValueNode` structure
- [ ] Implement `set` and `get` CLI commands

**Deliverable**: `reddit-kv set foo bar` and `reddit-kv get foo` work

### Phase 5: KV Operations - APPEND

- [ ] Implement `Append(key, value, parentPath)`
  - Find post by key
  - Navigate to parent comment using path
  - Create comment as reply
- [ ] Implement `append` CLI command with `--parent` flag
- [ ] Path parsing (comma-separated integers)

**Deliverable**: Can build tree structures

### Phase 6: KV Operations - DELETE and KEYS

- [ ] Implement `Delete(key)`
  - Find post by key
  - Delete post (and all comments with it)
- [ ] Implement `Keys()`
  - List posts in subreddit
  - Return titles
- [ ] Implement `delete` and `keys` CLI commands

**Deliverable**: Full CRUD operations work

### Phase 7: Polish and Error Handling

- [ ] Better error messages
- [ ] Rate limit handling with backoff
- [ ] Input validation
- [ ] Help text and examples for CLI
- [ ] Optional: key caching for faster lookups

**Deliverable**: Production-ish quality

## Dependencies

```go
require (
    github.com/spf13/cobra v1.8.0      // CLI framework
    golang.org/x/oauth2 v0.16.0        // OAuth2 client
)
```

## Configuration File

Location: `~/.config/reddit-kv/config.json`

```json
{
  "client_id": "your-client-id",
  "client_secret": "your-client-secret",
  "username": "your-reddit-username",
  "password": "your-reddit-password",
  "subreddit": "your-kv-subreddit",
  "access_token": "current-access-token",
  "refresh_token": "refresh-token",
  "token_expiry": "2024-01-15T12:00:00Z"
}
```

## Testing Strategy

### Manual Testing

1. Create test subreddit (e.g., `r/testkvstore123`)
2. Run through all operations manually
3. Verify in Reddit UI that posts/comments appear correctly

### Automated Testing (if time permits)

- Mock Reddit API responses
- Unit tests for tree transformation logic
- Integration tests against real Reddit (with test subreddit)

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Reddit API changes | Use stable v1 API, handle errors gracefully |
| Rate limiting | Implement backoff, warn user |
| ToS violation | Clear disclaimers, educational use only |
| Token expiry | Implement automatic refresh |

## Next Steps

Ready to begin **Phase 1: Project Setup and Core Types**.

Start with:
1. `go mod init`
2. Create directory structure
3. Define types
4. Set up cobra CLI skeleton

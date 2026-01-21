# MEMORY.md - Design Decisions and Concepts

## Core Concepts

### Data Model

| Reddit Concept | KV Concept | Notes |
|----------------|------------|-------|
| Subreddit | Database | User creates and controls the subreddit |
| Post Title | Key | User-specified key, must be unique within subreddit |
| Post ID | Internal ID | Reddit-generated, used for API operations |
| Comment Tree | Value | Hierarchical value structure |

### Value Structure

Values are represented as trees, mapped directly from Reddit's comment structure:

```go
type ValueNode struct {
    Value    string      `json:"value"`
    Children []ValueNode `json:"children"`
}
```

**Interpretations:**
- **Scalar**: Single node, empty children
- **Array**: Linear chain (each node has 0 or 1 child)
- **Tree**: Branching structure (nodes can have multiple children)

### Key Constraints

- Keys are Reddit post titles
- Reddit allows duplicate post titles, but we treat keys as unique
- On `SET`, if key exists, we **overwrite** (delete old post, create new)
- Keys are strings, no size limit specified (Reddit's title limit applies)

## Design Decisions

### DD-001: Language Choice - Go

**Decision**: Use Go for implementation.

**Rationale**:
- Human preference
- Excellent for CLI tools (single binary distribution)
- Good HTTP/OAuth libraries (golang.org/x/oauth2)
- Clean concurrency model for API operations

### DD-002: Authentication - OAuth2 Password Grant

**Decision**: Use Reddit's OAuth2 "script" app type with password grant.

**Rationale**:
- Simplest flow for personal/CLI use
- No callback URL needed
- User provides credentials once, we manage token refresh
- Tokens stored locally in config file

### DD-003: All Values Are Strings

**Decision**: No type system. All keys and values are strings.

**Rationale**:
- Matches memcached simplicity
- Users can encode/decode as needed (JSON, base64, etc.)
- Reduces complexity

### DD-004: No TTL/Expiration

**Decision**: Keys don't expire.

**Rationale**:
- Keeps implementation simple
- Reddit doesn't auto-delete posts anyway
- User can manually delete if needed

### DD-005: User-Controlled Subreddit

**Decision**: User must create and provide their own subreddit.

**Rationale**:
- Avoids polluting public subreddits
- User has full control/moderation
- Can be private subreddit for "security"

### DD-006: Tree Representation (Option B)

**Decision**: Return values as nested JSON structure preserving tree hierarchy.

**Rationale**:
- Preserves full information from Reddit's comment structure
- More flexible than flat array
- Natural mapping to comment trees

## API Design

### CLI Commands

| Command | Description | Reddit Operation |
|---------|-------------|------------------|
| `auth` | Configure OAuth credentials | N/A |
| `set <key> <value>` | Create/update key with value | Create post + comment |
| `get <key>` | Retrieve value tree | Fetch post + comments |
| `append <key> <value> [--parent=path]` | Add value to tree | Add comment |
| `delete <key>` | Remove key | Delete post |
| `keys` | List all keys | List posts in subreddit |

### Library Interface

```go
type Client interface {
    Set(key, value string) error
    Get(key string) (*ValueNode, error)
    Append(key, value string, parentPath []int) error
    Delete(key string) error
    Keys() ([]string, error)
}
```

### Path Notation

For `--parent` flag, paths are comma-separated indices:
- `0` = first child of root
- `0,1` = second child of first child of root
- Empty/nil = append as new root-level sibling

## Reddit API Notes

### Endpoints Used

- `POST /api/v1/access_token` - OAuth token
- `POST /api/submit` - Create post
- `POST /api/comment` - Add comment
- `GET /r/{subreddit}/comments/{post_id}` - Get post + comments
- `GET /r/{subreddit}/new` - List posts
- `POST /api/del` - Delete post

### Rate Limits

- ~60 requests per minute for OAuth apps
- Should implement exponential backoff
- Consider caching post ID lookups

### User Agent

Reddit requires a unique User-Agent:
```
reddit-kv/0.1.0 (by /u/yourusername)
```

## Resolved Questions

1. **Upsert behavior**: `SET` overwrites existing keys (delete + recreate)

## Open Questions

1. **Caching**: Should we cache post ID lookups locally?
2. **Concurrent access**: Multiple clients hitting same subreddit - any locking needed?

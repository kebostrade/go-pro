package graph

import (
	"context"
	"fmt"
	"sync"

	"basic/projects/graphql/pkg/models"
)

// Resolver is the root resolver
type Resolver struct {
	DB               *models.DB
	Loaders          *Loaders
	mu               sync.RWMutex
	userCreatedSubs  []chan *models.User
	postCreatedSubs  []chan *models.Post
	commentAddedSubs []commentSub
}

// commentSub holds a subscription for comments
type commentSub struct {
	postID string
	ch     chan *models.Comment
}

// NewResolver creates a new resolver
func NewResolver(db *models.DB) *Resolver {
	return &Resolver{
		DB:      db,
		Loaders: NewLoaders(db),
	}
}

// Query resolvers

func (r *Resolver) User(ctx context.Context, id string) (*models.User, error) {
	return r.DB.UserByID(id)
}

func (r *Resolver) Users(ctx context.Context, role *models.Role, active *bool, page *PageInput) ([]*models.User, error) {
	users := r.DB.GetAllUsers()

	// Filter by role
	if role != nil {
		filtered := make([]*models.User, 0)
		for _, u := range users {
			if u.Role == *role {
				filtered = append(filtered, u)
			}
		}
		users = filtered
	}

	// Filter by active
	if active != nil {
		filtered := make([]*models.User, 0)
		for _, u := range users {
			if u.Active == *active {
				filtered = append(filtered, u)
			}
		}
		users = filtered
	}

	// Pagination
	if page != nil {
		limit := page.Limit
		offset := page.Offset
		if limit <= 0 {
			limit = 10
		}
		if offset < 0 {
			offset = 0
		}
		end := offset + limit
		if end > len(users) {
			end = len(users)
		}
		if offset >= len(users) {
			users = []*models.User{}
		} else {
			users = users[offset:end]
		}
	}

	return users, nil
}

func (r *Resolver) Me(ctx context.Context) (*models.User, error) {
	userID := UserFromContext(ctx)
	if userID == "" {
		return nil, fmt.Errorf("unauthorized")
	}
	return r.DB.UserByID(userID)
}

func (r *Resolver) Post(ctx context.Context, id string) (*models.Post, error) {
	return r.DB.PostByID(id)
}

func (r *Resolver) Posts(ctx context.Context, authorID *string, published *bool, tags []string, search *string, page *PageInput) ([]*models.Post, error) {
	posts := r.DB.GetAllPosts()

	// Filter by author
	if authorID != nil {
		filtered := make([]*models.Post, 0)
		for _, p := range posts {
			if p.AuthorID == *authorID {
				filtered = append(filtered, p)
			}
		}
		posts = filtered
	}

	// Filter by published
	if published != nil {
		filtered := make([]*models.Post, 0)
		for _, p := range posts {
			if p.Published == *published {
				filtered = append(filtered, p)
			}
		}
		posts = filtered
	}

	// Filter by tags
	if len(tags) > 0 {
		filtered := make([]*models.Post, 0)
		for _, p := range posts {
			for _, tag := range tags {
				for _, pTag := range p.Tags {
					if tag == pTag {
						filtered = append(filtered, p)
						break
					}
				}
			}
		}
		posts = filtered
	}

	// Filter by search
	if search != nil && *search != "" {
		filtered := make([]*models.Post, 0)
		for _, p := range posts {
			if contains(p.Title, *search) || contains(p.Content, *search) {
				filtered = append(filtered, p)
			}
		}
		posts = filtered
	}

	// Pagination
	if page != nil {
		limit := page.Limit
		offset := page.Offset
		if limit <= 0 {
			limit = 10
		}
		if offset < 0 {
			offset = 0
		}
		end := offset + limit
		if end > len(posts) {
			end = len(posts)
		}
		if offset >= len(posts) {
			posts = []*models.Post{}
		} else {
			posts = posts[offset:end]
		}
	}

	return posts, nil
}

func (r *Resolver) GetComment(ctx context.Context, id string) (*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.DB.Comments {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("comment not found")
}

func (r *Resolver) Comments(ctx context.Context, postID string) ([]*models.Comment, error) {
	return r.DB.CommentsByPostID(postID), nil
}

func (r *Resolver) Stats(ctx context.Context) (*models.Stats, error) {
	return r.DB.Stats(), nil
}

// Mutation resolvers

func (r *Resolver) CreateUser(ctx context.Context, input CreateUserInput) (*models.User, error) {
	// Check if email already exists
	for _, u := range r.DB.Users {
		if u.Email == input.Email {
			return nil, fmt.Errorf("email already exists")
		}
	}

	role := models.RoleUser
	if input.Role != "" {
		role = input.Role
	}

	user := r.DB.CreateUser(input.Username, input.Email, input.Password, role)

	// Publish user created event
	r.publishUserCreated(user)

	return user, nil
}

func (r *Resolver) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*models.User, error) {
	return r.DB.UpdateUser(id, input.Username, input.Email, input.Active, input.Role)
}

func (r *Resolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	err := r.DB.DeleteUser(id)
	return err == nil, err
}

func (r *Resolver) DeactivateUser(ctx context.Context, id string) (bool, error) {
	active := false
	user, err := r.DB.UpdateUser(id, nil, nil, &active, nil)
	if user != nil {
		user.Active = false
	}
	return err == nil, err
}

func (r *Resolver) CreatePost(ctx context.Context, input CreatePostInput) (*models.Post, error) {
	userID := UserFromContext(ctx)
	if userID == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	post := r.DB.CreatePost(userID, input.Title, input.Content, input.Published, tags)

	// Publish post created event
	r.publishPostCreated(post)

	return post, nil
}

func (r *Resolver) UpdatePost(ctx context.Context, id string, input UpdatePostInput) (*models.Post, error) {
	return r.DB.UpdatePost(id, input.Title, input.Content, input.Published, input.Tags)
}

func (r *Resolver) DeletePost(ctx context.Context, id string) (bool, error) {
	err := r.DB.DeletePost(id)
	return err == nil, err
}

func (r *Resolver) PublishPost(ctx context.Context, id string) (*models.Post, error) {
	published := true
	return r.DB.UpdatePost(id, nil, nil, &published, nil)
}

func (r *Resolver) UnpublishPost(ctx context.Context, id string) (*models.Post, error) {
	published := false
	return r.DB.UpdatePost(id, nil, nil, &published, nil)
}

func (r *Resolver) CreateComment(ctx context.Context, input CreateCommentInput) (*models.Comment, error) {
	userID := UserFromContext(ctx)
	if userID == "" {
		return nil, fmt.Errorf("unauthorized")
	}

	comment := r.DB.CreateComment(userID, input.PostID, input.Content)

	// Publish comment added event
	r.publishCommentAdded(comment)

	return comment, nil
}

func (r *Resolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	err := r.DB.DeleteComment(id)
	return err == nil, err
}

// Field resolvers

func (r *Resolver) UserPosts(ctx context.Context, user *models.User, limit *int, offset *int) ([]*models.Post, error) {
	posts := r.DB.PostsByAuthorID(user.ID)

	l := 10
	o := 0
	if limit != nil {
		l = *limit
	}
	if offset != nil {
		o = *offset
	}

	end := o + l
	if end > len(posts) {
		end = len(posts)
	}
	if o >= len(posts) {
		return []*models.Post{}, nil
	}

	return posts[o:end], nil
}

func (r *Resolver) PostAuthor(ctx context.Context, post *models.Post) (*models.User, error) {
	return r.DB.UserByID(post.AuthorID)
}

func (r *Resolver) PostComments(ctx context.Context, post *models.Post) ([]*models.Comment, error) {
	return r.DB.CommentsByPostID(post.ID), nil
}

func (r *Resolver) CommentAuthor(ctx context.Context, comment *models.Comment) (*models.User, error) {
	return r.DB.UserByID(comment.AuthorID)
}

func (r *Resolver) CommentPost(ctx context.Context, comment *models.Comment) (*models.Post, error) {
	return r.DB.PostByID(comment.PostID)
}

// Subscription resolvers

func (r *Resolver) UserCreated(ctx context.Context) (<-chan *models.User, error) {
	ch := make(chan *models.User, 1)
	r.mu.Lock()
	r.userCreatedSubs = append(r.userCreatedSubs, ch)
	r.mu.Unlock()
	return ch, nil
}

func (r *Resolver) PostCreated(ctx context.Context) (<-chan *models.Post, error) {
	ch := make(chan *models.Post, 1)
	r.mu.Lock()
	r.postCreatedSubs = append(r.postCreatedSubs, ch)
	r.mu.Unlock()
	return ch, nil
}

func (r *Resolver) CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error) {
	ch := make(chan *models.Comment, 1)
	r.mu.Lock()
	r.commentAddedSubs = append(r.commentAddedSubs, commentSub{postID: postID, ch: ch})
	r.mu.Unlock()
	return ch, nil
}

// Subscription management
func (r *Resolver) publishUserCreated(user *models.User) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ch := range r.userCreatedSubs {
		select {
		case ch <- user:
		default:
		}
	}
}

func (r *Resolver) publishPostCreated(post *models.Post) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ch := range r.postCreatedSubs {
		select {
		case ch <- post:
		default:
		}
	}
}

func (r *Resolver) publishCommentAdded(comment *models.Comment) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, sub := range r.commentAddedSubs {
		if sub.postID == comment.PostID {
			select {
			case sub.ch <- comment:
			default:
			}
		}
	}
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Loaders

func NewLoaders(db *models.DB) *Loaders {
	return &Loaders{
		UserLoader: newUserLoader(db),
		PostLoader: newPostLoader(db),
	}
}

type Loaders struct {
	UserLoader *UserLoader
	PostLoader *PostLoader
}

type userResult struct {
	user *models.User
	err  error
}

type UserLoader struct {
	load func(keys []string) []userResult
}

func newUserLoader(db *models.DB) *UserLoader {
	return &UserLoader{
		load: func(keys []string) []userResult {
			results := make([]userResult, len(keys))
			for i, id := range keys {
				if u, err := db.UserByID(id); err != nil {
					results[i] = userResult{err: err}
				} else {
					results[i] = userResult{user: u}
				}
			}
			return results
		},
	}
}

func (l *UserLoader) Load(ctx context.Context, id string) (*models.User, error) {
	results := l.load([]string{id})
	return results[0].user, results[0].err
}

type postResult struct {
	posts []*models.Post
	err   error
}

type PostLoader struct {
	load func(authorID string) postResult
}

func newPostLoader(db *models.DB) *PostLoader {
	return &PostLoader{
		load: func(authorID string) postResult {
			return postResult{posts: db.PostsByAuthorID(authorID)}
		},
	}
}

func (l *PostLoader) LoadByAuthor(ctx context.Context, authorID string) []*models.Post {
	return l.load(authorID).posts
}

// Context key for user ID
type contextKey string

const userContextKey contextKey = "user_id"

func UserFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(userContextKey).(string); ok {
		return userID
	}
	return ""
}

func ContextWithUser(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userContextKey, userID)
}

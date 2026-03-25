package main

import (
	"context"
	"fmt"
	"time"
)

// ============================================
// Generated Models
// ============================================

// Role represents user role enum
type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
	RoleGuest Role = "GUEST"
)

// PostSort represents post sort enum
type PostSort string

const (
	SortCreatedAtAsc  PostSort = "CREATED_AT_ASC"
	SortCreatedAtDesc PostSort = "CREATED_AT_DESC"
	SortUpdatedAtAsc  PostSort = "UPDATED_AT_ASC"
	SortUpdatedAtDesc PostSort = "UPDATED_AT_DESC"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Post represents a blog post
type Post struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	AuthorID  string    `json:"authorId"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"authorId"`
	PostID    string    `json:"postId"`
	CreatedAt time.Time `json:"createdAt"`
}

// Stats represents system statistics
type Stats struct {
	Users          int `json:"users"`
	Posts          int `json:"posts"`
	Comments       int `json:"comments"`
	PublishedPosts int `json:"publishedPosts"`
}

// ============================================
// Input Types
// ============================================

// CreateUserInput represents input for creating a user
type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

// UpdateUserInput represents input for updating a user
type UpdateUserInput struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Active   *bool   `json:"active,omitempty"`
	Role     *Role   `json:"role,omitempty"`
}

// CreatePostInput represents input for creating a post
type CreatePostInput struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Published bool     `json:"published"`
	Tags      []string `json:"tags"`
}

// UpdatePostInput represents input for updating a post
type UpdatePostInput struct {
	Title     *string  `json:"title,omitempty"`
	Content   *string  `json:"content,omitempty"`
	Published *bool    `json:"published,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

// CreateCommentInput represents input for creating a comment
type CreateCommentInput struct {
	Content string `json:"content"`
	PostID  string `json:"postId"`
}

// PageInput represents pagination input
type PageInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ============================================
// Mock Database
// ============================================

// MockDB simulates a database
type MockDB struct {
	Users    map[string]*User
	Posts    map[string]*Post
	Comments map[string]*Comment
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	db := &MockDB{
		Users:    make(map[string]*User),
		Posts:    make(map[string]*Post),
		Comments: make(map[string]*Comment),
	}

	// Seed initial data
	now := time.Now()

	// Create users
	db.Users["1"] = &User{
		ID:        "1",
		Username:  "admin",
		Email:     "admin@example.com",
		Role:      RoleAdmin,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	db.Users["2"] = &User{
		ID:        "2",
		Username:  "alice",
		Email:     "alice@example.com",
		Role:      RoleUser,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	db.Users["3"] = &User{
		ID:        "3",
		Username:  "bob",
		Email:     "bob@example.com",
		Role:      RoleUser,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create posts
	db.Posts["1"] = &Post{
		ID:        "1",
		Title:     "First Post",
		Content:   "This is the first post",
		Published: true,
		AuthorID:  "2",
		Tags:      []string{"welcome", "intro"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	db.Posts["2"] = &Post{
		ID:        "2",
		Title:     "GraphQL Tutorial",
		Content:   "Learn GraphQL with Go",
		Published: true,
		AuthorID:  "2",
		Tags:      []string{"graphql", "go", "tutorial"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	db.Posts["3"] = &Post{
		ID:        "3",
		Title:     "Draft Post",
		Content:   "This is a draft",
		Published: false,
		AuthorID:  "3",
		Tags:      []string{"draft"},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Create comments
	db.Comments["1"] = &Comment{
		ID:        "1",
		Content:   "Great post!",
		AuthorID:  "3",
		PostID:    "1",
		CreatedAt: now,
	}

	return db
}

// ============================================
// Database Methods
// ============================================

// UserByID retrieves a user by ID
func (db *MockDB) UserByID(id string) (*User, error) {
	user, ok := db.Users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetAllUsers retrieves all users
func (db *MockDB) GetAllUsers() ([]*User, error) {
	users := make([]*User, 0, len(db.Users))
	for _, user := range db.Users {
		users = append(users, user)
	}
	return users, nil
}

// CreateUser creates a new user
func (db *MockDB) CreateUser(input CreateUserInput) (*User, error) {
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	user := &User{
		ID:        id,
		Username:  input.Username,
		Email:     input.Email,
		Role:      input.Role,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Users[id] = user
	return user, nil
}

// UpdateUser updates an existing user
func (db *MockDB) UpdateUser(id string, input UpdateUserInput) (*User, error) {
	user, err := db.UserByID(id)
	if err != nil {
		return nil, err
	}

	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Active != nil {
		user.Active = *input.Active
	}
	if input.Role != nil {
		user.Role = *input.Role
	}

	user.UpdatedAt = time.Now()
	return user, nil
}

// DeleteUser deletes a user
func (db *MockDB) DeleteUser(id string) error {
	if _, ok := db.Users[id]; !ok {
		return fmt.Errorf("user not found")
	}
	delete(db.Users, id)
	return nil
}

// PostByID retrieves a post by ID
func (db *MockDB) PostByID(id string) (*Post, error) {
	post, ok := db.Posts[id]
	if !ok {
		return nil, fmt.Errorf("post not found")
	}
	return post, nil
}

// GetAllPosts retrieves all posts
func (db *MockDB) GetAllPosts() ([]*Post, error) {
	posts := make([]*Post, 0, len(db.Posts))
	for _, post := range db.Posts {
		posts = append(posts, post)
	}
	return posts, nil
}

// CreatePost creates a new post
func (db *MockDB) CreatePost(authorID string, input CreatePostInput) (*Post, error) {
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	post := &Post{
		ID:        id,
		Title:     input.Title,
		Content:   input.Content,
		Published: input.Published,
		AuthorID:  authorID,
		Tags:      input.Tags,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	db.Posts[id] = post
	return post, nil
}

// UpdatePost updates an existing post
func (db *MockDB) UpdatePost(id string, input UpdatePostInput) (*Post, error) {
	post, err := db.PostByID(id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		post.Title = *input.Title
	}
	if input.Content != nil {
		post.Content = *input.Content
	}
	if input.Published != nil {
		post.Published = *input.Published
	}
	if input.Tags != nil {
		post.Tags = input.Tags
	}

	post.UpdatedAt = time.Now()
	return post, nil
}

// DeletePost deletes a post
func (db *MockDB) DeletePost(id string) error {
	if _, ok := db.Posts[id]; !ok {
		return fmt.Errorf("post not found")
	}
	delete(db.Posts, id)
	return nil
}

// CommentsByPostID retrieves comments for a post
func (db *MockDB) CommentsByPostID(postID string) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	for _, comment := range db.Comments {
		if comment.PostID == postID {
			comments = append(comments, comment)
		}
	}
	return comments, nil
}

// CreateComment creates a new comment
func (db *MockDB) CreateComment(authorID string, input CreateCommentInput) (*Comment, error) {
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	comment := &Comment{
		ID:        id,
		Content:   input.Content,
		AuthorID:  authorID,
		PostID:    input.PostID,
		CreatedAt: time.Now(),
	}
	db.Comments[id] = comment
	return comment, nil
}

// Stats returns system statistics
func (db *MockDB) Stats() *Stats {
	publishedCount := 0
	for _, post := range db.Posts {
		if post.Published {
			publishedCount++
		}
	}

	return &Stats{
		Users:          len(db.Users),
		Posts:          len(db.Posts),
		Comments:       len(db.Comments),
		PublishedPosts: publishedCount,
	}
}

// ============================================
// Context Keys
// ============================================

type contextKey string

const (
	userContextKey contextKey = "user"
)

// UserFromContext retrieves user from context
func UserFromContext(ctx context.Context) *User {
	user, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		return nil
	}
	return user
}

// ContextWithUser adds user to context
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

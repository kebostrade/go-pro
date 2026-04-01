package models

import (
	"fmt"
	"sync"
	"time"
)

// Role represents user role enum
type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
	RoleGuest Role = "GUEST"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
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

// DB represents an in-memory database
type DB struct {
	Users    map[string]*User
	Posts    map[string]*Post
	Comments map[string]*Comment
	mu       sync.RWMutex
}

// NewDB creates a new mock database with seed data
func NewDB() *DB {
	db := &DB{
		Users:    make(map[string]*User),
		Posts:    make(map[string]*Post),
		Comments: make(map[string]*Comment),
	}

	now := time.Now()

	db.Users["1"] = &User{ID: "1", Username: "admin", Email: "admin@example.com", Password: "admin123", Role: RoleAdmin, Active: true, CreatedAt: now, UpdatedAt: now}
	db.Users["2"] = &User{ID: "2", Username: "alice", Email: "alice@example.com", Password: "alice123", Role: RoleUser, Active: true, CreatedAt: now, UpdatedAt: now}
	db.Users["3"] = &User{ID: "3", Username: "bob", Email: "bob@example.com", Password: "bob123", Role: RoleUser, Active: true, CreatedAt: now, UpdatedAt: now}

	db.Posts["1"] = &Post{ID: "1", Title: "First Post", Content: "This is the first post", Published: true, AuthorID: "2", Tags: []string{"welcome", "intro"}, CreatedAt: now, UpdatedAt: now}
	db.Posts["2"] = &Post{ID: "2", Title: "GraphQL Tutorial", Content: "Learn GraphQL with Go", Published: true, AuthorID: "2", Tags: []string{"graphql", "go", "tutorial"}, CreatedAt: now, UpdatedAt: now}
	db.Posts["3"] = &Post{ID: "3", Title: "Draft Post", Content: "This is a draft", Published: false, AuthorID: "3", Tags: []string{"draft"}, CreatedAt: now, UpdatedAt: now}

	db.Comments["1"] = &Comment{ID: "1", Content: "Great post!", AuthorID: "3", PostID: "1", CreatedAt: now}

	return db
}

func (db *DB) UserByID(id string) (*User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if u, ok := db.Users[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (db *DB) GetAllUsers() []*User {
	db.mu.RLock()
	defer db.mu.RUnlock()
	users := make([]*User, 0, len(db.Users))
	for _, u := range db.Users {
		users = append(users, u)
	}
	return users
}

func (db *DB) CreateUser(username, email, password string, role Role) *User {
	db.mu.Lock()
	defer db.mu.Unlock()
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	u := &User{ID: id, Username: username, Email: email, Password: password, Role: role, Active: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	db.Users[id] = u
	return u
}

func (db *DB) UpdateUser(id string, username, email *string, active *bool, role *Role) (*User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	u, ok := db.Users[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	if username != nil {
		u.Username = *username
	}
	if email != nil {
		u.Email = *email
	}
	if active != nil {
		u.Active = *active
	}
	if role != nil {
		u.Role = *role
	}
	u.UpdatedAt = time.Now()
	return u, nil
}

func (db *DB) DeleteUser(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Users[id]; !ok {
		return fmt.Errorf("user not found")
	}
	delete(db.Users, id)
	return nil
}

func (db *DB) PostByID(id string) (*Post, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	if p, ok := db.Posts[id]; ok {
		return p, nil
	}
	return nil, fmt.Errorf("post not found")
}

func (db *DB) GetAllPosts() []*Post {
	db.mu.RLock()
	defer db.mu.RUnlock()
	posts := make([]*Post, 0, len(db.Posts))
	for _, p := range db.Posts {
		posts = append(posts, p)
	}
	return posts
}

func (db *DB) CreatePost(authorID, title, content string, published bool, tags []string) *Post {
	db.mu.Lock()
	defer db.mu.Unlock()
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	if tags == nil {
		tags = []string{}
	}
	p := &Post{ID: id, Title: title, Content: content, Published: published, AuthorID: authorID, Tags: tags, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	db.Posts[id] = p
	return p
}

func (db *DB) UpdatePost(id string, title, content *string, published *bool, tags []string) (*Post, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	p, ok := db.Posts[id]
	if !ok {
		return nil, fmt.Errorf("post not found")
	}
	if title != nil {
		p.Title = *title
	}
	if content != nil {
		p.Content = *content
	}
	if published != nil {
		p.Published = *published
	}
	if tags != nil {
		p.Tags = tags
	}
	p.UpdatedAt = time.Now()
	return p, nil
}

func (db *DB) DeletePost(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Posts[id]; !ok {
		return fmt.Errorf("post not found")
	}
	delete(db.Posts, id)
	return nil
}

func (db *DB) PostsByAuthorID(authorID string) []*Post {
	db.mu.RLock()
	defer db.mu.RUnlock()
	posts := make([]*Post, 0)
	for _, p := range db.Posts {
		if p.AuthorID == authorID {
			posts = append(posts, p)
		}
	}
	return posts
}

func (db *DB) CommentsByPostID(postID string) []*Comment {
	db.mu.RLock()
	defer db.mu.RUnlock()
	comments := make([]*Comment, 0)
	for _, c := range db.Comments {
		if c.PostID == postID {
			comments = append(comments, c)
		}
	}
	return comments
}

func (db *DB) CreateComment(authorID, postID, content string) *Comment {
	db.mu.Lock()
	defer db.mu.Unlock()
	id := fmt.Sprintf("%d", time.Now().UnixNano())
	c := &Comment{ID: id, Content: content, AuthorID: authorID, PostID: postID, CreatedAt: time.Now()}
	db.Comments[id] = c
	return c
}

func (db *DB) DeleteComment(id string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, ok := db.Comments[id]; !ok {
		return fmt.Errorf("comment not found")
	}
	delete(db.Comments, id)
	return nil
}

func (db *DB) Stats() *Stats {
	db.mu.RLock()
	defer db.mu.RUnlock()
	published := 0
	for _, p := range db.Posts {
		if p.Published {
			published++
		}
	}
	return &Stats{Users: len(db.Users), Posts: len(db.Posts), Comments: len(db.Comments), PublishedPosts: published}
}

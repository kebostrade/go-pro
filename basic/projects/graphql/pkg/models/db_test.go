package models_test

import (
	"testing"
	"time"

	"basic/projects/graphql/pkg/models"
)

func TestNewDB(t *testing.T) {
	db := models.NewDB()
	if db == nil {
		t.Fatal("NewDB() returned nil")
	}
	if len(db.Users) != 3 {
		t.Errorf("expected 3 users, got %d", len(db.Users))
	}
	if len(db.Posts) != 3 {
		t.Errorf("expected 3 posts, got %d", len(db.Posts))
	}
	if len(db.Comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(db.Comments))
	}
}

func TestDBUserByID(t *testing.T) {
	db := models.NewDB()

	user, err := db.UserByID("1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.Username != "admin" {
		t.Errorf("expected username 'admin', got '%s'", user.Username)
	}

	_, err = db.UserByID("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestDBGetAllUsers(t *testing.T) {
	db := models.NewDB()
	users := db.GetAllUsers()
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestDBCreateUser(t *testing.T) {
	db := models.NewDB()
	initialCount := len(db.Users)

	user := db.CreateUser("newuser", "new@example.com", "password123", models.RoleUser)
	if user.Username != "newuser" {
		t.Errorf("expected username 'newuser', got '%s'", user.Username)
	}
	if len(db.Users) != initialCount+1 {
		t.Errorf("expected %d users after create, got %d", initialCount+1, len(db.Users))
	}
}

func TestDBUpdateUser(t *testing.T) {
	db := models.NewDB()

	newName := "updatedadmin"
	user, err := db.UpdateUser("1", &newName, nil, nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if user.Username != "updatedadmin" {
		t.Errorf("expected username 'updatedadmin', got '%s'", user.Username)
	}
}

func TestDBDeleteUser(t *testing.T) {
	db := models.NewDB()
	initialCount := len(db.Users)

	err := db.DeleteUser("1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(db.Users) != initialCount-1 {
		t.Errorf("expected %d users after delete, got %d", initialCount-1, len(db.Users))
	}

	err = db.DeleteUser("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent user")
	}
}

func TestDBPostByID(t *testing.T) {
	db := models.NewDB()

	post, err := db.PostByID("1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if post.Title != "First Post" {
		t.Errorf("expected title 'First Post', got '%s'", post.Title)
	}

	_, err = db.PostByID("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent post")
	}
}

func TestDBGetAllPosts(t *testing.T) {
	db := models.NewDB()
	posts := db.GetAllPosts()
	if len(posts) != 3 {
		t.Errorf("expected 3 posts, got %d", len(posts))
	}
}

func TestDBCreatePost(t *testing.T) {
	db := models.NewDB()
	initialCount := len(db.Posts)

	post := db.CreatePost("2", "New Post", "Content here", true, []string{"test"})
	if post.Title != "New Post" {
		t.Errorf("expected title 'New Post', got '%s'", post.Title)
	}
	if len(db.Posts) != initialCount+1 {
		t.Errorf("expected %d posts after create, got %d", initialCount+1, len(db.Posts))
	}
}

func TestDBUpdatePost(t *testing.T) {
	db := models.NewDB()

	newTitle := "Updated Title"
	post, err := db.UpdatePost("1", &newTitle, nil, nil, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if post.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got '%s'", post.Title)
	}
}

func TestDBDeletePost(t *testing.T) {
	db := models.NewDB()
	initialCount := len(db.Posts)

	err := db.DeletePost("1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(db.Posts) != initialCount-1 {
		t.Errorf("expected %d posts after delete, got %d", initialCount-1, len(db.Posts))
	}
}

func TestDBPostsByAuthorID(t *testing.T) {
	db := models.NewDB()
	posts := db.PostsByAuthorID("2")
	if len(posts) != 2 {
		t.Errorf("expected 2 posts by author 2, got %d", len(posts))
	}
}

func TestDBCommentsByPostID(t *testing.T) {
	db := models.NewDB()
	comments := db.CommentsByPostID("1")
	if len(comments) != 1 {
		t.Errorf("expected 1 comment for post 1, got %d", len(comments))
	}
}

func TestDBCreateComment(t *testing.T) {
	db := models.NewDB()
	initialCount := len(db.Comments)

	comment := db.CreateComment("3", "1", "New comment")
	if comment.Content != "New comment" {
		t.Errorf("expected content 'New comment', got '%s'", comment.Content)
	}
	if len(db.Comments) != initialCount+1 {
		t.Errorf("expected %d comments after create, got %d", initialCount+1, len(db.Comments))
	}
}

func TestDBDeleteComment(t *testing.T) {
	db := models.NewDB()
	initialCount := len(db.Comments)

	err := db.DeleteComment("1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(db.Comments) != initialCount-1 {
		t.Errorf("expected %d comments after delete, got %d", initialCount-1, len(db.Comments))
	}
}

func TestDBStats(t *testing.T) {
	db := models.NewDB()
	stats := db.Stats()
	if stats.Users != 3 {
		t.Errorf("expected 3 users, got %d", stats.Users)
	}
	if stats.Posts != 3 {
		t.Errorf("expected 3 posts, got %d", stats.Posts)
	}
	if stats.Comments != 1 {
		t.Errorf("expected 1 comment, got %d", stats.Comments)
	}
	if stats.PublishedPosts != 2 {
		t.Errorf("expected 2 published posts, got %d", stats.PublishedPosts)
	}
}

func TestUserTimestamps(t *testing.T) {
	db := models.NewDB()
	user, _ := db.UserByID("1")
	now := time.Now()
	if user.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
	if user.CreatedAt.After(now) {
		t.Error("CreatedAt should not be in the future")
	}
}

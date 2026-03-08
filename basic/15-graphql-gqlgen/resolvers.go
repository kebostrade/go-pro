package main

import (
	"context"
	"fmt"
	"strings"
)

// ============================================
// Resolver Interfaces
// ============================================
// These interfaces define the resolver methods implemented below

type QueryResolver interface {
	User(ctx context.Context, id string) (*User, error)
	Users(ctx context.Context, role *Role, active *bool, page *PageInput, sort *PostSort) ([]*User, error)
	Me(ctx context.Context) (*User, error)
	Post(ctx context.Context, id string) (*Post, error)
	Posts(ctx context.Context, authorId *string, published *bool, tags []string, search *string, page *PageInput, sort *PostSort) ([]*Post, error)
	Comment(ctx context.Context, id string) (*Comment, error)
	Comments(ctx context.Context, postId string) ([]*Comment, error)
	SearchUsers(ctx context.Context, query string, limit *int) ([]*User, error)
	SearchPosts(ctx context.Context, query string, limit *int) ([]*Post, error)
	Stats(ctx context.Context) (*Stats, error)
}

type MutationResolver interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*User, error)
	UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	DeactivateUser(ctx context.Context, id string) (bool, error)
	CreatePost(ctx context.Context, input CreatePostInput) (*Post, error)
	UpdatePost(ctx context.Context, id string, input UpdatePostInput) (*Post, error)
	DeletePost(ctx context.Context, id string) (bool, error)
	PublishPost(ctx context.Context, id string) (*Post, error)
	UnpublishPost(ctx context.Context, id string) (*Post, error)
	CreateComment(ctx context.Context, input CreateCommentInput) (*Comment, error)
	DeleteComment(ctx context.Context, id string) (bool, error)
	UploadImage(ctx context.Context, file interface{}) (string, error)
}

type UserResolver interface {
	Posts(ctx context.Context, obj *User, limit *int, offset *int) ([]*Post, error)
}

type PostResolver interface {
	Author(ctx context.Context, obj *Post) (*User, error)
	Comments(ctx context.Context, obj *Post) ([]*Comment, error)
}

type CommentResolver interface {
	Author(ctx context.Context, obj *Comment) (*User, error)
	Post(ctx context.Context, obj *Comment) (*Post, error)
}

// ============================================
// Resolver Structure
// ============================================

type resolvers struct {
	db *MockDB
}

// ============================================
// Query Resolvers
// ============================================

func (r *queryResolver) User(ctx context.Context, id string) (*User, error) {
	// Fetch user from database
	user, err := r.db.UserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	return user, nil
}

func (r *queryResolver) Users(ctx context.Context, role *Role, active *bool, page *PageInput, sort *PostSort) ([]*User, error) {
	// Get all users
	users, err := r.db.GetAllUsers()
	if err != nil {
		return nil, err
	}

	// Filter by role
	if role != nil {
		filtered := make([]*User, 0)
		for _, user := range users {
			if user.Role == *role {
				filtered = append(filtered, user)
			}
		}
		users = filtered
	}

	// Filter by active status
	if active != nil {
		filtered := make([]*User, 0)
		for _, user := range users {
			if user.Active == *active {
				filtered = append(filtered, user)
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

		if offset >= len(users) {
			return []*User{}, nil
		}

		end := offset + limit
		if end > len(users) {
			end = len(users)
		}

		users = users[offset:end]
	}

	return users, nil
}

func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	// Get user from context (authenticated)
	user := UserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("not authenticated")
	}

	return user, nil
}

func (r *queryResolver) Post(ctx context.Context, id string) (*Post, error) {
	post, err := r.db.PostByID(id)
	if err != nil {
		return nil, fmt.Errorf("post not found: %s", id)
	}

	return post, nil
}

func (r *queryResolver) Posts(ctx context.Context, authorId *string, published *bool, tags []string, search *string, page *PageInput, sort *PostSort) ([]*Post, error) {
	// Get all posts
	posts, err := r.db.GetAllPosts()
	if err != nil {
		return nil, err
	}

	// Filter by author
	if authorId != nil {
		filtered := make([]*Post, 0)
		for _, post := range posts {
			if post.AuthorID == *authorId {
				filtered = append(filtered, post)
			}
		}
		posts = filtered
	}

	// Filter by published status
	if published != nil {
		filtered := make([]*Post, 0)
		for _, post := range posts {
			if post.Published == *published {
				filtered = append(filtered, post)
			}
		}
		posts = filtered
	}

	// Filter by tags
	if len(tags) > 0 {
		filtered := make([]*Post, 0)
		for _, post := range posts {
			if containsAllTags(post.Tags, tags) {
				filtered = append(filtered, post)
			}
		}
		posts = filtered
	}

	// Search by title/content
	if search != nil {
		filtered := make([]*Post, 0)
		searchLower := strings.ToLower(*search)
		for _, post := range posts {
			if strings.Contains(strings.ToLower(post.Title), searchLower) ||
				strings.Contains(strings.ToLower(post.Content), searchLower) {
				filtered = append(filtered, post)
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

		if offset >= len(posts) {
			return []*Post{}, nil
		}

		end := offset + limit
		if end > len(posts) {
			end = len(posts)
		}

		posts = posts[offset:end]
	}

	return posts, nil
}

func (r *queryResolver) Comment(ctx context.Context, id string) (*Comment, error) {
	comment, ok := r.db.Comments[id]
	if !ok {
		return nil, fmt.Errorf("comment not found: %s", id)
	}

	return comment, nil
}

func (r *queryResolver) Comments(ctx context.Context, postId string) ([]*Comment, error) {
	comments, err := r.db.CommentsByPostID(postId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *queryResolver) SearchUsers(ctx context.Context, query string, limit *int) ([]*User, error) {
	users, err := r.db.GetAllUsers()
	if err != nil {
		return nil, err
	}

	queryLower := strings.ToLower(query)
	results := make([]*User, 0)

	for _, user := range users {
		if strings.Contains(strings.ToLower(user.Username), queryLower) ||
			strings.Contains(strings.ToLower(user.Email), queryLower) {
			results = append(results, user)
		}
	}

	if limit != nil && *limit > 0 && len(results) > *limit {
		results = results[:*limit]
	}

	return results, nil
}

func (r *queryResolver) SearchPosts(ctx context.Context, query string, limit *int) ([]*Post, error) {
	posts, err := r.db.GetAllPosts()
	if err != nil {
		return nil, err
	}

	queryLower := strings.ToLower(query)
	results := make([]*Post, 0)

	for _, post := range posts {
		if strings.Contains(strings.ToLower(post.Title), queryLower) ||
			strings.Contains(strings.ToLower(post.Content), queryLower) {
			results = append(results, post)
		}
	}

	if limit != nil && *limit > 0 && len(results) > *limit {
		results = results[:*limit]
	}

	return results, nil
}

func (r *queryResolver) Stats(ctx context.Context) (*Stats, error) {
	return r.db.Stats(), nil
}

// ============================================
// Mutation Resolvers
// ============================================

func (r *mutationResolver) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	// Get authenticated user from context
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Check if user is admin
	if authUser.Role != RoleAdmin {
		return nil, fmt.Errorf("forbidden: only admins can create users")
	}

	// Validate input
	if input.Username == "" || input.Email == "" {
		return nil, fmt.Errorf("username and email are required")
	}

	// Create user
	user, err := r.db.CreateUser(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*User, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Check permission (admin or own account)
	if authUser.Role != RoleAdmin && authUser.ID != id {
		return nil, fmt.Errorf("forbidden")
	}

	// Update user
	user, err := r.db.UpdateUser(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (bool, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return false, fmt.Errorf("unauthorized")
	}

	// Only admins can delete users
	if authUser.Role != RoleAdmin {
		return false, fmt.Errorf("forbidden: only admins can delete users")
	}

	// Delete user
	err := r.db.DeleteUser(id)
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) DeactivateUser(ctx context.Context, id string) (bool, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return false, fmt.Errorf("unauthorized")
	}

	// Only admins can deactivate users
	if authUser.Role != RoleAdmin {
		return false, fmt.Errorf("forbidden")
	}

	// Deactivate user
	active := false
	_, err := r.db.UpdateUser(id, UpdateUserInput{Active: &active})
	if err != nil {
		return false, fmt.Errorf("failed to deactivate user: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) CreatePost(ctx context.Context, input CreatePostInput) (*Post, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Validate input
	if input.Title == "" || input.Content == "" {
		return nil, fmt.Errorf("title and content are required")
	}

	// Create post
	post, err := r.db.CreatePost(authUser.ID, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (r *mutationResolver) UpdatePost(ctx context.Context, id string, input UpdatePostInput) (*Post, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get post
	post, err := r.db.PostByID(id)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	// Check permission (admin or author)
	if authUser.Role != RoleAdmin && post.AuthorID != authUser.ID {
		return nil, fmt.Errorf("forbidden")
	}

	// Update post
	post, err = r.db.UpdatePost(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update post: %w", err)
	}

	return post, nil
}

func (r *mutationResolver) DeletePost(ctx context.Context, id string) (bool, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return false, fmt.Errorf("unauthorized")
	}

	// Get post
	post, err := r.db.PostByID(id)
	if err != nil {
		return false, fmt.Errorf("post not found")
	}

	// Check permission (admin or author)
	if authUser.Role != RoleAdmin && post.AuthorID != authUser.ID {
		return false, fmt.Errorf("forbidden")
	}

	// Delete post
	err = r.db.DeletePost(id)
	if err != nil {
		return false, fmt.Errorf("failed to delete post: %w", err)
	}

	return true, nil
}

func (r *mutationResolver) PublishPost(ctx context.Context, id string) (*Post, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get post
	post, err := r.db.PostByID(id)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	// Check permission
	if authUser.Role != RoleAdmin && post.AuthorID != authUser.ID {
		return nil, fmt.Errorf("forbidden")
	}

	// Publish post
	published := true
	post, err = r.db.UpdatePost(id, UpdatePostInput{Published: &published})
	if err != nil {
		return nil, fmt.Errorf("failed to publish post: %w", err)
	}

	return post, nil
}

func (r *mutationResolver) UnpublishPost(ctx context.Context, id string) (*Post, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get post
	post, err := r.db.PostByID(id)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	// Check permission
	if authUser.Role != RoleAdmin && post.AuthorID != authUser.ID {
		return nil, fmt.Errorf("forbidden")
	}

	// Unpublish post
	published := false
	post, err = r.db.UpdatePost(id, UpdatePostInput{Published: &published})
	if err != nil {
		return nil, fmt.Errorf("failed to unpublish post: %w", err)
	}

	return post, nil
}

func (r *mutationResolver) CreateComment(ctx context.Context, input CreateCommentInput) (*Comment, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	// Validate input
	if input.Content == "" || input.PostID == "" {
		return nil, fmt.Errorf("content and postId are required")
	}

	// Check if post exists
	_, err := r.db.PostByID(input.PostID)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}

	// Create comment
	comment, err := r.db.CreateComment(authUser.ID, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	return comment, nil
}

func (r *mutationResolver) DeleteComment(ctx context.Context, id string) (bool, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return false, fmt.Errorf("unauthorized")
	}

	// Get comment
	comment, ok := r.db.Comments[id]
	if !ok {
		return false, fmt.Errorf("comment not found")
	}

	// Check permission (admin or author)
	if authUser.Role != RoleAdmin && comment.AuthorID != authUser.ID {
		return false, fmt.Errorf("forbidden")
	}

	// Delete comment
	delete(r.db.Comments, id)
	return true, nil
}

func (r *mutationResolver) UploadImage(ctx context.Context, file interface{}) (string, error) {
	// Get authenticated user
	authUser := UserFromContext(ctx)
	if authUser == nil {
		return "", fmt.Errorf("unauthorized")
	}

	// In production, save to S3 or cloud storage
	// For demo, return a fake URL
	url := fmt.Sprintf("https://example.com/images/%d.jpg", time.Now().UnixNano())

	return url, nil
}

// ============================================
// Field Resolvers
// ============================================

func (r *userResolver) Posts(ctx context.Context, obj *User, limit *int, offset *int) ([]*Post, error) {
	// Get all posts
	posts, err := r.db.GetAllPosts()
	if err != nil {
		return nil, err
	}

	// Filter by author
	filtered := make([]*Post, 0)
	for _, post := range posts {
		if post.AuthorID == obj.ID {
			filtered = append(filtered, post)
		}
	}

	// Pagination
	if limit != nil {
		if offset == nil {
			offset = new(int)
		}

		if *offset >= len(filtered) {
			return []*Post{}, nil
		}

		end := *offset + *limit
		if end > len(filtered) {
			end = len(filtered)
		}

		filtered = filtered[*offset:end]
	}

	return filtered, nil
}

func (r *postResolver) Author(ctx context.Context, obj *Post) (*User, error) {
	return r.db.UserByID(obj.AuthorID)
}

func (r *postResolver) Comments(ctx context.Context, obj *Post) ([]*Comment, error) {
	return r.db.CommentsByPostID(obj.ID)
}

func (r *commentResolver) Author(ctx context.Context, obj *Comment) (*User, error) {
	return r.db.UserByID(obj.AuthorID)
}

func (r *commentResolver) Post(ctx context.Context, obj *Comment) (*Post, error) {
	return r.db.PostByID(obj.PostID)
}

// ============================================
// Helper Functions
// ============================================

func containsAllTags(postTags []string, requiredTags []string) bool {
	tagMap := make(map[string]bool)
	for _, tag := range postTags {
		tagMap[tag] = true
	}

	for _, required := range requiredTags {
		if !tagMap[required] {
			return false
		}
	}

	return true
}

// ============================================
// Resolver Interface Implementation
// ============================================

func (r *resolvers) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *resolvers) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *resolvers) User() UserResolver {
	return &userResolver{r}
}

func (r *resolvers) Post() PostResolver {
	return &postResolver{r}
}

func (r *resolvers) Comment() CommentResolver {
	return &commentResolver{r}
}

type queryResolver struct{ *resolvers }
type mutationResolver struct{ *resolvers }
type userResolver struct{ *resolvers }
type postResolver struct{ *resolvers }
type commentResolver struct{ *resolvers }

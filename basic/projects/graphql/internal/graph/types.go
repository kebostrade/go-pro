package graph

import "basic/projects/graphql/pkg/models"

// Role represents the user role enum
type Role = models.Role

const (
	RoleAdmin Role = models.RoleAdmin
	RoleUser  Role = models.RoleUser
	RoleGuest Role = models.RoleGuest
)

// PageInput for pagination
type PageInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// PostSort for sorting posts
type PostSort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// CreateUserInput for creating a user
type CreateUserInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

// UpdateUserInput for updating a user
type UpdateUserInput struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Active   *bool   `json:"active"`
	Role     *Role   `json:"role"`
}

// CreatePostInput for creating a post
type CreatePostInput struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Published bool     `json:"published"`
	Tags      []string `json:"tags"`
}

// UpdatePostInput for updating a post
type UpdatePostInput struct {
	Title     *string  `json:"title"`
	Content   *string  `json:"content"`
	Published *bool    `json:"published"`
	Tags      []string `json:"tags"`
}

// CreateCommentInput for creating a comment
type CreateCommentInput struct {
	Content string `json:"content"`
	PostID  string `json:"postId"`
}

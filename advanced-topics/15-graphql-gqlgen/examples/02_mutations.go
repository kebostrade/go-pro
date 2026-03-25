package main

import (
	"context"
	"fmt"
)

// ============================================
// GraphQL Mutation Examples
// ============================================
// This file demonstrates various mutation patterns in GraphQL

// Example 1: Create User
func demoCreateUser(ctx context.Context) {
	fmt.Println("\n=== Demo 1: Create User ===")

	mutation := `
		mutation CreateUser($input: CreateUserInput!) {
			createUser(input: $input) {
				id
				username
				email
				role
				active
				createdAt
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"username": "newuser",
			"email":    "newuser@example.com",
			"password": "securepassword123",
			"role":     "USER",
		},
	}

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: Newly created user with generated ID")
	fmt.Println("Use case: User registration, admin user creation")
}

// Example 2: Update User
func demoUpdateUser(ctx context.Context) {
	fmt.Println("\n=== Demo 2: Update User ===")

	mutation := `
		mutation UpdateUser($id: ID!, $input: UpdateUserInput!) {
			updateUser(id: $id, input: $input) {
				id
				username
				email
				active
				updatedAt
			}
		}
	`

	variables := map[string]interface{}{
		"id": "2",
		"input": map[string]interface{}{
			"username": stringPtr("alice_updated"),
			"active":   boolPtr(true),
		},
	}

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: Updated user with new values")
	fmt.Println("Use case: Profile updates, account management")
}

// Example 3: Delete User
func demoDeleteUser(ctx context.Context) {
	fmt.Println("\n=== Demo 3: Delete User ===")

	mutation := `
		mutation DeleteUser($id: ID!) {
			deleteUser(id: $id)
		}
	`

	variables := map[string]interface{}{
		"id": "3",
	}

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: true if successful")
	fmt.Println("Use case: Account deletion, user management")
}

// Example 4: Create Post
func demoCreatePost(ctx context.Context) {
	fmt.Println("\n=== Demo 4: Create Post ===")

	mutation := `
		mutation CreatePost($input: CreatePostInput!) {
			createPost(input: $input) {
				id
				title
				content
				published
				tags
				author {
					username
				}
				createdAt
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"title":     "My First GraphQL Post",
			"content":   "Learning GraphQL with Go is amazing!",
			"published": true,
			"tags":      []string{"graphql", "go", "tutorial"},
		},
	}

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: Newly created post linked to authenticated user")
	fmt.Println("Use case: Content creation, blogging")
}

// Example 5: Publish Post
func demoPublishPost(ctx context.Context) {
	fmt.Println("\n=== Demo 5: Publish Post ===")

	mutation := `
		mutation PublishPost($id: ID!) {
			publishPost(id: $id) {
				id
				title
				published
				updatedAt
			}
		}
	`

	variables := map[string]interface{}{
		"id": "3",
	}

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: Post with published=true")
	fmt.Println("Use case: Content publishing workflow")
}

// Example 6: Create Comment
func demoCreateComment(ctx context.Context) {
	fmt.Println("\n=== Demo 6: Create Comment ===")

	mutation := `
		mutation CreateComment($input: CreateCommentInput!) {
			createComment(input: $input) {
				id
				content
				author {
					username
				}
				post {
					title
				}
				createdAt
			}
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"content": "Great post! Thanks for sharing.",
			"postId":  "1",
		},
	}

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: Newly created comment")
	fmt.Println("Use case: User engagement, discussions")
}

// Example 7: Batch Mutations
func demoBatchMutations(ctx context.Context) {
	fmt.Println("\n=== Demo 7: Batch Mutations ===")

	mutation := `
		mutation {
			# Create user
			createUser(input: {
				username: "charlie"
				email: "charlie@example.com"
				password: "password123"
				role: USER
			}) {
				id
				username
			}

			# Create post
			createPost(input: {
				title: "Hello World"
				content: "My first post"
				published: true
				tags: ["hello"]
			}) {
				id
				title
			}
		}
	`

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Println("\nResult: Multiple mutations in single request")
	fmt.Println("Use case: Complex operations, reduce round trips")
	fmt.Println("Note: Mutations execute sequentially, not parallel")
}

// Example 8: Mutation with Error Handling
func demoMutationWithError(ctx context.Context) {
	fmt.Println("\n=== Demo 8: Mutation with Error ===")

	mutation := `
		mutation CreateUser {
			createUser(input: {
				username: ""
				email: "invalid"
				password: "123"
				role: USER
			}) {
				id
				username
			}
		}
	`

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Println("\nExpected Error: Validation errors")
	fmt.Println("Result:")
	fmt.Println(`{
  "data": { "createUser": null },
  "errors": [
    {
      "message": "username and email are required",
      "path": ["createUser"]
    }
  ]
}`)
	fmt.Println("\nUse case: Input validation, error handling")
}

// Example 9: Conditional Mutation
func demoConditionalMutation(ctx context.Context) {
	fmt.Println("\n=== Demo 9: Conditional Mutation ===")

	// Only publish if user has permission
	mutation := `
		mutation PublishPostConditional($postId: ID!, $force: Boolean!) {
			publishPost(id: $postId) @include(if: $force) {
				id
				published
			}
		}
	`

	variables := map[string]interface{}{
		"postId": "3",
		"force":  true,
	}

	fmt.Printf("Mutation:\n%s\n", mutation)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nUse case: Conditional mutations based on business logic")
}

// Example 10: Optimistic UI Update Pattern
func demoOptimisticUpdate(ctx context.Context) {
	fmt.Println("\n=== Demo 10: Optimistic UI Update ===")

	fmt.Println("Pattern: Update UI immediately, then send mutation")
	fmt.Println("1. User clicks 'Like' button")
	fmt.Println("2. UI updates: Like count +1 (optimistic)")
	fmt.Println("3. Send mutation:")
	mutation := `
		mutation LikePost($postId: ID!) {
			likePost(postId: $postId) {
				id
				likeCount
			}
		}
	`
	fmt.Printf("%s\n", mutation)
	fmt.Println("4. On success: Keep optimistic update")
	fmt.Println("5. On error: Revert UI, show error")
	fmt.Println("\nUse case: Better perceived performance")
}

// ============================================
// Helper Functions
// ============================================

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

// ============================================
// Usage Examples
// ============================================

func ExampleMutations() {
	ctx := context.Background()

	fmt.Println("✏️  GraphQL Mutation Examples")
	fmt.Println("============================")

	demoCreateUser(ctx)
	demoUpdateUser(ctx)
	demoDeleteUser(ctx)
	demoCreatePost(ctx)
	demoPublishPost(ctx)
	demoCreateComment(ctx)
	demoBatchMutations(ctx)
	demoMutationWithError(ctx)
	demoConditionalMutation(ctx)
	demoOptimisticUpdate(ctx)

	fmt.Println("\n✅ Mutation examples completed")
}

/*
Common Mutation Patterns:

1. Create:
   mutation { createUser(input: {...}) { id username } }

2. Update:
   mutation { updateUser(id: "1", input: {...}) { id username } }

3. Delete:
   mutation { deleteUser(id: "1") }

4. Toggle:
   mutation { publishPost(id: "1") { id published } }

5. Batch:
   mutation { createUser(...) { id } createPost(...) { id } }

Best Practices:

1. Always validate input in resolvers
2. Return meaningful error messages
3. Use transactions for related mutations
4. Implement proper authentication/authorization
5. Return created/updated objects for confirmation
6. Handle partial failures gracefully
7. Use input types for complex inputs
8. Document mutation side effects

Security Considerations:

1. Validate all inputs
2. Check permissions before mutations
3. Rate limit mutations
4. Sanitize error messages (don't leak internal info)
5. Use HTTPS in production
6. Implement CSRF protection
7. Audit log mutations
8. Implement idempotency where possible

Error Handling Patterns:

1. Validation Errors: Return specific field errors
2. Not Found: Return null with error
3. Unauthorized: 401 with clear message
4. Forbidden: 403 with permission info
5. Server Error: Generic message, log details
6. Partial Success: Return data + errors
*/

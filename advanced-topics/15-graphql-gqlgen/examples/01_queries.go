package main

import (
	"context"
	"fmt"
)

// ============================================
// GraphQL Query Examples
// ============================================
// This file demonstrates various query patterns in GraphQL

// Example 1: Basic Query - Fetch all users
func demoBasicQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 1: Basic Query ===")

	query := `
		query {
			users {
				id
				username
				email
				role
				active
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)
	fmt.Println("\nResult: Returns array of users with requested fields")
	fmt.Println("Use case: Simple data fetching without filters")
}

// Example 2: Query with Arguments - Filter users
func demoQueryWithArguments(ctx context.Context) {
	fmt.Println("\n=== Demo 2: Query with Arguments ===")

	query := `
		query GetActiveUsers($role: Role!, $active: Boolean!) {
			users(role: $role, active: $active, page: {limit: 10, offset: 0}) {
				id
				username
				email
				role
			}
		}
	`

	variables := map[string]interface{}{
		"role":   "USER",
		"active": true,
	}

	fmt.Printf("Query:\n%s\n", query)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: Only active users with USER role")
	fmt.Println("Use case: Filtering and pagination")
}

// Example 3: Nested Query - User with posts and comments
func demoNestedQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 3: Nested Query ===")

	query := `
		query GetUserWithPosts($userId: ID!) {
			user(id: $userId) {
				id
				username
				email
				posts(limit: 5) {
					id
					title
					content
					published
					tags
					comments {
						id
						content
						author {
							username
						}
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userId": "2",
	}

	fmt.Printf("Query:\n%s\n", query)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: User with nested posts and comments")
	fmt.Println("Use case: Fetching related data in single request")
}

// Example 4: Multiple Root Queries - Batch fetching
func demoBatchQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 4: Batch Queries ===")

	query := `
		query GetDashboardData {
			stats {
				users
				posts
				comments
				publishedPosts
			}
			posts(published: true, limit: 5) {
				id
				title
				author {
					username
				}
			}
			users(active: true, limit: 5) {
				id
				username
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)
	fmt.Println("\nResult: Multiple queries in single request")
	fmt.Println("Use case: Dashboard data, reduces HTTP requests")
}

// Example 5: Query with Search
func demoSearchQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 5: Search Query ===")

	query := `
		query SearchContent($searchQuery: String!) {
			searchPosts(query: $searchQuery, limit: 10) {
				id
				title
				content
				tags
				author {
					username
				}
			}
			searchUsers(query: $searchQuery, limit: 5) {
				id
				username
				email
			}
		}
	`

	variables := map[string]interface{}{
		"searchQuery": "graphql",
	}

	fmt.Printf("Query:\n%s\n", query)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: Posts and users matching search term")
	fmt.Println("Use case: Search functionality")
}

// Example 6: Query with Sorting
func demoSortedQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 6: Sorted Query ===")

	query := `
		query {
			posts(sort: CREATED_AT_DESC, page: {limit: 10}) {
				id
				title
				createdAt
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)
	fmt.Println("\nResult: Posts sorted by creation date (newest first)")
	fmt.Println("Use case: Ordered lists, feeds")
}

// Example 7: Query with Fragments
func demoFragmentQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 7: Query with Fragments ===")

	query := `
		fragment UserInfo on User {
			id
			username
			email
			role
		}

		fragment PostInfo on Post {
			id
			title
			content
			published
			tags
		}

		query GetUserData {
			users {
				...UserInfo
			}
			posts {
				...PostInfo
				author {
					...UserInfo
				}
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)
	fmt.Println("\nResult: Reusable field definitions")
	fmt.Println("Use case: DRY principle, consistent field selection")
}

// Example 8: Query with Directives (Conditional Fields)
func demoDirectiveQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 8: Query with Directives ===")

	query := `
		query GetUser($withEmail: Boolean!, $withPosts: Boolean!) {
			user(id: "2") {
				id
				username
				email @include(if: $withEmail)
				posts @skip(if: $withPosts) {
					id
					title
				}
			}
		}
	`

	variables := map[string]interface{}{
		"withEmail": true,
		"withPosts": false,
	}

	fmt.Printf("Query:\n%s\n", query)
	fmt.Printf("Variables: %v\n", variables)
	fmt.Println("\nResult: Email included, posts excluded")
	fmt.Println("Use case: Conditional field selection based on user permissions")
}

// Example 9: Query with Aliases
func demoAliasQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 9: Query with Aliases ===")

	query := `
		query {
			adminUser: user(id: "1") {
				id
				username
				role
			}
			regularUser: user(id: "2") {
				id
				username
				role
			}
			publishedPosts: posts(published: true) {
				id
				title
			}
			draftPosts: posts(published: false) {
				id
				title
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)
	fmt.Println("\nResult: Same field with different names")
	fmt.Println("Use case: Fetching same resource with different parameters")
}

// Example 10: Introspection Query
func demoIntrospectionQuery(ctx context.Context) {
	fmt.Println("\n=== Demo 10: Introspection Query ===")

	query := `
		query {
			__schema {
				types {
					name
					kind
					description
					fields {
						name
						type {
							name
						}
					}
				}
			}
		}
	`

	fmt.Printf("Query:\n%s\n", query)
	fmt.Println("\nResult: Schema metadata")
	fmt.Println("Use case: Documentation, tooling, autocomplete")
}

// ============================================
// Usage Examples
// ============================================

func ExampleQueries() {
	ctx := context.Background()

	fmt.Println("📝 GraphQL Query Examples")
	fmt.Println("==========================")

	demoBasicQuery(ctx)
	demoQueryWithArguments(ctx)
	demoNestedQuery(ctx)
	demoBatchQuery(ctx)
	demoSearchQuery(ctx)
	demoSortedQuery(ctx)
	demoFragmentQuery(ctx)
	demoDirectiveQuery(ctx)
	demoAliasQuery(ctx)
	demoIntrospectionQuery(ctx)

	fmt.Println("\n✅ Query examples completed")
}

/*
Common Query Patterns:

1. Simple Fetch:
   query { users { id username } }

2. Filtered Fetch:
   query { users(role: USER, active: true) { id username } }

3. Nested Fetch:
   query { user(id: "1") { id posts { title } } }

4. Batch Fetch:
   query { users { id } posts { id } stats { users } }

5. Search:
   query { searchPosts(query: "graphql") { id title } }

6. Pagination:
   query { posts(page: {limit: 10, offset: 20}) { id title } }

7. Sorting:
   query { posts(sort: CREATED_AT_DESC) { id title } }

Performance Tips:

1. Request only needed fields (over-fetching is wasteful)
2. Use batching to reduce round trips
3. Implement pagination for large lists
4. Use data loaders to prevent N+1 queries
5. Add query complexity limits
6. Cache frequently accessed data
7. Use persisted queries for security
8. Monitor and optimize slow queries
*/

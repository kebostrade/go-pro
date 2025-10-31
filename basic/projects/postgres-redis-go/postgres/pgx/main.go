package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                    PostgreSQL with pgx/v5 - Tutorial                         ║
║                                                                              ║
║  pgx is a pure Go PostgreSQL driver and toolkit with:                       ║
║  • High performance (faster than database/sql)                              ║
║  • Native PostgreSQL features support                                       ║
║  • Connection pooling with pgxpool                                          ║
║  • Prepared statements and batch operations                                 ║
║  • Copy protocol support                                                    ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

type User struct {
	ID        int
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   time.Time
}

func main() {
	fmt.Println("🐘 PostgreSQL with pgx/v5 Tutorial")
	fmt.Println("=" + string(make([]byte, 50)))

	ctx := context.Background()

	// Example 1: Simple Connection
	fmt.Println("\n📌 Example 1: Simple Connection")
	simpleConnection(ctx)

	// Example 2: Connection Pool
	fmt.Println("\n📌 Example 2: Connection Pool (Recommended)")
	pool := createConnectionPool(ctx)
	defer pool.Close()

	// Example 3: CRUD Operations
	fmt.Println("\n📌 Example 3: CRUD Operations")
	crudOperations(ctx, pool)

	// Example 4: Transactions
	fmt.Println("\n📌 Example 4: Transactions")
	transactionExample(ctx, pool)

	// Example 5: Batch Operations
	fmt.Println("\n📌 Example 5: Batch Operations")
	batchOperations(ctx, pool)

	// Example 6: Prepared Statements
	fmt.Println("\n📌 Example 6: Prepared Statements")
	preparedStatements(ctx, pool)

	// Example 7: Query Rows
	fmt.Println("\n📌 Example 7: Query Multiple Rows")
	queryRows(ctx, pool)

	fmt.Println("\n✅ All examples completed successfully!")
}

// Example 1: Simple connection (not recommended for production)
func simpleConnection(ctx context.Context) {
	connString := "postgres://postgres:postgres@localhost:5432/go_tutorial?sslmode=disable"
	
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Printf("❌ Unable to connect: %v\n", err)
		return
	}
	defer conn.Close(ctx)

	var version string
	err = conn.QueryRow(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		log.Printf("❌ Query failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Connected! PostgreSQL version: %s\n", version[:50]+"...")
}

// Example 2: Connection pool (recommended for production)
func createConnectionPool(ctx context.Context) *pgxpool.Pool {
	connString := "postgres://postgres:postgres@localhost:5432/go_tutorial?sslmode=disable"

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatalf("❌ Unable to parse config: %v\n", err)
	}

	// Configure connection pool
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("❌ Unable to create pool: %v\n", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Unable to ping database: %v\n", err)
	}

	fmt.Printf("✅ Connection pool created (Max: %d, Min: %d)\n", 
		config.MaxConns, config.MinConns)

	return pool
}

// Example 3: CRUD Operations
func crudOperations(ctx context.Context, pool *pgxpool.Pool) {
	// Create table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Printf("❌ Create table failed: %v\n", err)
		return
	}

	// INSERT
	var userID int
	err = pool.QueryRow(ctx,
		"INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id",
		"john_doe", "john@example.com",
	).Scan(&userID)
	if err != nil {
		log.Printf("❌ Insert failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Created user with ID: %d\n", userID)

	// SELECT
	var user User
	err = pool.QueryRow(ctx,
		"SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Created At, &user.UpdatedAt)
	if err != nil {
		log.Printf("❌ Select failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Found user: %s (%s)\n", user.Username, user.Email)

	// UPDATE
	tag, err := pool.Exec(ctx,
		"UPDATE users SET email = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		"john.doe@example.com", userID,
	)
	if err != nil {
		log.Printf("❌ Update failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Updated %d row(s)\n", tag.RowsAffected())

	// DELETE
	tag, err = pool.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		log.Printf("❌ Delete failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Deleted %d row(s)\n", tag.RowsAffected())
}

// Example 4: Transactions
func transactionExample(ctx context.Context, pool *pgxpool.Pool) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Printf("❌ Begin transaction failed: %v\n", err)
		return
	}
	defer tx.Rollback(ctx) // Rollback if not committed

	// Insert multiple users in a transaction
	users := []struct {
		username string
		email    string
	}{
		{"alice", "alice@example.com"},
		{"bob", "bob@example.com"},
		{"charlie", "charlie@example.com"},
	}

	for _, u := range users {
		_, err := tx.Exec(ctx,
			"INSERT INTO users (username, email) VALUES ($1, $2)",
			u.username, u.email,
		)
		if err != nil {
			log.Printf("❌ Insert in transaction failed: %v\n", err)
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		log.Printf("❌ Commit failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Transaction committed: inserted %d users\n", len(users))
}

// Example 5: Batch Operations
func batchOperations(ctx context.Context, pool *pgxpool.Pool) {
	batch := &pgx.Batch{}

	// Add multiple queries to batch
	for i := 1; i <= 5; i++ {
		batch.Queue(
			"INSERT INTO users (username, email) VALUES ($1, $2)",
			fmt.Sprintf("user%d", i),
			fmt.Sprintf("user%d@example.com", i),
		)
	}

	// Execute batch
	results := pool.SendBatch(ctx, batch)
	defer results.Close()

	// Process results
	var count int
	for i := 0; i < batch.Len(); i++ {
		_, err := results.Exec()
		if err != nil {
			log.Printf("❌ Batch operation %d failed: %v\n", i, err)
			continue
		}
		count++
	}

	fmt.Printf("✅ Batch completed: %d/%d operations successful\n", count, batch.Len())
}

// Example 6: Prepared Statements
func preparedStatements(ctx context.Context, pool *pgxpool.Pool) {
	// Prepare statement
	stmt := "SELECT id, username, email FROM users WHERE username = $1"

	var user User
	err := pool.QueryRow(ctx, stmt, "alice").Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		log.Printf("❌ Prepared statement failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Found user via prepared statement: %s\n", user.Username)
}

// Example 7: Query Multiple Rows
func queryRows(ctx context.Context, pool *pgxpool.Pool) {
	rows, err := pool.Query(ctx, "SELECT id, username, email FROM users LIMIT 10")
	if err != nil {
		log.Printf("❌ Query rows failed: %v\n", err)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			log.Printf("❌ Scan failed: %v\n", err)
			continue
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("❌ Rows error: %v\n", err)
		return
	}

	fmt.Printf("✅ Retrieved %d users\n", len(users))
	for _, u := range users {
		fmt.Printf("   - %s (%s)\n", u.Username, u.Email)
	}
}


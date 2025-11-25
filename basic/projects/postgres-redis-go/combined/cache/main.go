package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║              PostgreSQL + Redis: Cache-Aside Pattern                        ║
║                                                                              ║
║  Cache-Aside (Lazy Loading):                                                ║
║  1. Check cache first                                                        ║
║  2. If miss, query database                                                  ║
║  3. Store result in cache                                                    ║
║  4. Return result                                                            ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserService struct {
	db    *pgxpool.Pool
	cache *redis.Client
}

func main() {
	fmt.Println("💾 Cache-Aside Pattern Tutorial")
	fmt.Println("=" + string(make([]byte, 50)))

	ctx := context.Background()

	// Connect to PostgreSQL
	db := connectPostgres(ctx)
	defer db.Close()

	// Connect to Redis
	rdb := connectRedis()
	defer rdb.Close()

	// Create service
	service := &UserService{
		db:    db,
		cache: rdb,
	}

	// Setup database
	setupDatabase(ctx, db)

	// Example 1: Cache-Aside Pattern
	fmt.Println("\n📌 Example 1: Cache-Aside Pattern")
	cacheAsideExample(ctx, service)

	// Example 2: Cache Invalidation
	fmt.Println("\n📌 Example 2: Cache Invalidation")
	cacheInvalidationExample(ctx, service)

	// Example 3: Write-Through Cache
	fmt.Println("\n📌 Example 3: Write-Through Cache")
	writeThroughExample(ctx, service)

	// Example 4: Cache Statistics
	fmt.Println("\n📌 Example 4: Cache Statistics")
	cacheStatistics(ctx, service)

	fmt.Println("\n✅ All cache pattern examples completed!")
}

func connectPostgres(ctx context.Context) *pgxpool.Pool {
	connString := "postgres://postgres:postgres@localhost:5432/go_tutorial?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Failed to ping PostgreSQL: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
	return pool
}

func connectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	fmt.Println("✅ Connected to Redis")
	return rdb
}

func setupDatabase(ctx context.Context, db *pgxpool.Pool) {
	// Create table
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("❌ Failed to create table: %v", err)
	}

	// Insert sample data
	users := []struct {
		username string
		email    string
	}{
		{"alice", "alice@example.com"},
		{"bob", "bob@example.com"},
		{"charlie", "charlie@example.com"},
	}

	for _, u := range users {
		_, err := db.Exec(ctx,
			"INSERT INTO users (username, email) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING",
			u.username, u.email,
		)
		if err != nil {
			log.Printf("⚠️  Failed to insert user: %v", err)
		}
	}

	fmt.Println("✅ Database setup completed")
}

// Example 1: Cache-Aside Pattern
func cacheAsideExample(ctx context.Context, service *UserService) {
	userID := 1

	// First call - cache miss
	fmt.Println("\n🔍 First call (cache miss):")
	start := time.Now()
	user1, err := service.GetUser(ctx, userID)
	duration1 := time.Since(start)
	if err != nil {
		log.Printf("❌ Failed to get user: %v\n", err)
		return
	}
	fmt.Printf("✅ Retrieved user: %s (%s) in %v\n", user1.Username, user1.Email, duration1)

	// Second call - cache hit
	fmt.Println("\n🔍 Second call (cache hit):")
	start = time.Now()
	user2, err := service.GetUser(ctx, userID)
	duration2 := time.Since(start)
	if err != nil {
		log.Printf("❌ Failed to get user: %v\n", err)
		return
	}
	fmt.Printf("✅ Retrieved user: %s (%s) in %v\n", user2.Username, user2.Email, duration2)

	fmt.Printf("\n📊 Performance improvement: %.2fx faster\n", float64(duration1)/float64(duration2))
}

// GetUser implements cache-aside pattern
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// 1. Check cache
	cached, err := s.cache.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache hit
		var user User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			fmt.Println("   💚 Cache HIT")
			return &user, nil
		}
	}

	// 2. Cache miss - query database
	fmt.Println("   💔 Cache MISS - querying database")
	var user User
	err = s.db.QueryRow(ctx,
		"SELECT id, username, email, created_at FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	// 3. Store in cache
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = s.cache.Set(ctx, cacheKey, data, 5*time.Minute).Err()
	if err != nil {
		log.Printf("⚠️  Failed to cache user: %v", err)
	}

	return &user, nil
}

// Example 2: Cache Invalidation
func cacheInvalidationExample(ctx context.Context, service *UserService) {
	userID := 1

	// Get user (cache it)
	user, err := service.GetUser(ctx, userID)
	if err != nil {
		log.Printf("❌ Failed to get user: %v\n", err)
		return
	}
	fmt.Printf("✅ User before update: %s\n", user.Email)

	// Update user
	newEmail := "alice.updated@example.com"
	err = service.UpdateUser(ctx, userID, newEmail)
	if err != nil {
		log.Printf("❌ Failed to update user: %v\n", err)
		return
	}
	fmt.Printf("✅ User updated with new email: %s\n", newEmail)

	// Get user again (should reflect update)
	user, err = service.GetUser(ctx, userID)
	if err != nil {
		log.Printf("❌ Failed to get user: %v\n", err)
		return
	}
	fmt.Printf("✅ User after update: %s\n", user.Email)
}

// UpdateUser with cache invalidation
func (s *UserService) UpdateUser(ctx context.Context, id int, email string) error {
	// 1. Update database
	_, err := s.db.Exec(ctx,
		"UPDATE users SET email = $1 WHERE id = $2",
		email, id,
	)
	if err != nil {
		return err
	}

	// 2. Invalidate cache
	cacheKey := fmt.Sprintf("user:%d", id)
	err = s.cache.Del(ctx, cacheKey).Err()
	if err != nil {
		log.Printf("⚠️  Failed to invalidate cache: %v", err)
	} else {
		fmt.Println("   🗑️  Cache invalidated")
	}

	return nil
}

// Example 3: Write-Through Cache
func writeThroughExample(ctx context.Context, service *UserService) {
	user := &User{
		Username: "david",
		Email:    "david@example.com",
	}

	err := service.CreateUser(ctx, user)
	if err != nil {
		log.Printf("❌ Failed to create user: %v\n", err)
		return
	}

	fmt.Printf("✅ User created: ID=%d, Username=%s\n", user.ID, user.Username)

	// Verify cache
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	exists, err := service.cache.Exists(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("❌ Failed to check cache: %v\n", err)
		return
	}

	if exists > 0 {
		fmt.Println("✅ User is cached (write-through)")
	}
}

// CreateUser with write-through caching
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	// 1. Insert into database
	err := s.db.QueryRow(ctx,
		"INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, created_at",
		user.Username, user.Email,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	// 2. Write to cache
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = s.cache.Set(ctx, cacheKey, data, 5*time.Minute).Err()
	if err != nil {
		log.Printf("⚠️  Failed to cache new user: %v", err)
	} else {
		fmt.Println("   💾 User cached (write-through)")
	}

	return nil
}

// Example 4: Cache Statistics
func cacheStatistics(ctx context.Context, service *UserService) {
	// Simulate multiple requests
	userIDs := []int{1, 2, 3, 1, 2, 1, 3, 1}

	hits := 0
	misses := 0

	for _, id := range userIDs {
		cacheKey := fmt.Sprintf("user:%d", id)
		_, err := service.cache.Get(ctx, cacheKey).Result()
		if err == nil {
			hits++
		} else {
			misses++
			// Simulate database query and cache
			service.GetUser(ctx, id)
		}
	}

	total := hits + misses
	hitRate := float64(hits) / float64(total) * 100

	fmt.Printf("\n📊 Cache Statistics:\n")
	fmt.Printf("   Total requests: %d\n", total)
	fmt.Printf("   Cache hits: %d\n", hits)
	fmt.Printf("   Cache misses: %d\n", misses)
	fmt.Printf("   Hit rate: %.2f%%\n", hitRate)
}


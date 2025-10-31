package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                    Redis Patterns - Tutorial                                 ║
║                                                                              ║
║  Common Redis patterns:                                                     ║
║  • Distributed Locks                                                         ║
║  • Rate Limiting                                                             ║
║  • Leaderboards                                                              ║
║  • Session Management                                                        ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

func main() {
	fmt.Println("🎯 Redis Patterns Tutorial")
	fmt.Println("=" + string(make([]byte, 50)))

	ctx := context.Background()
	rdb := connectRedis()
	defer rdb.Close()

	// Example 1: Distributed Lock
	fmt.Println("\n📌 Example 1: Distributed Lock")
	distributedLock(ctx, rdb)

	// Example 2: Rate Limiting
	fmt.Println("\n📌 Example 2: Rate Limiting")
	rateLimiting(ctx, rdb)

	// Example 3: Leaderboard
	fmt.Println("\n📌 Example 3: Leaderboard")
	leaderboard(ctx, rdb)

	// Example 4: Session Management
	fmt.Println("\n📌 Example 4: Session Management")
	sessionManagement(ctx, rdb)

	// Example 5: Counter Pattern
	fmt.Println("\n📌 Example 5: Counter Pattern")
	counterPattern(ctx, rdb)

	fmt.Println("\n✅ All pattern examples completed!")
}

func connectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}

	fmt.Println("✅ Connected to Redis")
	return rdb
}

// Example 1: Distributed Lock
func distributedLock(ctx context.Context, rdb *redis.Client) {
	lockKey := "resource:lock"
	lockValue := "unique-lock-id-123"
	lockTTL := 10 * time.Second

	// Acquire lock
	acquired, err := rdb.SetNX(ctx, lockKey, lockValue, lockTTL).Result()
	if err != nil {
		log.Printf("❌ Lock acquisition failed: %v\n", err)
		return
	}

	if acquired {
		fmt.Println("✅ Lock acquired successfully")

		// Simulate work
		fmt.Println("   Performing critical operation...")
		time.Sleep(2 * time.Second)

		// Release lock
		script := `
			if redis.call("get", KEYS[1]) == ARGV[1] then
				return redis.call("del", KEYS[1])
			else
				return 0
			end
		`
		result, err := rdb.Eval(ctx, script, []string{lockKey}, lockValue).Result()
		if err != nil {
			log.Printf("❌ Lock release failed: %v\n", err)
			return
		}

		if result.(int64) == 1 {
			fmt.Println("✅ Lock released successfully")
		} else {
			fmt.Println("⚠️  Lock was already released or expired")
		}
	} else {
		fmt.Println("❌ Failed to acquire lock (already locked)")
	}
}

// Example 2: Rate Limiting (Fixed Window)
func rateLimiting(ctx context.Context, rdb *redis.Client) {
	userID := "user:123"
	limit := 5
	window := 60 * time.Second

	for i := 1; i <= 7; i++ {
		allowed := checkRateLimit(ctx, rdb, userID, limit, window)
		if allowed {
			fmt.Printf("✅ Request %d: Allowed\n", i)
		} else {
			fmt.Printf("❌ Request %d: Rate limit exceeded\n", i)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func checkRateLimit(ctx context.Context, rdb *redis.Client, userID string, limit int, window time.Duration) bool {
	key := fmt.Sprintf("rate_limit:%s", userID)

	// Increment counter
	count, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		log.Printf("❌ Rate limit check failed: %v\n", err)
		return false
	}

	// Set expiration on first request
	if count == 1 {
		rdb.Expire(ctx, key, window)
	}

	return count <= int64(limit)
}

// Example 3: Leaderboard
func leaderboard(ctx context.Context, rdb *redis.Client) {
	leaderboardKey := "game:leaderboard"

	// Add players with scores
	players := map[string]float64{
		"Alice":   1500,
		"Bob":     2000,
		"Charlie": 1800,
		"David":   2200,
		"Eve":     1900,
	}

	for player, score := range players {
		err := rdb.ZAdd(ctx, leaderboardKey, redis.Z{
			Score:  score,
			Member: player,
		}).Err()
		if err != nil {
			log.Printf("❌ Failed to add player: %v\n", err)
			continue
		}
	}
	fmt.Printf("✅ Added %d players to leaderboard\n", len(players))

	// Get top 3 players
	top3, err := rdb.ZRevRangeWithScores(ctx, leaderboardKey, 0, 2).Result()
	if err != nil {
		log.Printf("❌ Failed to get top players: %v\n", err)
		return
	}

	fmt.Println("🏆 Top 3 Players:")
	for i, player := range top3 {
		fmt.Printf("   #%d: %s (%.0f points)\n", i+1, player.Member, player.Score)
	}

	// Get player rank
	rank, err := rdb.ZRevRank(ctx, leaderboardKey, "Charlie").Result()
	if err != nil {
		log.Printf("❌ Failed to get rank: %v\n", err)
		return
	}
	fmt.Printf("✅ Charlie's rank: #%d\n", rank+1)

	// Increment player score
	newScore, err := rdb.ZIncrBy(ctx, leaderboardKey, 300, "Charlie").Result()
	if err != nil {
		log.Printf("❌ Failed to increment score: %v\n", err)
		return
	}
	fmt.Printf("✅ Charlie's new score: %.0f\n", newScore)
}

// Example 4: Session Management
func sessionManagement(ctx context.Context, rdb *redis.Client) {
	sessionID := "session:abc123"
	sessionTTL := 30 * time.Minute

	// Create session
	sessionData := map[string]interface{}{
		"user_id":    "123",
		"username":   "john_doe",
		"email":      "john@example.com",
		"login_time": time.Now().Unix(),
	}

	err := rdb.HSet(ctx, sessionID, sessionData).Err()
	if err != nil {
		log.Printf("❌ Session creation failed: %v\n", err)
		return
	}

	// Set expiration
	err = rdb.Expire(ctx, sessionID, sessionTTL).Err()
	if err != nil {
		log.Printf("❌ Setting expiration failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Session created with %v TTL\n", sessionTTL)

	// Get session data
	session, err := rdb.HGetAll(ctx, sessionID).Result()
	if err != nil {
		log.Printf("❌ Get session failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Session data: %v\n", session)

	// Update session activity
	err = rdb.HSet(ctx, sessionID, "last_activity", time.Now().Unix()).Err()
	if err != nil {
		log.Printf("❌ Update session failed: %v\n", err)
		return
	}

	// Refresh TTL
	err = rdb.Expire(ctx, sessionID, sessionTTL).Err()
	if err != nil {
		log.Printf("❌ Refresh TTL failed: %v\n", err)
		return
	}
	fmt.Println("✅ Session activity updated and TTL refreshed")

	// Delete session (logout)
	err = rdb.Del(ctx, sessionID).Err()
	if err != nil {
		log.Printf("❌ Delete session failed: %v\n", err)
		return
	}
	fmt.Println("✅ Session deleted (user logged out)")
}

// Example 5: Counter Pattern
func counterPattern(ctx context.Context, rdb *redis.Client) {
	// Page views counter
	pageKey := "page:home:views"

	// Increment views
	for i := 0; i < 10; i++ {
		views, err := rdb.Incr(ctx, pageKey).Result()
		if err != nil {
			log.Printf("❌ Increment failed: %v\n", err)
			continue
		}
		if i == 9 {
			fmt.Printf("✅ Page views: %d\n", views)
		}
	}

	// Daily counter with expiration
	today := time.Now().Format("2006-01-02")
	dailyKey := fmt.Sprintf("page:home:views:%s", today)

	views, err := rdb.Incr(ctx, dailyKey).Result()
	if err != nil {
		log.Printf("❌ Daily increment failed: %v\n", err)
		return
	}

	// Set expiration to end of day
	tomorrow := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	ttl := time.Until(tomorrow)
	rdb.Expire(ctx, dailyKey, ttl)

	fmt.Printf("✅ Daily views (%s): %d (expires in %v)\n", today, views, ttl.Round(time.Second))
}


package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                      Redis Basics - Tutorial                                 ║
║                                                                              ║
║  Redis is an in-memory data structure store used as:                        ║
║  • Database                                                                  ║
║  • Cache                                                                     ║
║  • Message broker                                                            ║
║  • Session store                                                             ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	fmt.Println("🔴 Redis Basics Tutorial")
	fmt.Println("=" + string(make([]byte, 50)))

	ctx := context.Background()

	// Connect to Redis
	rdb := connectRedis()
	defer rdb.Close()

	// Example 1: String Operations
	fmt.Println("\n📌 Example 1: String Operations")
	stringOperations(ctx, rdb)

	// Example 2: Hash Operations
	fmt.Println("\n📌 Example 2: Hash Operations")
	hashOperations(ctx, rdb)

	// Example 3: List Operations
	fmt.Println("\n📌 Example 3: List Operations")
	listOperations(ctx, rdb)

	// Example 4: Set Operations
	fmt.Println("\n📌 Example 4: Set Operations")
	setOperations(ctx, rdb)

	// Example 5: Sorted Set Operations
	fmt.Println("\n📌 Example 5: Sorted Set Operations")
	sortedSetOperations(ctx, rdb)

	// Example 6: Expiration and TTL
	fmt.Println("\n📌 Example 6: Expiration and TTL")
	expirationOperations(ctx, rdb)

	// Example 7: JSON Storage
	fmt.Println("\n📌 Example 7: JSON Storage")
	jsonOperations(ctx, rdb)

	// Example 8: Pipelining
	fmt.Println("\n📌 Example 8: Pipelining")
	pipelining(ctx, rdb)

	// Example 9: Transactions
	fmt.Println("\n📌 Example 9: Transactions (MULTI/EXEC)")
	transactions(ctx, rdb)

	fmt.Println("\n✅ All Redis basic examples completed!")
}

func connectRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		MaxRetries:   3,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	fmt.Println("✅ Connected to Redis")
	return rdb
}

// Example 1: String Operations
func stringOperations(ctx context.Context, rdb *redis.Client) {
	// SET
	err := rdb.Set(ctx, "user:1:name", "John Doe", 0).Err()
	if err != nil {
		log.Printf("❌ SET failed: %v\n", err)
		return
	}
	fmt.Println("✅ SET user:1:name = 'John Doe'")

	// GET
	val, err := rdb.Get(ctx, "user:1:name").Result()
	if err != nil {
		log.Printf("❌ GET failed: %v\n", err)
		return
	}
	fmt.Printf("✅ GET user:1:name = '%s'\n", val)

	// INCR
	count, err := rdb.Incr(ctx, "page:views").Result()
	if err != nil {
		log.Printf("❌ INCR failed: %v\n", err)
		return
	}
	fmt.Printf("✅ INCR page:views = %d\n", count)

	// MSET (Multiple SET)
	err = rdb.MSet(ctx,
		"user:2:name", "Alice",
		"user:3:name", "Bob",
	).Err()
	if err != nil {
		log.Printf("❌ MSET failed: %v\n", err)
		return
	}
	fmt.Println("✅ MSET multiple keys")

	// MGET (Multiple GET)
	vals, err := rdb.MGet(ctx, "user:1:name", "user:2:name", "user:3:name").Result()
	if err != nil {
		log.Printf("❌ MGET failed: %v\n", err)
		return
	}
	fmt.Printf("✅ MGET = %v\n", vals)
}

// Example 2: Hash Operations
func hashOperations(ctx context.Context, rdb *redis.Client) {
	// HSET
	err := rdb.HSet(ctx, "user:100", map[string]interface{}{
		"username": "john_doe",
		"email":    "john@example.com",
		"age":      30,
	}).Err()
	if err != nil {
		log.Printf("❌ HSET failed: %v\n", err)
		return
	}
	fmt.Println("✅ HSET user:100 with multiple fields")

	// HGET
	username, err := rdb.HGet(ctx, "user:100", "username").Result()
	if err != nil {
		log.Printf("❌ HGET failed: %v\n", err)
		return
	}
	fmt.Printf("✅ HGET user:100 username = '%s'\n", username)

	// HGETALL
	user, err := rdb.HGetAll(ctx, "user:100").Result()
	if err != nil {
		log.Printf("❌ HGETALL failed: %v\n", err)
		return
	}
	fmt.Printf("✅ HGETALL user:100 = %v\n", user)

	// HINCRBY
	newAge, err := rdb.HIncrBy(ctx, "user:100", "age", 1).Result()
	if err != nil {
		log.Printf("❌ HINCRBY failed: %v\n", err)
		return
	}
	fmt.Printf("✅ HINCRBY user:100 age = %d\n", newAge)
}

// Example 3: List Operations
func listOperations(ctx context.Context, rdb *redis.Client) {
	// LPUSH (Left Push)
	err := rdb.LPush(ctx, "tasks", "task1", "task2", "task3").Err()
	if err != nil {
		log.Printf("❌ LPUSH failed: %v\n", err)
		return
	}
	fmt.Println("✅ LPUSH 3 tasks")

	// RPUSH (Right Push)
	err = rdb.RPush(ctx, "tasks", "task4").Err()
	if err != nil {
		log.Printf("❌ RPUSH failed: %v\n", err)
		return
	}
	fmt.Println("✅ RPUSH task4")

	// LRANGE (Get range)
	tasks, err := rdb.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		log.Printf("❌ LRANGE failed: %v\n", err)
		return
	}
	fmt.Printf("✅ LRANGE tasks = %v\n", tasks)

	// LPOP (Left Pop)
	task, err := rdb.LPop(ctx, "tasks").Result()
	if err != nil {
		log.Printf("❌ LPOP failed: %v\n", err)
		return
	}
	fmt.Printf("✅ LPOP = '%s'\n", task)

	// LLEN (List Length)
	length, err := rdb.LLen(ctx, "tasks").Result()
	if err != nil {
		log.Printf("❌ LLEN failed: %v\n", err)
		return
	}
	fmt.Printf("✅ LLEN tasks = %d\n", length)
}

// Example 4: Set Operations
func setOperations(ctx context.Context, rdb *redis.Client) {
	// SADD
	err := rdb.SAdd(ctx, "tags", "golang", "redis", "tutorial").Err()
	if err != nil {
		log.Printf("❌ SADD failed: %v\n", err)
		return
	}
	fmt.Println("✅ SADD 3 tags")

	// SMEMBERS
	tags, err := rdb.SMembers(ctx, "tags").Result()
	if err != nil {
		log.Printf("❌ SMEMBERS failed: %v\n", err)
		return
	}
	fmt.Printf("✅ SMEMBERS tags = %v\n", tags)

	// SISMEMBER
	exists, err := rdb.SIsMember(ctx, "tags", "golang").Result()
	if err != nil {
		log.Printf("❌ SISMEMBER failed: %v\n", err)
		return
	}
	fmt.Printf("✅ SISMEMBER tags 'golang' = %v\n", exists)

	// SCARD (Set Cardinality)
	count, err := rdb.SCard(ctx, "tags").Result()
	if err != nil {
		log.Printf("❌ SCARD failed: %v\n", err)
		return
	}
	fmt.Printf("✅ SCARD tags = %d\n", count)
}

// Example 5: Sorted Set Operations
func sortedSetOperations(ctx context.Context, rdb *redis.Client) {
	// ZADD
	err := rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "player1"}).Err()
	rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 200, Member: "player2"})
	rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 150, Member: "player3"})
	if err != nil {
		log.Printf("❌ ZADD failed: %v\n", err)
		return
	}
	fmt.Println("✅ ZADD 3 players to leaderboard")

	// ZRANGE (Ascending)
	players, err := rdb.ZRangeWithScores(ctx, "leaderboard", 0, -1).Result()
	if err != nil {
		log.Printf("❌ ZRANGE failed: %v\n", err)
		return
	}
	fmt.Println("✅ ZRANGE leaderboard (ascending):")
	for _, p := range players {
		fmt.Printf("   %s: %.0f\n", p.Member, p.Score)
	}

	// ZREVRANGE (Descending)
	topPlayers, err := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 2).Result()
	if err != nil {
		log.Printf("❌ ZREVRANGE failed: %v\n", err)
		return
	}
	fmt.Println("✅ ZREVRANGE top 3 players:")
	for i, p := range topPlayers {
		fmt.Printf("   #%d: %s (%.0f)\n", i+1, p.Member, p.Score)
	}
}

// Example 6: Expiration Operations
func expirationOperations(ctx context.Context, rdb *redis.Client) {
	// SET with expiration
	err := rdb.Set(ctx, "session:abc123", "user_data", 10*time.Second).Err()
	if err != nil {
		log.Printf("❌ SET with expiration failed: %v\n", err)
		return
	}
	fmt.Println("✅ SET session:abc123 with 10s expiration")

	// TTL (Time To Live)
	ttl, err := rdb.TTL(ctx, "session:abc123").Result()
	if err != nil {
		log.Printf("❌ TTL failed: %v\n", err)
		return
	}
	fmt.Printf("✅ TTL session:abc123 = %v\n", ttl)

	// EXPIRE (Set expiration on existing key)
	err = rdb.Expire(ctx, "user:1:name", 60*time.Second).Err()
	if err != nil {
		log.Printf("❌ EXPIRE failed: %v\n", err)
		return
	}
	fmt.Println("✅ EXPIRE user:1:name in 60s")
}

// Example 7: JSON Operations
func jsonOperations(ctx context.Context, rdb *redis.Client) {
	user := User{
		ID:       1,
		Username: "john_doe",
		Email:    "john@example.com",
	}

	// Serialize to JSON
	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("❌ JSON marshal failed: %v\n", err)
		return
	}

	// Store JSON
	err = rdb.Set(ctx, "user:json:1", data, 0).Err()
	if err != nil {
		log.Printf("❌ SET JSON failed: %v\n", err)
		return
	}
	fmt.Println("✅ Stored user as JSON")

	// Retrieve and deserialize
	jsonData, err := rdb.Get(ctx, "user:json:1").Result()
	if err != nil {
		log.Printf("❌ GET JSON failed: %v\n", err)
		return
	}

	var retrievedUser User
	err = json.Unmarshal([]byte(jsonData), &retrievedUser)
	if err != nil {
		log.Printf("❌ JSON unmarshal failed: %v\n", err)
		return
	}
	fmt.Printf("✅ Retrieved user from JSON: %+v\n", retrievedUser)
}

// Example 8: Pipelining
func pipelining(ctx context.Context, rdb *redis.Client) {
	pipe := rdb.Pipeline()

	// Queue multiple commands
	pipe.Set(ctx, "key1", "value1", 0)
	pipe.Set(ctx, "key2", "value2", 0)
	pipe.Set(ctx, "key3", "value3", 0)
	pipe.Incr(ctx, "counter")
	pipe.Get(ctx, "key1")

	// Execute pipeline
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("❌ Pipeline failed: %v\n", err)
		return
	}

	fmt.Printf("✅ Pipeline executed: %d commands\n", len(cmds))
}

// Example 9: Transactions
func transactions(ctx context.Context, rdb *redis.Client) {
	err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		// Get current value
		val, err := tx.Get(ctx, "balance").Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// Transaction
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, "balance", val+100, 0)
			pipe.Incr(ctx, "transaction_count")
			return nil
		})
		return err
	}, "balance")

	if err != nil {
		log.Printf("❌ Transaction failed: %v\n", err)
		return
	}

	fmt.Println("✅ Transaction completed successfully")
}


package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/*
╔══════════════════════════════════════════════════════════════════════════════╗
║                         GORM ORM - Tutorial                                  ║
║                                                                              ║
║  GORM is a full-featured ORM library for Go with:                           ║
║  • Auto migrations                                                           ║
║  • Associations (Has One, Has Many, Belongs To, Many To Many)               ║
║  • Hooks (Before/After Create/Update/Delete/Find)                           ║
║  • Preloading (Eager loading)                                               ║
║  • Transactions, Nested Transactions, Save Point                            ║
╚══════════════════════════════════════════════════════════════════════════════╝
*/

// Models
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Profile   Profile   `gorm:"constraint:OnDelete:CASCADE;"`
	Posts     []Post    `gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Profile struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"uniqueIndex;not null"`
	FirstName string
	LastName  string
	Bio       string
	Avatar    string
}

type Post struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"not null"`
	Content   string
	AuthorID  uint      `gorm:"not null"`
	Author    User      `gorm:"foreignKey:AuthorID"`
	Tags      []Tag     `gorm:"many2many:post_tags;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Tag struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"uniqueIndex;not null"`
	Posts []Post `gorm:"many2many:post_tags;"`
}

func main() {
	fmt.Println("🗄️  GORM ORM Tutorial")
	fmt.Println("=" + string(make([]byte, 50)))

	// Connect to database
	db := connectDB()

	// Example 1: Auto Migration
	fmt.Println("\n📌 Example 1: Auto Migration")
	autoMigration(db)

	// Example 2: Create Records
	fmt.Println("\n📌 Example 2: Create Records")
	createRecords(db)

	// Example 3: Query Records
	fmt.Println("\n📌 Example 3: Query Records")
	queryRecords(db)

	// Example 4: Update Records
	fmt.Println("\n📌 Example 4: Update Records")
	updateRecords(db)

	// Example 5: Delete Records
	fmt.Println("\n📌 Example 5: Delete Records")
	deleteRecords(db)

	// Example 6: Associations
	fmt.Println("\n📌 Example 6: Associations")
	associations(db)

	// Example 7: Preloading
	fmt.Println("\n📌 Example 7: Preloading (Eager Loading)")
	preloading(db)

	// Example 8: Transactions
	fmt.Println("\n📌 Example 8: Transactions")
	transactions(db)

	// Example 9: Hooks
	fmt.Println("\n📌 Example 9: Hooks")
	hooks(db)

	// Example 10: Advanced Queries
	fmt.Println("\n📌 Example 10: Advanced Queries")
	advancedQueries(db)

	fmt.Println("\n✅ All GORM examples completed!")
}

func connectDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=go_tutorial port=5432 sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL with GORM")
	return db
}

func autoMigration(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Profile{}, &Post{}, &Tag{})
	if err != nil {
		log.Printf("❌ Migration failed: %v\n", err)
		return
	}
	fmt.Println("✅ Auto migration completed")
}

func createRecords(db *gorm.DB) {
	// Create single record
	user := User{
		Username: "john_doe",
		Email:    "john@example.com",
	}
	result := db.Create(&user)
	if result.Error != nil {
		log.Printf("❌ Create failed: %v\n", result.Error)
		return
	}
	fmt.Printf("✅ Created user: ID=%d, Rows affected=%d\n", user.ID, result.RowsAffected)

	// Create multiple records
	users := []User{
		{Username: "alice", Email: "alice@example.com"},
		{Username: "bob", Email: "bob@example.com"},
		{Username: "charlie", Email: "charlie@example.com"},
	}
	db.Create(&users)
	fmt.Printf("✅ Created %d users in batch\n", len(users))
}

func queryRecords(db *gorm.DB) {
	// Find by primary key
	var user User
	db.First(&user, 1) // Find user with ID = 1
	fmt.Printf("✅ Found user by ID: %s\n", user.Username)

	// Find by condition
	db.Where("username = ?", "alice").First(&user)
	fmt.Printf("✅ Found user by username: %s (%s)\n", user.Username, user.Email)

	// Find all
	var users []User
	db.Find(&users)
	fmt.Printf("✅ Found %d total users\n", len(users))

	// Find with conditions
	db.Where("created_at > ?", time.Now().Add(-24*time.Hour)).Find(&users)
	fmt.Printf("✅ Found %d users created in last 24h\n", len(users))
}

func updateRecords(db *gorm.DB) {
	// Update single column
	db.Model(&User{}).Where("username = ?", "john_doe").Update("email", "john.doe@example.com")
	fmt.Println("✅ Updated single column")

	// Update multiple columns
	db.Model(&User{}).Where("username = ?", "alice").Updates(User{
		Email: "alice.updated@example.com",
	})
	fmt.Println("✅ Updated multiple columns")

	// Update with map
	db.Model(&User{}).Where("username = ?", "bob").Updates(map[string]interface{}{
		"email": "bob.updated@example.com",
	})
	fmt.Println("✅ Updated with map")
}

func deleteRecords(db *gorm.DB) {
	// Soft delete (if DeletedAt field exists)
	db.Delete(&User{}, 1)
	fmt.Println("✅ Soft deleted user with ID 1")

	// Delete with condition
	db.Where("username = ?", "charlie").Delete(&User{})
	fmt.Println("✅ Deleted user by condition")

	// Permanent delete
	db.Unscoped().Where("username = ?", "test").Delete(&User{})
	fmt.Println("✅ Permanently deleted user")
}

func associations(db *gorm.DB) {
	// Create user with profile
	user := User{
		Username: "jane_doe",
		Email:    "jane@example.com",
		Profile: Profile{
			FirstName: "Jane",
			LastName:  "Doe",
			Bio:       "Software Engineer",
		},
	}
	db.Create(&user)
	fmt.Println("✅ Created user with profile (Has One)")

	// Create post with tags
	post := Post{
		Title:    "My First Post",
		Content:  "Hello, World!",
		AuthorID: user.ID,
		Tags: []Tag{
			{Name: "golang"},
			{Name: "tutorial"},
		},
	}
	db.Create(&post)
	fmt.Println("✅ Created post with tags (Many To Many)")
}

func preloading(db *gorm.DB) {
	// Preload associations
	var users []User
	db.Preload("Profile").Preload("Posts").Find(&users)
	fmt.Printf("✅ Loaded %d users with profiles and posts\n", len(users))

	// Nested preloading
	var posts []Post
	db.Preload("Author.Profile").Preload("Tags").Find(&posts)
	fmt.Printf("✅ Loaded %d posts with authors and tags\n", len(posts))
}

func transactions(db *gorm.DB) {
	// Manual transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// Create user
		user := User{Username: "tx_user", Email: "tx@example.com"}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// Create profile
		profile := Profile{UserID: user.ID, FirstName: "TX", LastName: "User"}
		if err := tx.Create(&profile).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		fmt.Printf("❌ Transaction failed: %v\n", err)
	} else {
		fmt.Println("✅ Transaction completed successfully")
	}
}

func hooks(db *gorm.DB) {
	// Hooks are defined in model methods
	// Example: BeforeCreate, AfterCreate, BeforeUpdate, etc.
	fmt.Println("✅ Hooks can be defined in model methods")
	fmt.Println("   Example: func (u *User) BeforeCreate(tx *gorm.DB) error { ... }")
}

func advancedQueries(db *gorm.DB) {
	var users []User

	// Select specific fields
	db.Select("username", "email").Find(&users)
	fmt.Printf("✅ Selected specific fields for %d users\n", len(users))

	// Order
	db.Order("created_at desc").Find(&users)
	fmt.Println("✅ Ordered by created_at desc")

	// Limit and Offset
	db.Limit(5).Offset(0).Find(&users)
	fmt.Println("✅ Limited to 5 users with offset 0")

	// Group and Having
	type Result struct {
		Date  time.Time
		Count int
	}
	var results []Result
	db.Model(&User{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Group("DATE(created_at)").
		Having("COUNT(*) > ?", 0).
		Scan(&results)
	fmt.Printf("✅ Grouped by date: %d results\n", len(results))

	// Raw SQL
	db.Raw("SELECT * FROM users WHERE username LIKE ?", "%doe%").Scan(&users)
	fmt.Printf("✅ Raw SQL query: found %d users\n", len(users))
}


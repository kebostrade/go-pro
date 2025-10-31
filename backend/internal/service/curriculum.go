// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides functionality for the GO-PRO Learning Platform.
package service

import (
	"context"
	"time"

	"go-pro-backend/internal/cache"
	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/errors"
	"go-pro-backend/internal/messaging"
	"go-pro-backend/pkg/logger"
	"go-pro-backend/pkg/validator"
)

// ErrLessonNotFound is returned when a lesson is not found.
var ErrLessonNotFound = errors.NewNotFoundError("lesson not found")

// curriculumService implements the CurriculumService interface.
type curriculumService struct {
	logger    logger.Logger
	validator validator.Validator
	cache     cache.CacheManager
	messaging *messaging.Service
}

// NewCurriculumService creates a new curriculum service.
func NewCurriculumService(config *Config) CurriculumService {
	return &curriculumService{
		logger:    config.Logger,
		validator: config.Validator,
		cache:     config.Cache,
		messaging: config.Messaging,
	}
}

// GetCurriculum returns the complete curriculum structure.
func (s *curriculumService) GetCurriculum(ctx context.Context) (*domain.Curriculum, error) {
	s.logger.Info(ctx, "Getting curriculum")

	// In a real implementation, this would come from a database or file system.
	// For now, we'll return the curriculum based on the syllabus.md structure.
	curriculum := &domain.Curriculum{
		ID:          "go-pro-curriculum",
		Title:       "GO-PRO: Complete Go Programming Mastery",
		Description: "Master Go programming from basics to advanced production systems. Learn Go's syntax, concurrency patterns, web development, testing, performance optimization, security, and best practices through hands-on exercises and real-world projects.",
		Duration:    "16 weeks",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Phases: []domain.CurriculumPhase{
			{
				ID:          "foundations",
				Title:       "Foundations",
				Description: "Master Go basics and core concepts",
				Weeks:       "Weeks 1-2",
				Icon:        "zap",
				Color:       "text-blue-500",
				Order:       1,
				Progress:    80, // Updated progress
				Lessons: []domain.CurriculumLesson{
					{
						ID:          1,
						Title:       "Go Syntax and Basic Types",
						Description: "Go installation, basic syntax, primitive types, constants and iota",
						Duration:    "3-4 hours",
						Exercises:   5,
						Difficulty:  domain.DifficultyBeginner,
						Completed:   true,
						Locked:      false,
						Order:       1,
					},
					{
						ID:          2,
						Title:       "Variables, Constants, and Functions",
						Description: "Variable declarations, scope, function definitions, multiple return values",
						Duration:    "4-5 hours",
						Exercises:   6,
						Difficulty:  domain.DifficultyBeginner,
						Completed:   true,
						Locked:      false,
						Order:       2,
					},
					{
						ID:          3,
						Title:       "Control Structures and Loops",
						Description: "if/else, switch statements, for loops, defer statements",
						Duration:    "3-4 hours",
						Exercises:   7,
						Difficulty:  domain.DifficultyBeginner,
						Completed:   true,
						Locked:      false,
						Order:       3,
					},
					{
						ID:          4,
						Title:       "Arrays, Slices, and Maps",
						Description: "Data structures, manipulation, memory considerations",
						Duration:    "5-6 hours",
						Exercises:   8,
						Difficulty:  domain.DifficultyBeginner,
						Completed:   true,
						Locked:      false,
						Order:       4,
					},
					{
						ID:          5,
						Title:       "Pointers and Memory Management",
						Description: "Pointer basics, memory allocation, garbage collection",
						Duration:    "4-5 hours",
						Exercises:   8,
						Difficulty:  domain.DifficultyBeginner,
						Completed:   true,
						Locked:      false,
						Order:       5,
					},
				},
			},
			{
				ID:          "intermediate",
				Title:       "Intermediate",
				Description: "Object-oriented concepts and concurrency",
				Weeks:       "Weeks 3-5",
				Icon:        "globe",
				Color:       "text-green-500",
				Order:       2,
				Progress:    60,
				Lessons: []domain.CurriculumLesson{
					{
						ID:          6,
						Title:       "Structs and Methods",
						Description: "Struct definition, methods, receivers, method sets",
						Duration:    "5-6 hours",
						Exercises:   8,
						Difficulty:  domain.DifficultyIntermediate,
						Completed:   true,
						Locked:      false,
						Order:       6,
					},
					{
						ID:          7,
						Title:       "Interfaces and Polymorphism",
						Description: "Interface definition, type assertions, composition",
						Duration:    "6-7 hours",
						Exercises:   9,
						Difficulty:  domain.DifficultyIntermediate,
						Completed:   false,
						Locked:      false,
						Order:       7,
					},
					{
						ID:          8,
						Title:       "Error Handling Patterns",
						Description: "Custom errors, wrapping, panic/recover, best practices",
						Duration:    "4-5 hours",
						Exercises:   7,
						Difficulty:  domain.DifficultyIntermediate,
						Completed:   false,
						Locked:      true,
						Order:       8,
					},
					{
						ID:          9,
						Title:       "Goroutines and Channels",
						Description: "Concurrency, channel operations, select statements",
						Duration:    "7-8 hours",
						Exercises:   10,
						Difficulty:  domain.DifficultyIntermediate,
						Completed:   false,
						Locked:      true,
						Order:       9,
					},
					{
						ID:          10,
						Title:       "Packages and Modules",
						Description: "Package organization, Go modules, dependency management",
						Duration:    "5-6 hours",
						Exercises:   6,
						Difficulty:  domain.DifficultyIntermediate,
						Completed:   false,
						Locked:      true,
						Order:       10,
					},
				},
			},
			{
				ID:          "advanced",
				Title:       "Advanced",
				Description: "Production-ready development skills",
				Weeks:       "Weeks 6-8",
				Icon:        "trending-up",
				Color:       "text-purple-500",
				Order:       3,
				Progress:    0,
				Lessons: []domain.CurriculumLesson{
					{
						ID:          11,
						Title:       "Advanced Concurrency Patterns",
						Description: "Worker pools, pipelines, context package, sync primitives",
						Duration:    "8-9 hours",
						Exercises:   12,
						Difficulty:  domain.DifficultyAdvanced,
						Completed:   false,
						Locked:      true,
						Order:       11,
					},
					{
						ID:          12,
						Title:       "Testing and Benchmarking",
						Description: "Unit testing, table-driven tests, benchmarking, profiling",
						Duration:    "6-7 hours",
						Exercises:   8,
						Difficulty:  domain.DifficultyAdvanced,
						Completed:   false,
						Locked:      true,
						Order:       12,
					},
					{
						ID:          13,
						Title:       "HTTP Servers and REST APIs",
						Description: "HTTP server basics, routing, middleware, authentication",
						Duration:    "8-9 hours",
						Exercises:   10,
						Difficulty:  domain.DifficultyAdvanced,
						Completed:   false,
						Locked:      true,
						Order:       13,
					},
					{
						ID:          14,
						Title:       "Database Integration",
						Description: "Database/sql package, connection pooling, transactions",
						Duration:    "7-8 hours",
						Exercises:   9,
						Difficulty:  domain.DifficultyAdvanced,
						Completed:   false,
						Locked:      true,
						Order:       14,
					},
					{
						ID:          15,
						Title:       "Microservices Architecture",
						Description: "Design principles, service communication, monitoring",
						Duration:    "9-10 hours",
						Exercises:   11,
						Difficulty:  domain.DifficultyAdvanced,
						Completed:   false,
						Locked:      false,
						Order:       15,
					},
				},
			},
			{
				ID:          "expert",
				Title:       "Expert",
				Description: "Advanced patterns and production systems",
				Weeks:       "Weeks 9-12",
				Icon:        "award",
				Color:       "text-orange-500",
				Order:       4,
				Progress:    0,
				Lessons: []domain.CurriculumLesson{
					{
						ID:          16,
						Title:       "Performance Optimization and Profiling",
						Description: "Memory optimization, CPU profiling, benchmarking techniques",
						Duration:    "8-10 hours",
						Exercises:   10,
						Difficulty:  domain.DifficultyExpert,
						Completed:   false,
						Locked:      false,
						Order:       16,
					},
					{
						ID:          17,
						Title:       "Security Best Practices",
						Description: "Authentication, encryption, vulnerability prevention",
						Duration:    "7-9 hours",
						Exercises:   9,
						Difficulty:  domain.DifficultyExpert,
						Completed:   false,
						Locked:      false,
						Order:       17,
					},
					{
						ID:          18,
						Title:       "Deployment and DevOps",
						Description: "Docker, CI/CD, cloud deployment, monitoring",
						Duration:    "9-11 hours",
						Exercises:   12,
						Difficulty:  domain.DifficultyExpert,
						Completed:   false,
						Locked:      false,
						Order:       18,
					},
					{
						ID:          19,
						Title:       "Advanced Design Patterns",
						Description: "Functional programming, generics, architectural patterns",
						Duration:    "8-10 hours",
						Exercises:   11,
						Difficulty:  domain.DifficultyExpert,
						Completed:   false,
						Locked:      false,
						Order:       19,
					},
					{
						ID:          20,
						Title:       "Building Production Systems",
						Description: "Complete system design, observability, scalability",
						Duration:    "10-12 hours",
						Exercises:   15,
						Difficulty:  domain.DifficultyExpert,
						Completed:   false,
						Locked:      false,
						Order:       20,
					},
				},
			},
		},
		Projects: []domain.Project{
			{
				ID:          "cli-task-manager",
				Title:       "CLI Task Manager",
				Description: "Command-line application with file persistence",
				Duration:    "1 week",
				Difficulty:  domain.DifficultyIntermediate,
				Skills:      []string{"File I/O", "JSON handling", "CLI design"},
				Points:      100,
				Completed:   false,
				Locked:      false,
				Order:       1,
			},
			{
				ID:          "rest-api-server",
				Title:       "REST API with Database",
				Description: "Full REST API with PostgreSQL integration",
				Duration:    "1.5 weeks",
				Difficulty:  domain.DifficultyAdvanced,
				Skills:      []string{"HTTP servers", "Database integration", "Testing"},
				Points:      150,
				Completed:   false,
				Locked:      false,
				Order:       2,
			},
			{
				ID:          "realtime-chat",
				Title:       "Real-time Chat Server",
				Description: "WebSocket-based chat with concurrent users",
				Duration:    "1.5 weeks",
				Difficulty:  domain.DifficultyAdvanced,
				Skills:      []string{"Concurrency", "WebSockets", "Real-time communication"},
				Points:      200,
				Completed:   false,
				Locked:      false,
				Order:       3,
			},
			{
				ID:          "microservices-system",
				Title:       "Microservices System",
				Description: "Multi-service system with API gateway",
				Duration:    "2 weeks",
				Difficulty:  domain.DifficultyAdvanced,
				Skills:      []string{"Microservices", "Service mesh", "Monitoring"},
				Points:      250,
				Completed:   false,
				Locked:      false,
				Order:       4,
			},
			{
				ID:          "distributed-cache",
				Title:       "Distributed Cache System",
				Description: "High-performance distributed cache similar to Redis",
				Duration:    "3 weeks",
				Difficulty:  domain.DifficultyExpert,
				Skills:      []string{"Distributed Systems", "Networking", "Performance"},
				Points:      300,
				Completed:   false,
				Locked:      false,
				Order:       5,
			},
			{
				ID:          "event-driven-system",
				Title:       "Event-Driven Architecture",
				Description: "Complete event-driven system with CQRS and event sourcing",
				Duration:    "3 weeks",
				Difficulty:  domain.DifficultyExpert,
				Skills:      []string{"Event Sourcing", "CQRS", "Message Queues"},
				Points:      350,
				Completed:   false,
				Locked:      false,
				Order:       6,
			},
			{
				ID:          "observability-platform",
				Title:       "Monitoring & Observability Platform",
				Description: "Comprehensive monitoring platform for distributed systems",
				Duration:    "3 weeks",
				Difficulty:  domain.DifficultyExpert,
				Skills:      []string{"Observability", "Time Series", "Real-time Analytics"},
				Points:      400,
				Completed:   false,
				Locked:      false,
				Order:       7,
			},
		},
	}

	s.logger.Info(ctx, "Curriculum retrieved successfully")

	return curriculum, nil
}

// GetLessonDetail returns detailed information about a specific lesson.
func (s *curriculumService) GetLessonDetail(ctx context.Context, lessonID int) (*domain.LessonDetail, error) {
	s.logger.Info(ctx, "Getting lesson detail", "lesson_id", lessonID)

	// Mock lesson data - in a real implementation, this would come from a database or file system.
	lessonData := map[int]*domain.LessonDetail{
		1: {
			ID:          1,
			Title:       "Go Syntax and Basic Types",
			Description: "Learn the fundamental syntax of Go and work with basic data types including integers, floats, strings, and booleans.",
			Duration:    "3-4 hours",
			Difficulty:  domain.DifficultyBeginner,
			Phase:       "Foundations",
			Objectives: []string{
				"Set up a Go development environment",
				"Understand Go's basic syntax and program structure",
				"Work with primitive data types (int, float, string, bool)",
				"Declare and use constants",
				"Perform type conversions",
				"Use the iota identifier for enumerated constants",
			},
			Theory: `# Go Program Structure

Every Go program starts with a package declaration, followed by imports, and then the program code:

` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
` + "```" + `

## Basic Types

Go has several built-in basic types:

### Numeric Types
- **Integers**: int, int8, int16, int32, int64
- **Unsigned integers**: uint, uint8, uint16, uint32, uint64
- **Floating point**: float32, float64
- **Complex numbers**: complex64, complex128

### Other Types
- **Boolean**: bool (true or false)
- **String**: string (UTF-8 encoded)
- **Byte**: byte (alias for uint8)
- **Rune**: rune (alias for int32, represents Unicode code points)`,
			CodeExample: `package main

import "fmt"

func main() {
    // Basic variable declarations.
    var name string = "Go Programming"
    var version float64 = 1.21
    var isAwesome bool = true
    
    // Short variable declaration.
    year := 2024
    
    // Constants.
    const MaxUsers = 1000
    
    // Type conversion.
    var x int = 42
    var y float64 = float64(x)
    
    // Print values.
    fmt.Printf("Language: %s\n", name)
    fmt.Printf("Version: %.2f\n", version)
    fmt.Printf("Year: %d\n", year)
    fmt.Printf("Is Awesome: %t\n", isAwesome)
    fmt.Printf("Max Users: %d\n", MaxUsers)
    fmt.Printf("Converted: %.1f\n", y)
}`,
			Solution: `package main

import "fmt"

func main() {
    // Basic variable declarations.
    var name string = "Go Programming"
    var version float64 = 1.21
    var isAwesome bool = true
    
    // Short variable declaration.
    year := 2024
    
    // Constants.
    const MaxUsers = 1000
    
    // Type conversion.
    var x int = 42
    var y float64 = float64(x)
    
    // Print values.
    fmt.Printf("Language: %s\n", name)
    fmt.Printf("Version: %.2f\n", version)
    fmt.Printf("Year: %d\n", year)
    fmt.Printf("Is Awesome: %t\n", isAwesome)
    fmt.Printf("Max Users: %d\n", MaxUsers)
    fmt.Printf("Converted: %.1f\n", y)
    
    // Additional examples.
    
    // Multiple variable declaration.
    var (
        firstName = "John"
        lastName  = "Doe"
        age       = 30
    )
    
    fmt.Printf("Full Name: %s %s, Age: %d\n", firstName, lastName, age)
    
    // Enumerated constants.
    const (
        Red = iota
        Green
        Blue
    )
    
    fmt.Printf("Colors: Red=%d, Green=%d, Blue=%d\n", Red, Green, Blue)
}`,
			Exercises: []domain.LessonExercise{
				{
					ID:          "basic-variables",
					Title:       "Variable Declaration Practice",
					Description: "Practice declaring variables of different types and using type conversions.",
					Requirements: []string{
						"Declare a string variable for your name",
						"Declare an integer variable for your age",
						"Declare a boolean variable for whether you like programming",
						"Use short variable declaration for the current year",
						"Convert an integer to float64 and print both values",
					},
					InitialCode: `package main

import "fmt"

func main() {
    // TODO: Declare your variables here.
    
    // TODO: Print the values.
    
}`,
					Solution: `package main

import "fmt"

func main() {
    // Variable declarations.
    var name string = "Alice"
    var age int = 25
    var likesProgramming bool = true
    currentYear := 2024
    
    // Type conversion.
    var score int = 95
    var percentage float64 = float64(score)
    
    // Print values.
    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Likes Programming: %t\n", likesProgramming)
    fmt.Printf("Current Year: %d\n", currentYear)
    fmt.Printf("Score: %d, Percentage: %.1f%%\n", score, percentage)
}`,
				},
			},
			NextLessonID: func() *int { i := 2; return &i }(),
		},
	}

	lesson, exists := lessonData[lessonID]
	if !exists {
		s.logger.Warn(ctx, "Lesson not found", "lesson_id", lessonID)
		return nil, ErrLessonNotFound
	}

	s.logger.Info(ctx, "Lesson detail retrieved successfully", "lesson_id", lessonID)

	return lesson, nil
}

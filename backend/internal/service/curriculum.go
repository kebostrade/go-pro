// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package service provides functionality for the GO-PRO Learning Platform.
package service

import (
	"context"
	"fmt"
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
		Description: "Master Go programming from basics to advanced production systems. " +
			"Learn Go's syntax, concurrency patterns, web development, testing, " +
			"performance optimization, security, and best practices through hands-on exercises and real-world projects.",
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
	lessonData := s.generateLessonMockData()

	lesson, exists := lessonData[lessonID]
	if !exists {
		s.logger.Warn(ctx, "Lesson not found", "lesson_id", lessonID)
		return nil, ErrLessonNotFound
	}

	s.logger.Info(ctx, "Lesson detail retrieved successfully", "lesson_id", lessonID)

	return lesson, nil
}

// generateLessonMockData generates mock lesson data for all 20 lessons.
func (s *curriculumService) generateLessonMockData() map[int]*domain.LessonDetail {
	return map[int]*domain.LessonDetail{
		1:  s.getLessonData1(),
		2:  s.getLessonData2(),
		3:  s.getLessonData3(),
		4:  s.getLessonData4(),
		5:  s.getLessonData5(),
		6:  s.getLessonData6(),
		7:  s.getLessonData7(),
		8:  s.getLessonData8(),
		9:  s.getLessonData9(),
		10: s.getLessonData10(),
		11: s.getLessonData11(),
		12: s.getLessonData12(),
		13: s.getLessonData13(),
		14: s.getLessonData14(),
		15: s.getLessonData15(),
		16: s.getLessonData16(),
		17: s.getLessonData17(),
		18: s.getLessonData18(),
		19: s.getLessonData19(),
		20: s.getLessonData20(),
	}
}

func (s *curriculumService) getLessonData1() *domain.LessonDetail {
	return &domain.LessonDetail{
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

Every Go program starts with a package declaration, followed by imports, and then the program code.

## Basic Types

Go has several built-in basic types including integers, floats, strings, booleans, and more.`,
		CodeExample: `package main

import "fmt"

func main() {
    var name string = "Go Programming"
    year := 2024
    const MaxUsers = 1000

    fmt.Printf("Language: %s, Year: %d\n", name, year)
}`,
		Solution: `package main

import "fmt"

func main() {
    var name string = "Go Programming"
    var version float64 = 1.21
    year := 2024
    const MaxUsers = 1000

    fmt.Printf("Language: %s %.2f, Year: %d, Max Users: %d\n", name, version, year, MaxUsers)
}`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "basic-variables",
				Title:       "Variable Declaration Practice",
				Description: "Practice declaring variables of different types and using type conversions.",
				Requirements: []string{
					"Declare variables of different types",
					"Use type conversions",
					"Print values to console",
				},
				InitialCode: `package main

import "fmt"

func main() {
    // TODO: Declare your variables here
}`,
				Solution: `package main

import "fmt"

func main() {
    var name string = "Alice"
    var age int = 25
    fmt.Printf("Name: %s, Age: %d\n", name, age)
}`,
			},
		},
		NextLessonID: func() *int { i := 2; return &i }(),
	}
}

func (s *curriculumService) getLessonData2() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          2,
		Title:       "Variables, Constants, and Functions",
		Description: "Master variable declarations, scope, function definitions, and multiple return values in Go.",
		Duration:    "4-5 hours",
		Difficulty:  domain.DifficultyBeginner,
		Phase:       "Foundations",
		Objectives: []string{
			"Understand variable declaration methods",
			"Learn about variable scope",
			"Create and use functions",
			"Work with multiple return values",
			"Understand named return values",
		},
		Theory: `# Functions in Go

Functions are fundamental building blocks in Go. They can return multiple values and support various declaration styles.`,
		CodeExample: `package main

import "fmt"

func add(a, b int) int {
    return a + b
}

func main() {
    result := add(5, 3)
    fmt.Println("Result:", result)
}`,
		Solution: `package main

import "fmt"

func add(a, b int) int {
    return a + b
}

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

func main() {
    sum := add(5, 3)
    quotient, _ := divide(10, 2)
    fmt.Printf("Sum: %d, Quotient: %.2f\n", sum, quotient)
}`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "function-practice",
				Title:       "Function Practice",
				Description: "Create functions with multiple return values.",
				Requirements: []string{
					"Create a function that returns multiple values",
					"Handle errors properly",
				},
				InitialCode: `package main

func main() {
    // TODO: Implement functions
}`,
				Solution: `package main

import "fmt"

func calculate(a, b int) (int, int) {
    return a + b, a * b
}

func main() {
    sum, product := calculate(5, 3)
    fmt.Printf("Sum: %d, Product: %d\n", sum, product)
}`,
			},
		},
		NextLessonID: func() *int { i := 3; return &i }(),
		PrevLessonID: func() *int { i := 1; return &i }(),
	}
}

func (s *curriculumService) getLessonData3() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          3,
		Title:       "Control Structures and Loops",
		Description: "Learn about if/else statements, switch statements, for loops, and defer in Go.",
		Duration:    "3-4 hours",
		Difficulty:  domain.DifficultyBeginner,
		Phase:       "Foundations",
		Objectives: []string{
			"Use if/else statements effectively",
			"Master switch statements",
			"Work with for loops",
			"Understand the defer statement",
		},
		Theory: `# Control Flow in Go

Go provides standard control structures like if/else, switch, and for loops with a clean syntax.`,
		CodeExample: `package main

import "fmt"

func main() {
    for i := 0; i < 5; i++ {
        if i%2 == 0 {
            fmt.Println(i, "is even")
        } else {
            fmt.Println(i, "is odd")
        }
    }
}`,
		Solution: `package main

import "fmt"

func main() {
    for i := 0; i < 10; i++ {
        switch {
        case i%2 == 0:
            fmt.Println(i, "is even")
        default:
            fmt.Println(i, "is odd")
        }
    }
}`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "control-flow",
				Title:       "Control Flow Practice",
				Description: "Practice using loops and conditionals.",
				Requirements: []string{
					"Use for loops",
					"Implement conditional logic",
				},
				InitialCode: `package main

func main() {
    // TODO: Implement control flow
}`,
				Solution: `package main

import "fmt"

func main() {
    for i := 1; i <= 10; i++ {
        if i%3 == 0 && i%5 == 0 {
            fmt.Println("FizzBuzz")
        } else if i%3 == 0 {
            fmt.Println("Fizz")
        } else if i%5 == 0 {
            fmt.Println("Buzz")
        } else {
            fmt.Println(i)
        }
    }
}`,
			},
		},
		NextLessonID: func() *int { i := 4; return &i }(),
		PrevLessonID: func() *int { i := 2; return &i }(),
	}
}

func (s *curriculumService) getLessonData4() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          4,
		Title:       "Arrays, Slices, and Maps",
		Description: "Master Go's fundamental data structures including arrays, slices, and maps.",
		Duration:    "5-6 hours",
		Difficulty:  domain.DifficultyBeginner,
		Phase:       "Foundations",
		Objectives: []string{
			"Understand arrays and their limitations",
			"Master slice operations",
			"Work with maps effectively",
			"Learn about make and append functions",
		},
		Theory: `# Data Structures in Go

Go provides powerful built-in data structures including slices and maps for efficient data management.`,
		CodeExample: `package main

import "fmt"

func main() {
    slice := []int{1, 2, 3, 4, 5}
    numbers := make(map[string]int)
    numbers["one"] = 1
    numbers["two"] = 2

    fmt.Println("Slice:", slice)
    fmt.Println("Map:", numbers)
}`,
		Solution: `package main

import "fmt"

func main() {
    slice := make([]int, 0, 10)
    for i := 0; i < 5; i++ {
        slice = append(slice, i)
    }

    m := map[string]int{"a": 1, "b": 2, "c": 3}
    for k, v := range m {
        fmt.Printf("%s: %d\n", k, v)
    }
}`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "data-structures",
				Title:       "Data Structures Practice",
				Description: "Work with slices and maps.",
				Requirements: []string{
					"Create and manipulate slices",
					"Use maps for key-value storage",
				},
				InitialCode: `package main

func main() {
    // TODO: Work with slices and maps
}`,
				Solution: `package main

import "fmt"

func main() {
    nums := []int{1, 2, 3}
    nums = append(nums, 4, 5)

    dict := make(map[string]string)
    dict["hello"] = "world"
    fmt.Println(nums, dict)
}`,
			},
		},
		NextLessonID: func() *int { i := 5; return &i }(),
		PrevLessonID: func() *int { i := 3; return &i }(),
	}
}

func (s *curriculumService) getLessonData5() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          5,
		Title:       "Pointers and Memory Management",
		Description: "Learn about pointers, memory allocation, and Go's garbage collection.",
		Duration:    "4-5 hours",
		Difficulty:  domain.DifficultyBeginner,
		Phase:       "Foundations",
		Objectives: []string{
			"Understand pointer basics",
			"Learn about memory allocation",
			"Work with new and make",
			"Understand garbage collection",
		},
		Theory: `# Pointers in Go

Pointers allow you to pass references to values and records within your program.`,
		CodeExample: `package main

import "fmt"

func increment(x *int) {
    *x++
}

func main() {
    num := 5
    increment(&num)
    fmt.Println(num)
}`,
		Solution: `package main

import "fmt"

func swap(a, b *int) {
    *a, *b = *b, *a
}

func main() {
    x, y := 10, 20
    swap(&x, &y)
    fmt.Printf("x=%d, y=%d\n", x, y)
}`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "pointers",
				Title:       "Pointer Practice",
				Description: "Work with pointers and memory.",
				Requirements: []string{
					"Use pointers to modify values",
					"Understand pointer dereferencing",
				},
				InitialCode: `package main

func main() {
    // TODO: Practice with pointers
}`,
				Solution: `package main

import "fmt"

func double(x *int) {
    *x *= 2
}

func main() {
    num := 5
    double(&num)
    fmt.Println(num)
}`,
			},
		},
		NextLessonID: func() *int { i := 6; return &i }(),
		PrevLessonID: func() *int { i := 4; return &i }(),
	}
}

func (s *curriculumService) getLessonData6() *domain.LessonDetail {
	return &domain.LessonDetail{
		ID:          6,
		Title:       "Structs and Methods",
		Description: "Master struct definitions, methods, and receivers in Go.",
		Duration:    "5-6 hours",
		Difficulty:  domain.DifficultyIntermediate,
		Phase:       "Intermediate",
		Objectives: []string{
			"Define and use structs",
			"Create methods with receivers",
			"Understand value vs pointer receivers",
			"Work with embedded structs",
		},
		Theory: `# Structs and Methods

Structs are typed collections of fields useful for grouping data together to form records.`,
		CodeExample: `package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func (p Person) Greet() {
    fmt.Printf("Hello, I'm %s\n", p.Name)
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    p.Greet()
}`,
		Solution: `package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func (p *Person) Birthday() {
    p.Age++
}

func (p Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s, %d years old", p.Name, p.Age)
}

func main() {
    p := &Person{Name: "Bob", Age: 25}
    p.Birthday()
    fmt.Println(p.Greet())
}`,
		Exercises: []domain.LessonExercise{
			{
				ID:          "structs",
				Title:       "Struct Practice",
				Description: "Create structs with methods.",
				Requirements: []string{
					"Define a struct",
					"Add methods to the struct",
				},
				InitialCode: `package main

func main() {
    // TODO: Create structs
}`,
				Solution: `package main

import "fmt"

type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func main() {
    r := Rectangle{Width: 10, Height: 5}
    fmt.Println("Area:", r.Area())
}`,
			},
		},
		NextLessonID: func() *int { i := 7; return &i }(),
		PrevLessonID: func() *int { i := 5; return &i }(),
	}
}

// Comprehensive Lesson 7 is in curriculum_lessons_7_10.go
func (s *curriculumService) getLessonData7() *domain.LessonDetail {
	return s.getComprehensiveLessonData7()
}

// Comprehensive Lesson 8 is in curriculum_lessons_7_10.go
func (s *curriculumService) getLessonData8() *domain.LessonDetail {
	return s.getComprehensiveLessonData8()
}

// Comprehensive Lesson 9 is in curriculum_lessons_7_10.go
func (s *curriculumService) getLessonData9() *domain.LessonDetail {
	return s.getComprehensiveLessonData9()
}

// Comprehensive Lesson 10 is in curriculum_lessons_7_10.go
func (s *curriculumService) getLessonData10() *domain.LessonDetail {
	return s.getComprehensiveLessonData10()
}

// Lessons 11-20 with simpler implementations
func (s *curriculumService) getLessonData11() *domain.LessonDetail {
	return s.createGenericLesson(11, "Advanced Concurrency Patterns",
		"Master worker pools, pipelines, context package, and sync primitives.",
		"8-9 hours", domain.DifficultyAdvanced, "Advanced", 12, 10)
}

func (s *curriculumService) getLessonData12() *domain.LessonDetail {
	return s.createGenericLesson(12, "Testing and Benchmarking",
		"Learn unit testing, table-driven tests, benchmarking, and profiling.",
		"6-7 hours", domain.DifficultyAdvanced, "Advanced", 13, 11)
}

func (s *curriculumService) getLessonData13() *domain.LessonDetail {
	return s.createGenericLesson(13, "HTTP Servers and REST APIs", "Build HTTP servers with routing, middleware, and authentication.", "8-9 hours", domain.DifficultyAdvanced, "Advanced", 14, 12)
}

func (s *curriculumService) getLessonData14() *domain.LessonDetail {
	return s.createGenericLesson(14, "Database Integration", "Work with database/sql package, connection pooling, and transactions.", "7-8 hours", domain.DifficultyAdvanced, "Advanced", 15, 13)
}

func (s *curriculumService) getLessonData15() *domain.LessonDetail {
	return s.createGenericLesson(15, "Microservices Architecture", "Design principles, service communication, and monitoring.", "9-10 hours", domain.DifficultyAdvanced, "Advanced", 16, 14)
}

func (s *curriculumService) getLessonData16() *domain.LessonDetail {
	return s.createGenericLesson(16, "Performance Optimization and Profiling", "Memory optimization, CPU profiling, and benchmarking techniques.", "8-10 hours", domain.DifficultyAdvanced, "Expert", 17, 15)
}

func (s *curriculumService) getLessonData17() *domain.LessonDetail {
	return s.createGenericLesson(17, "Security Best Practices", "Authentication, encryption, and vulnerability prevention.", "7-9 hours", domain.DifficultyAdvanced, "Expert", 18, 16)
}

func (s *curriculumService) getLessonData18() *domain.LessonDetail {
	return s.createGenericLesson(18, "Deployment and DevOps", "Docker, CI/CD, cloud deployment, and monitoring.", "9-11 hours", domain.DifficultyAdvanced, "Expert", 19, 17)
}

func (s *curriculumService) getLessonData19() *domain.LessonDetail {
	return s.createGenericLesson(19, "Advanced Design Patterns", "Functional programming, generics, and architectural patterns.", "8-10 hours", domain.DifficultyAdvanced, "Expert", 20, 18)
}

func (s *curriculumService) getLessonData20() *domain.LessonDetail {
	lesson := s.createGenericLesson(20, "Building Production Systems", "Complete system design, observability, and scalability.", "10-12 hours", domain.DifficultyAdvanced, "Expert", 0, 19)
	lesson.NextLessonID = nil
	return lesson
}

// Helper function to create generic lesson data
func (s *curriculumService) createGenericLesson(
	id int,
	title, description, duration string,
	difficulty domain.Difficulty,
	phase string,
	nextID, prevID int,
) *domain.LessonDetail {
	lesson := &domain.LessonDetail{
		ID:          id,
		Title:       title,
		Description: description,
		Duration:    duration,
		Difficulty:  difficulty,
		Phase:       phase,
		Objectives: []string{
			fmt.Sprintf("Master core concepts of %s", title),
			fmt.Sprintf("Apply best practices in %s", title),
			"Build practical projects using learned skills",
		},
		Theory: fmt.Sprintf("# %s\n\nThis lesson covers essential concepts and practical applications.", title),
		CodeExample: `package main

import "fmt"

func main() {
    fmt.Println("Welcome to this lesson!")
}`,
		Solution: `package main

import "fmt"

func main() {
    fmt.Println("Completed lesson solution")
}`,
		Exercises: []domain.LessonExercise{
			{
				ID:          fmt.Sprintf("exercise-%d", id),
				Title:       fmt.Sprintf("%s Practice", title),
				Description: fmt.Sprintf("Practice the concepts learned in %s", title),
				Requirements: []string{
					"Complete the implementation",
					"Test your code",
					"Review the solution",
				},
				InitialCode: `package main

func main() {
    // TODO: Implement the exercise
}`,
				Solution: `package main

import "fmt"

func main() {
    fmt.Println("Exercise completed!")
}`,
			},
		},
	}

	if nextID > 0 {
		lesson.NextLessonID = &nextID
	}
	if prevID > 0 {
		lesson.PrevLessonID = &prevID
	}

	return lesson
}

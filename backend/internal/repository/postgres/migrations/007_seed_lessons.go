// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package migrations

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"go-pro-backend/internal/repository/postgres"
)

// seedLessonsData seeds the database with all 20 lessons from the curriculum.
func seedLessonsData() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     7,
		Description: "Seed lessons data from curriculum",
		Up: func(tx *sql.Tx) error {
			// First, create the course record (required by FK)
			courseQuery := `
				INSERT INTO courses (id, title, slug, description, difficulty, duration_hours, prerequisites, learning_outcomes, is_published, created_at, updated_at)
				VALUES (
					'go-pro-course-2025',
					'GO-PRO: Complete Go Programming Mastery',
					'go-pro-complete',
					'Master Go programming from basics to advanced production systems. Learn Go''s syntax, concurrency patterns, web development, testing, performance optimization, security, and best practices through hands-on exercises and real-world projects.',
					'beginner',
					320,
					'{}',
					ARRAY['Master Go fundamentals', 'Build production systems', 'Advanced concurrency', 'Microservices architecture', 'Performance optimization'],
					true,
					CURRENT_TIMESTAMP,
					CURRENT_TIMESTAMP
				)
				ON CONFLICT (id) DO NOTHING
			`
			if _, err := tx.Exec(courseQuery); err != nil {
				return fmt.Errorf("failed to create course: %w", err)
			}

			// Define lesson data structure
			type lessonData struct {
				id          int
				title       string
				description string
				duration    int // minutes
				difficulty  string
				phase       string
				objectives  []string
				theory      string
				codeExample string
				solution    string
				exercises   []map[string]interface{}
				nextID      *int
				prevID      *int
			}

			// Helper function to create int pointer
			intPtr := func(i int) *int {
				return &i
			}

			// Define all 20 lessons
			lessons := []lessonData{
				// Lesson 1: Go Syntax and Basic Types
				{
					id:          1,
					title:       "Go Syntax and Basic Types",
					description: "Learn the fundamental syntax of Go and work with basic data types including integers, floats, strings, and booleans.",
					duration:    210, // 3-4 hours average
					difficulty:  "beginner",
					phase:       "Foundations",
					objectives: []string{
						"Set up a Go development environment",
						"Understand Go's basic syntax and program structure",
						"Work with primitive data types (int, float, string, bool)",
						"Declare and use constants",
						"Perform type conversions",
						"Use the iota identifier for enumerated constants",
					},
					theory: `# Go Program Structure

Every Go program starts with a package declaration, followed by imports, and then the program code.

## Basic Types

Go has several built-in basic types including integers, floats, strings, booleans, and more.`,
					codeExample: `package main

import "fmt"

func main() {
    var name string = "Go Programming"
    year := 2024
    const MaxUsers = 1000

    fmt.Printf("Language: %s, Year: %d\n", name, year)
}`,
					solution: `package main

import "fmt"

func main() {
    var name string = "Go Programming"
    var version float64 = 1.21
    year := 2024
    const MaxUsers = 1000

    fmt.Printf("Language: %s %.2f, Year: %d, Max Users: %d\n", name, version, year, MaxUsers)
}`,
					exercises: []map[string]interface{}{
						{
							"id":          "basic-variables",
							"title":       "Variable Declaration Practice",
							"description": "Practice declaring variables of different types and using type conversions.",
							"requirements": []string{
								"Declare variables of different types",
								"Use type conversions",
								"Print values to console",
							},
							"initial_code": `package main

import "fmt"

func main() {
    // TODO: Declare your variables here
}`,
							"solution": `package main

import "fmt"

func main() {
    var name string = "Alice"
    var age int = 25
    fmt.Printf("Name: %s, Age: %d\n", name, age)
}`,
						},
					},
					nextID: intPtr(2),
					prevID: nil,
				},
				// Lesson 2: Variables, Constants, and Functions
				{
					id:          2,
					title:       "Variables, Constants, and Functions",
					description: "Master variable declarations, scope, function definitions, and multiple return values in Go.",
					duration:    270, // 4-5 hours
					difficulty:  "beginner",
					phase:       "Foundations",
					objectives: []string{
						"Understand variable declaration methods",
						"Learn about variable scope",
						"Create and use functions",
						"Work with multiple return values",
						"Understand named return values",
					},
					theory: `# Functions in Go

Functions are fundamental building blocks in Go. They can return multiple values and support various declaration styles.`,
					codeExample: `package main

import "fmt"

func add(a, b int) int {
    return a + b
}

func main() {
    result := add(5, 3)
    fmt.Println("Result:", result)
}`,
					solution: `package main

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
					exercises: []map[string]interface{}{
						{
							"id":          "function-practice",
							"title":       "Function Practice",
							"description": "Create functions with multiple return values.",
							"requirements": []string{
								"Create a function that returns multiple values",
								"Handle errors properly",
							},
							"initial_code": `package main

func main() {
    // TODO: Implement functions
}`,
							"solution": `package main

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
					nextID: intPtr(3),
					prevID: intPtr(1),
				},
				// Lesson 3: Control Structures and Loops
				{
					id:          3,
					title:       "Control Structures and Loops",
					description: "Learn about if/else statements, switch statements, for loops, and defer in Go.",
					duration:    210, // 3-4 hours
					difficulty:  "beginner",
					phase:       "Foundations",
					objectives: []string{
						"Use if/else statements effectively",
						"Master switch statements",
						"Work with for loops",
						"Understand the defer statement",
					},
					theory: `# Control Flow in Go

Go provides standard control structures like if/else, switch, and for loops with a clean syntax.`,
					codeExample: `package main

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
					solution: `package main

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
					exercises: []map[string]interface{}{
						{
							"id":          "control-flow",
							"title":       "Control Flow Practice",
							"description": "Practice using loops and conditionals.",
							"requirements": []string{
								"Use for loops",
								"Implement conditional logic",
							},
							"initial_code": `package main

func main() {
    // TODO: Implement control flow
}`,
							"solution": `package main

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
					nextID: intPtr(4),
					prevID: intPtr(2),
				},
				// Lesson 4: Arrays, Slices, and Maps
				{
					id:          4,
					title:       "Arrays, Slices, and Maps",
					description: "Master Go's fundamental data structures including arrays, slices, and maps.",
					duration:    330, // 5-6 hours
					difficulty:  "beginner",
					phase:       "Foundations",
					objectives: []string{
						"Understand arrays and their limitations",
						"Master slice operations",
						"Work with maps effectively",
						"Learn about make and append functions",
					},
					theory: `# Data Structures in Go

Go provides powerful built-in data structures including slices and maps for efficient data management.`,
					codeExample: `package main

import "fmt"

func main() {
    slice := []int{1, 2, 3, 4, 5}
    numbers := make(map[string]int)
    numbers["one"] = 1
    numbers["two"] = 2

    fmt.Println("Slice:", slice)
    fmt.Println("Map:", numbers)
}`,
					solution: `package main

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
					exercises: []map[string]interface{}{
						{
							"id":          "data-structures",
							"title":       "Data Structures Practice",
							"description": "Work with slices and maps.",
							"requirements": []string{
								"Create and manipulate slices",
								"Use maps for key-value storage",
							},
							"initial_code": `package main

func main() {
    // TODO: Work with slices and maps
}`,
							"solution": `package main

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
					nextID: intPtr(5),
					prevID: intPtr(3),
				},
				// Lesson 5: Pointers and Memory Management
				{
					id:          5,
					title:       "Pointers and Memory Management",
					description: "Learn about pointers, memory allocation, and Go's garbage collection.",
					duration:    270, // 4-5 hours
					difficulty:  "beginner",
					phase:       "Foundations",
					objectives: []string{
						"Understand pointer basics",
						"Learn about memory allocation",
						"Work with new and make",
						"Understand garbage collection",
					},
					theory: `# Pointers in Go

Pointers allow you to pass references to values and records within your program.`,
					codeExample: `package main

import "fmt"

func increment(x *int) {
    *x++
}

func main() {
    num := 5
    increment(&num)
    fmt.Println(num)
}`,
					solution: `package main

import "fmt"

func swap(a, b *int) {
    *a, *b = *b, *a
}

func main() {
    x, y := 10, 20
    swap(&x, &y)
    fmt.Printf("x=%d, y=%d\n", x, y)
}`,
					exercises: []map[string]interface{}{
						{
							"id":          "pointers",
							"title":       "Pointer Practice",
							"description": "Work with pointers and memory.",
							"requirements": []string{
								"Use pointers to modify values",
								"Understand pointer dereferencing",
							},
							"initial_code": `package main

func main() {
    // TODO: Practice with pointers
}`,
							"solution": `package main

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
					nextID: intPtr(6),
					prevID: intPtr(4),
				},
				// Lessons 6-20: Simplified structure (extend with full data as needed)
				{
					id:          6,
					title:       "Structs and Methods",
					description: "Master struct definitions, methods, and receivers in Go.",
					duration:    330,
					difficulty:  "intermediate",
					phase:       "Intermediate",
					objectives:  []string{"Define and use structs", "Create methods with receivers", "Understand value vs pointer receivers"},
					theory:      "# Structs and Methods\n\nStructs are typed collections of fields useful for grouping data together to form records.",
					codeExample: "package main\n\nimport \"fmt\"\n\ntype Person struct {\n    Name string\n    Age  int\n}\n\nfunc (p Person) Greet() {\n    fmt.Printf(\"Hello, I'm %s\\n\", p.Name)\n}\n\nfunc main() {\n    p := Person{Name: \"Alice\", Age: 30}\n    p.Greet()\n}",
					solution:    "package main\n\nimport \"fmt\"\n\ntype Person struct {\n    Name string\n    Age  int\n}\n\nfunc (p *Person) Birthday() {\n    p.Age++\n}\n\nfunc (p Person) Greet() string {\n    return fmt.Sprintf(\"Hello, I'm %s, %d years old\", p.Name, p.Age)\n}\n\nfunc main() {\n    p := &Person{Name: \"Bob\", Age: 25}\n    p.Birthday()\n    fmt.Println(p.Greet())\n}",
					exercises:   []map[string]interface{}{{"id": "structs", "title": "Struct Practice", "description": "Create structs with methods.", "requirements": []string{"Define a struct", "Add methods"}, "initial_code": "package main\n\nfunc main() {\n    // TODO\n}", "solution": "package main\n\nimport \"fmt\"\n\ntype Rectangle struct {\n    Width, Height float64\n}\n\nfunc (r Rectangle) Area() float64 {\n    return r.Width * r.Height\n}\n\nfunc main() {\n    r := Rectangle{Width: 10, Height: 5}\n    fmt.Println(\"Area:\", r.Area())\n}"}},
					nextID:      intPtr(7),
					prevID:      intPtr(5),
				},
				// Continue with lessons 7-20 (simplified for brevity)
			}

			// Add remaining lessons (7-20) - simplified data
			for i := 7; i <= 20; i++ {
				lessonTitles := map[int]string{
					7:  "Interfaces and Polymorphism",
					8:  "Error Handling Patterns",
					9:  "Goroutines and Channels",
					10: "Packages and Modules",
					11: "Advanced Concurrency Patterns",
					12: "Testing and Benchmarking",
					13: "HTTP Servers and REST APIs",
					14: "Database Integration",
					15: "Microservices Architecture",
					16: "Performance Optimization and Profiling",
					17: "Security Best Practices",
					18: "Deployment and DevOps",
					19: "Advanced Design Patterns",
					20: "Building Production Systems",
				}

				phase := "Intermediate"
				difficulty := "intermediate"
				if i >= 11 {
					phase = "Advanced"
					difficulty = "advanced"
				}
				if i >= 16 {
					phase = "Expert"
				}

				var nextID, prevID *int
				if i < 20 {
					nextID = intPtr(i + 1)
				}
				prevID = intPtr(i - 1)

				lessons = append(lessons, lessonData{
					id:          i,
					title:       lessonTitles[i],
					description: fmt.Sprintf("Learn about %s in Go programming.", lessonTitles[i]),
					duration:    300,
					difficulty:  difficulty,
					phase:       phase,
					objectives:  []string{fmt.Sprintf("Master %s concepts", lessonTitles[i]), "Apply best practices"},
					theory:      fmt.Sprintf("# %s\n\nThis lesson covers essential concepts.", lessonTitles[i]),
					codeExample: "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Welcome to this lesson!\")\n}",
					solution:    "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Completed lesson solution\")\n}",
					exercises:   []map[string]interface{}{{"id": fmt.Sprintf("exercise-%d", i), "title": fmt.Sprintf("%s Practice", lessonTitles[i]), "description": "Practice the concepts", "requirements": []string{"Complete implementation"}, "initial_code": "package main\n\nfunc main() {\n    // TODO\n}", "solution": "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Exercise completed!\")\n}"}},
					nextID:      nextID,
					prevID:      prevID,
				})
			}

			// Helper functions
			exercisesToJSON := func(exercises []map[string]interface{}) (string, error) {
				data, err := json.Marshal(exercises)
				if err != nil {
					return "", err
				}
				return string(data), nil
			}

			stringsToJSON := func(strs []string) (string, error) {
				data, err := json.Marshal(strs)
				if err != nil {
					return "", err
				}
				return string(data), nil
			}

			// Insert all lessons
			for _, lesson := range lessons {
				objectivesJSON, err := stringsToJSON(lesson.objectives)
				if err != nil {
					return fmt.Errorf("failed to marshal objectives for lesson %d: %w", lesson.id, err)
				}

				exercisesJSON, err := exercisesToJSON(lesson.exercises)
				if err != nil {
					return fmt.Errorf("failed to marshal exercises for lesson %d: %w", lesson.id, err)
				}

				insertQuery := `
					INSERT INTO lessons (
						id, course_id, title, slug, content, order_index, duration_minutes,
						description, difficulty, phase, objectives, theory, code_example,
						solution, exercises, next_lesson_id, prev_lesson_id,
						is_published, created_at, updated_at
					) VALUES (
						$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb, $12, $13, $14, $15::jsonb, $16, $17, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
					)
					ON CONFLICT (id) DO UPDATE SET
						title = EXCLUDED.title,
						description = EXCLUDED.description,
						duration_minutes = EXCLUDED.duration_minutes,
						difficulty = EXCLUDED.difficulty,
						phase = EXCLUDED.phase,
						objectives = EXCLUDED.objectives,
						theory = EXCLUDED.theory,
						code_example = EXCLUDED.code_example,
						solution = EXCLUDED.solution,
						exercises = EXCLUDED.exercises,
						next_lesson_id = EXCLUDED.next_lesson_id,
						prev_lesson_id = EXCLUDED.prev_lesson_id,
						updated_at = CURRENT_TIMESTAMP
				`

				slug := fmt.Sprintf("lesson-%d", lesson.id)
				content := fmt.Sprintf("%s\n\n%s", lesson.theory, lesson.codeExample)

				_, err = tx.Exec(
					insertQuery,
					lesson.id,
					"go-pro-course-2025",
					lesson.title,
					slug,
					content,
					lesson.id, // order_index same as id
					lesson.duration,
					lesson.description,
					lesson.difficulty,
					lesson.phase,
					objectivesJSON,
					lesson.theory,
					lesson.codeExample,
					lesson.solution,
					exercisesJSON,
					lesson.nextID,
					lesson.prevID,
				)
				if err != nil {
					return fmt.Errorf("failed to insert lesson %d: %w", lesson.id, err)
				}
			}

			return nil
		},
		Down: func(tx *sql.Tx) error {
			// Delete all lessons for this course
			if _, err := tx.Exec("DELETE FROM lessons WHERE course_id = 'go-pro-course-2025'"); err != nil {
				return fmt.Errorf("failed to delete lessons: %w", err)
			}

			// Delete the course
			if _, err := tx.Exec("DELETE FROM courses WHERE id = 'go-pro-course-2025'"); err != nil {
				return fmt.Errorf("failed to delete course: %w", err)
			}

			return nil
		},
	}
}

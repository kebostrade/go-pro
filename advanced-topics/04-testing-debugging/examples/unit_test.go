//go:build unit
// +build unit

package testingexamples

// This file contains unit test examples
// Run with: go test -tags=unit -v ./...

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

// ========================================
// EXAMPLE 1: BASIC TEST STRUCTURE
// ========================================

func TestBasicAssertions(t *testing.T) {
	// t.Error / t.Errorf - marks test as failed but continues
	if 1+1 != 2 {
		t.Error("expected 1+1 to equal 2")
	}

	// t.Fatal / t.Fatalf - marks test as failed and stops immediately
	if 2*2 != 4 {
		t.Fatal("critical failure: 2*2 should equal 4")
	}

	// t.Log / t.Logf - logging (only shown with -v flag)
	t.Log("this is a log message")
	t.Logf("formatted log: %d + %d = %d", 1, 1, 2)
}

// ========================================
// EXAMPLE 2: TABLE-DRIVEN TESTS
// ========================================

// Add function to test
func Add(a, b int) int {
	return a + b
}

// Divide function that can error
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 2, 3, 5},
		{"negative numbers", -2, -3, -5},
		{"mixed numbers", -2, 3, 1},
		{"zeros", 0, 0, 0},
		{"large numbers", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name          string
		a, b          float64
		expected      float64
		expectError   bool
		errorContains string
	}{
		{"simple division", 10, 2, 5, false, ""},
		{"division with remainder", 7, 2, 3.5, false, ""},
		{"negative numbers", -10, 2, -5, false, ""},
		{"division by zero", 10, 0, 0, true, "division by zero"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)

			if tt.expectError {
				if err == nil {
					t.Errorf("Divide(%f, %f) expected error containing '%s', got nil", tt.a, tt.b, tt.errorContains)
				} else if tt.errorContains != "" && !errors.Is(err, errors.New(tt.errorContains)) {
					// For simple string matching
					if !contains(err.Error(), tt.errorContains) {
						t.Errorf("Divide(%f, %f) error = %v, want error containing '%s'", tt.a, tt.b, err, tt.errorContains)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Divide(%f, %f) unexpected error: %v", tt.a, tt.b, err)
				}
				if result != tt.expected {
					t.Errorf("Divide(%f, %f) = %f; want %f", tt.a, tt.b, result, tt.expected)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ========================================
// EXAMPLE 3: TEST SETUP AND TEARDOWN
// ========================================

var globalSetupDone bool

// TestMain runs before and after all tests in the package
func TestMain(m *testing.M) {
	fmt.Println("\n=== TestMain: Setting up tests ===")
	// Setup: Initialize resources, database connections, etc.
	globalSetupDone = true

	// Run all tests
	code := m.Run()

	fmt.Println("=== TestMain: Tearing down tests ===")
	// Teardown: Close connections, cleanup, etc.

	// Exit with the test result code
	// (non-zero means tests failed)
	fmt.Printf("\n=== TestMain: Exiting with code %d ===\n", code)
}

func TestGlobalSetup(t *testing.T) {
	if !globalSetupDone {
		t.Error("TestMain setup was not called")
	}
}

// Example with setup/teardown for specific test
type TestDatabase struct {
	data map[string]string
}

func setupTestDB(t *testing.T) *TestDatabase {
	t.Helper() // Marks this as a test helper function
	// Line numbers in errors will point to caller, not here

	db := &TestDatabase{
		data: make(map[string]string),
	}

	// Setup: Add test data
	db.data["user1"] = "Alice"
	db.data["user2"] = "Bob"

	t.Log("Test database setup complete")
	return db
}

func TestWithSetup(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		// Teardown: Cleanup
		db.data = nil
		t.Log("Test database cleanup complete")
	}()

	// Test code
	if len(db.data) != 2 {
		t.Errorf("expected 2 users, got %d", len(db.data))
	}
}

// ========================================
// EXAMPLE 4: TEST HELPERS
// ========================================

// Helper functions reduce duplication and improve readability

// assertEqual checks if two values are equal
func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// assertNotEqual checks if two values are not equal
func assertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("values are equal: %v, expected them to be different", got)
	}
}

// assertNil checks if a value is nil
func assertNil(t *testing.T, got interface{}) {
	t.Helper()
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

// assertNotNil checks if a value is not nil
func assertNotNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil {
		t.Error("expected non-nil value, got nil")
	}
}

// assertError checks if an error occurred (or not)
func assertError(t *testing.T, err error, want bool) {
	t.Helper()
	if (err != nil) != want {
		t.Errorf("error = %v, want error? %v", err, want)
	}
}

// assertPanic checks if a function panics
func assertPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, but function didn't panic")
		}
	}()
	f()
}

func TestAssertions(t *testing.T) {
	// Test assertEqual
	assertEqual(t, 1, 1)
	assertEqual(t, "hello", "hello")

	// Test assertNotEqual
	assertNotEqual(t, 1, 2)
	assertNotEqual(t, "hello", "world")

	// Test assertNil
	assertNil(t, nil)

	// Test assertError
	assertError(t, errors.New("test error"), true)
	assertError(t, nil, false)
}

func TestPanic(t *testing.T) {
	assertPanic(t, func() {
		panic("expected panic")
	})
}

// ========================================
// EXAMPLE 5: MOCKING WITH INTERFACES
// ========================================

// Database interface defines the contract
type Database interface {
	GetUser(id int) (*User, error)
	SaveUser(user *User) error
}

// User represents a user in the system
type User struct {
	ID    int
	Name  string
	Email string
}

// UserService uses a Database
type UserService struct {
	db Database
}

func NewUserService(db Database) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUser(id int) (*User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}
	return s.db.GetUser(id)
}

func (s *UserService) SaveUser(user *User) error {
	if user.Name == "" {
		return errors.New("user name is required")
	}
	return s.db.SaveUser(user)
}

// MockDatabase is a mock implementation for testing
type MockDatabase struct {
	users     map[int]*User
	getError  bool
	saveError bool
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		users: make(map[int]*User),
	}
}

func (m *MockDatabase) GetUser(id int) (*User, error) {
	if m.getError {
		return nil, errors.New("database error")
	}
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockDatabase) SaveUser(user *User) error {
	if m.saveError {
		return errors.New("database error")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockDatabase) AddUser(user *User) {
	m.users[user.ID] = user
}

func TestUserService_GetUser(t *testing.T) {
	tests := []struct {
		name        string
		userID      int
		setupMock   func(*MockDatabase)
		expectError bool
		errorMsg    string
		checkUser   func(*testing.T, *User)
	}{
		{
			name:   "successful user retrieval",
			userID: 1,
			setupMock: func(m *MockDatabase) {
				m.AddUser(&User{ID: 1, Name: "Alice", Email: "alice@example.com"})
			},
			expectError: false,
			checkUser: func(t *testing.T, user *User) {
				t.Helper()
				assertEqual(t, user.Name, "Alice")
				assertEqual(t, user.Email, "alice@example.com")
			},
		},
		{
			name:        "user not found",
			userID:      999,
			setupMock:   func(m *MockDatabase) {},
			expectError: true,
			errorMsg:    "user not found",
		},
		{
			name:   "database error",
			userID: 1,
			setupMock: func(m *MockDatabase) {
				m.getError = true
			},
			expectError: true,
			errorMsg:    "database error",
		},
		{
			name:        "invalid user ID",
			userID:      -1,
			setupMock:   func(m *MockDatabase) {},
			expectError: true,
			errorMsg:    "invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := NewMockDatabase()
			tt.setupMock(mockDB)

			service := NewUserService(mockDB)
			user, err := service.GetUser(tt.userID)

			if tt.expectError {
				assertError(t, err, true)
				if err != nil && tt.errorMsg != "" {
					if !contains(err.Error(), tt.errorMsg) {
						t.Errorf("error = %v, want error containing '%s'", err, tt.errorMsg)
					}
				}
			} else {
				assertError(t, err, false)
				if tt.checkUser != nil {
					tt.checkUser(t, user)
				}
			}
		})
	}
}

func TestUserService_SaveUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *User
		setupMock   func(*MockDatabase)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful user save",
			user: &User{ID: 1, Name: "Bob", Email: "bob@example.com"},
			setupMock: func(m *MockDatabase) {
				// No special setup needed
			},
			expectError: false,
		},
		{
			name: "empty user name",
			user: &User{ID: 1, Name: "", Email: "test@example.com"},
			setupMock: func(m *MockDatabase) {},
			expectError: true,
			errorMsg: "user name is required",
		},
		{
			name: "database error on save",
			user: &User{ID: 1, Name: "Charlie", Email: "charlie@example.com"},
			setupMock: func(m *MockDatabase) {
				m.saveError = true
			},
			expectError: true,
			errorMsg: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := NewMockDatabase()
			tt.setupMock(mockDB)

			service := NewUserService(mockDB)
			err := service.SaveUser(tt.user)

			if tt.expectError {
				assertError(t, err, true)
				if err != nil && tt.errorMsg != "" {
					if !contains(err.Error(), tt.errorMsg) {
						t.Errorf("error = %v, want error containing '%s'", err, tt.errorMsg)
					}
				}
			} else {
				assertError(t, err, false)
			}
		})
	}
}

// ========================================
// EXAMPLE 6: RACE DETECTION
// ========================================

// UnsafeCounter has a race condition
type UnsafeCounter struct {
	value int
}

func (c *UnsafeCounter) Increment() {
	c.value++ // Race condition: concurrent writes
}

func (c *UnsafeCounter) Value() int {
	return c.value
}

// SafeCounter uses mutex for thread safety
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func TestUnsafeCounter(t *testing.T) {
	// This test will FAIL with race detector
	// Run: go test -race
	/*
		counter := &UnsafeCounter{}

		// Start multiple goroutines
		for i := 0; i < 100; i++ {
			go func() {
				for j := 0; j < 100; j++ {
					counter.Increment()
				}
			}()
		}

		// Wait a bit for goroutines to finish
		time.Sleep(100 * time.Millisecond)

		expected := 100 * 100
		if counter.Value() != expected {
			t.Errorf("expected %d, got %d", expected, counter.Value())
		}
	*/
	t.Skip("Skipping unsafe counter test - run with -race to see the issue")
}

func TestSafeCounter(t *testing.T) {
	counter := &SafeCounter{}

	// Use WaitGroup to wait for all goroutines
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	expected := 100 * 100
	assertEqual(t, counter.Value(), expected)
}

// ========================================
// EXAMPLE 7: SUBTESTS
// ========================================

func Process(input string) string {
	return input
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"lowercase input", "hello", "hello"},
		{"uppercase input", "WORLD", "WORLD"},
		{"mixed case", "Go", "Go"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Process(tt.input)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

// ========================================
// EXAMPLE 8: TEST SKIP AND PARALLEL
// ========================================

func TestSlowOperation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping slow operation in short mode")
	}

	// Simulate slow operation
	time.Sleep(100 * time.Millisecond)
	t.Log("slow operation complete")
}

func TestParallelOperations(t *testing.T) {
	tests := []struct {
		name string
		data int
	}{
		{"test 1", 1},
		{"test 2", 2},
		{"test 3", 3},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run this test in parallel with others

			// Simulate some work
			time.Sleep(time.Duration(tt.data*10) * time.Millisecond)
			t.Logf("completed test %s with data %d", tt.name, tt.data)
		})
	}
}

// ========================================
// EXAMPLE 9: TESTING GOROUTINES
// ========================================

func ProcessWithChannel(input <-chan int, output chan<- int) {
	for num := range input {
		output <- num * 2
	}
	close(output)
}

func TestProcessWithChannel(t *testing.T) {
	input := make(chan int, 3)
	output := make(chan int, 3)

	// Send data to input channel
	input <- 1
	input <- 2
	input <- 3
	close(input)

	// Process in goroutine
	go ProcessWithChannel(input, output)

	// Collect results
	var results []int
	for result := range output {
		results = append(results, result)
	}

	// Verify results
	expected := []int{2, 4, 6}
	if len(results) != len(expected) {
		t.Fatalf("got %d results, expected %d", len(results), len(expected))
	}

	for i, result := range results {
		if result != expected[i] {
			t.Errorf("result[%d] = %d, expected %d", i, result, expected[i])
		}
	}
}

// ========================================
// EXAMPLE 10: TEMPORARY FILES IN TESTS
// ========================================

func TestFileOperations(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Cleanup after test

	// Write test data
	testData := "Hello, World!"
	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Read and verify
	data, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}

	if string(data) != testData {
		t.Errorf("got %s, want %s", string(data), testData)
	}
}

// ========================================
// EXAMPLE 11: EXAMPLE-BASED TESTS
// ========================================

func ExampleAdd() {
	result := Add(2, 3)
	fmt.Println(result)
	// Output: 5
}

func ExampleDivide() {
	result, err := Divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(result)
	// Output: 5
}

// ========================================
// EXAMPLE 12: CUSTOM TEST TYPES
// ========================================

// Calculator for testing
type Calculator struct {
	result int
}

func NewCalculator() *Calculator {
	return &Calculator{result: 0}
}

func (c *Calculator) Add(value int) *Calculator {
	c.result += value
	return c
}

func (c *Calculator) Subtract(value int) *Calculator {
	c.result -= value
	return c
}

func (c *Calculator) Multiply(value int) *Calculator {
	c.result *= value
	return c
}

func (c *Calculator) Result() int {
	return c.result
}

func TestCalculator(t *testing.T) {
	tests := []struct {
		name     string
		ops      func(*Calculator)
		expected int
	}{
		{
			name: "simple addition",
			ops: func(c *Calculator) {
				c.Add(5).Add(3)
			},
			expected: 8,
		},
		{
			name: "addition and subtraction",
			ops: func(c *Calculator) {
				c.Add(10).Subtract(3)
			},
			expected: 7,
		},
		{
			name: "complex operations",
			ops: func(c *Calculator) {
				c.Add(10).Multiply(2).Subtract(5)
			},
			expected: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewCalculator()
			tt.ops(calc)
			assertEqual(t, calc.Result(), tt.expected)
		})
	}
}

// ========================================
// EXAMPLE 13: TESTING ERROR WRAPPING
// ========================================

func processFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Simulate some processing
	if filename == "error.txt" {
		return fmt.Errorf("processing error in %s", filename)
	}

	return nil
}

func TestProcessFile(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		expectError bool
		checkError  func(*testing.T, error)
	}{
		{
			name:        "non-existent file",
			filename:    "nonexistent.txt",
			expectError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				// Check if error wraps os.ErrNotExist
				if !errors.Is(err, os.ErrNotExist) {
					// Might be wrapped, so check the message
					if !contains(err.Error(), "no such file") {
						t.Errorf("expected file not found error, got %v", err)
					}
				}
			},
		},
		{
			name:        "processing error",
			filename:    "error.txt",
			expectError: true,
			checkError: func(t *testing.T, err error) {
				t.Helper()
				if !contains(err.Error(), "processing error") {
					t.Errorf("expected processing error, got %v", err)
				}
			},
		},
		{
			name:        "create temp file for success",
			filename:    "",
			expectError: false,
			checkError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filename string

			// For success case, create a temp file
			if !tt.expectError && tt.filename == "" {
				tmpFile, err := os.CreateTemp("", "test_*.txt")
				if err != nil {
					t.Fatalf("failed to create temp file: %v", err)
				}
				defer os.Remove(tmpFile.Name())
				tmpFile.Close()
				filename = tmpFile.Name()
			} else {
				filename = tt.filename
			}

			err := processFile(filename)

			if tt.expectError {
				assertError(t, err, true)
				if tt.checkError != nil {
					tt.checkError(t, err)
				}
			} else {
				assertError(t, err, false)
			}
		})
	}
}

// ========================================
// EXAMPLE 14: TESTING WITH TIMEOUT
// ========================================

func TestWithTimeout(t *testing.T) {
	// This test has a timeout of 100ms
	// If it takes longer, it will fail
	deadline := time.Now().Add(100 * time.Millisecond)

	// Simulate work
	done := make(chan bool)
	go func() {
		time.Sleep(50 * time.Millisecond)
		done <- true
	}()

	select {
	case <-done:
		// Test completed in time
		t.Log("operation completed within timeout")
	case <-time.After(200 * time.Millisecond):
		t.Error("operation timed out")
	}

	// Check we're still before deadline
	if time.Now().After(deadline) {
		t.Error("test exceeded deadline")
	}
}

// ========================================
// EXAMPLE 15: TESTING RANDOM BEHAVIOR
// ========================================

func RandomInRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func TestRandomInRange(t *testing.T) {
	// Seed for reproducibility (in real tests, use fixed seed)
	rand.Seed(time.Now().UnixNano())

	min, max := 1, 10
	numIterations := 1000

	for i := 0; i < numIterations; i++ {
		result := RandomInRange(min, max)

		if result < min || result > max {
			t.Errorf("result %d is outside range [%d, %d]", result, min, max)
		}
	}

	t.Logf("Completed %d iterations", numIterations)
}

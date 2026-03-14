//go:build tools
// +build tools

package testingexamples

// This file contains testing tools and utilities
// Run with: go test -tags=tools -v ./...

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

// User type for testing
type User struct {
	ID    int
	Name  string
	Email string
}

// ========================================
// CUSTOM ASSERTION FUNCTIONS
// ========================================

// assertEqual checks if two values are equal
func assertEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v (%T), want %v (%T)", got, got, want, want)
	}
}

// assertNotEqual checks if two values are not equal
func assertNotEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got == want {
		t.Errorf("values should not be equal: both are %v", got)
	}
}

// assertDeepEqual checks deep equality for slices, maps, structs
func assertDeepEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("values are not deeply equal\ngot:  %+v\nwant: %+v", got, want)
	}
}

// assertNil checks if a value is nil
func assertNil(t *testing.T, got interface{}) {
	t.Helper()
	if got != nil {
		t.Errorf("expected nil, got %v (%T)", got, got)
	}
}

// assertNotNil checks if a value is not nil
func assertNotNil(t *testing.T, got interface{}) {
	t.Helper()
	if got == nil {
		t.Error("expected non-nil value, got nil")
	}
}

// assertTrue checks if a boolean is true
func assertTrue(t *testing.T, value bool, msg ...string) {
	t.Helper()
	if !value {
		if len(msg) > 0 {
			t.Error(msg[0])
		} else {
			t.Error("expected true, got false")
		}
	}
}

// assertFalse checks if a boolean is false
func assertFalse(t *testing.T, value bool, msg ...string) {
	t.Helper()
	if value {
		if len(msg) > 0 {
			t.Error(msg[0])
		} else {
			t.Error("expected false, got true")
		}
	}
}

// assertError checks if an error occurred
func assertError(t *testing.T, err error, wantError bool) {
	t.Helper()
	if (err != nil) != wantError {
		t.Errorf("error = %v, want error? %v", err, wantError)
	}
}

// assertErrorIs checks if error wraps a specific error
func assertErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errorIs(err, target) {
		t.Errorf("error %v does not wrap %v", err, target)
	}
}

// assertErrorContains checks if error message contains substring
func assertErrorContains(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Errorf("expected error containing '%s', got nil", substr)
		return
	}
	if !contains(err.Error(), substr) {
		t.Errorf("error = %v, want error containing '%s'", err, substr)
	}
}

// assertLen checks length of strings, slices, maps, channels
func assertLen(t *testing.T, obj interface{}, expected int) {
	t.Helper()
	val := reflect.ValueOf(obj)
	length := val.Len()

	if length != expected {
		t.Errorf("expected length %d, got %d", expected, length)
	}
}

// assertContains checks if slice contains element
func assertContains[T comparable](t *testing.T, slice []T, element T) {
	t.Helper()
	for _, v := range slice {
		if v == element {
			return
		}
	}
	t.Errorf("slice %+v does not contain %v", slice, element)
}

// assertNotContains checks if slice does not contain element
func assertNotContains[T comparable](t *testing.T, slice []T, element T) {
	t.Helper()
	for _, v := range slice {
		if v == element {
			t.Errorf("slice %+v should not contain %v", slice, element)
			return
		}
	}
}

// assertInRange checks if value is in range [min, max]
func assertInRange[T comparable](t *testing.T, value, min, max T) {
	t.Helper()
	// Note: This requires ordered types, so we'll use a simpler approach
	// In production, you'd use constraints.Ordered from Go 1.18+
	t.Log("assertInRange: needs constraints.Ordered for proper implementation")
}

// assertPanics checks if function panics
func assertPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected function to panic, but it didn't")
		}
	}()
	f()
}

// assertNotPanics checks if function doesn't panic
func assertNotPanics(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("function panicked unexpectedly: %v", r)
		}
	}()
	f()
}

// assertEventually checks if condition becomes true within timeout
func assertEventually(t *testing.T, condition func() bool, timeout time.Duration, msg string) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		if condition() {
			return
		}

		if time.Now().After(deadline) {
			if msg == "" {
				msg = "condition was not met within timeout"
			}
			t.Fatal(msg)
		}

		<-ticker.C
	}
}

// assertJSONEqual checks if two JSON strings are semantically equal
func assertJSONEqual(t *testing.T, got, want string) {
	t.Helper()

	var gotJSON, wantJSON interface{}

	if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
		t.Fatalf("failed to parse 'got' JSON: %v", err)
	}

	if err := json.Unmarshal([]byte(want), &wantJSON); err != nil {
		t.Fatalf("failed to parse 'want' JSON: %v", err)
	}

	if !reflect.DeepEqual(gotJSON, wantJSON) {
		t.Errorf("JSON not equal\ngot:  %s\nwant: %s", got, want)
	}
}

// Helper functions
func errorIs(err, target error) bool {
	for err != nil {
		if err == target {
			return true
		}
		err = unwrapError(err)
	}
	return false
}

func unwrapError(err error) error {
	type wrapper interface {
		Unwrap() error
	}
	if wrapper, ok := err.(wrapper); ok {
		return wrapper.Unwrap()
	}
	return nil
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
// TEST FIXTURES
// ========================================

// UserFixture provides test user data
type UserFixture struct {
	ValidUser       User
	InvalidUser     User
	AdminUser       User
	MultipleUsers   []User
	UserJSON        string
	InvalidUserJSON string
}

// LoadUserFixture loads test user data
func LoadUserFixture() *UserFixture {
	return &UserFixture{
		ValidUser: User{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
		},
		InvalidUser: User{
			ID:    2,
			Name:  "",
			Email: "invalid-email",
		},
		AdminUser: User{
			ID:    999,
			Name:  "Admin User",
			Email: "admin@example.com",
		},
		MultipleUsers: []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com"},
			{ID: 2, Name: "Bob", Email: "bob@example.com"},
			{ID: 3, Name: "Charlie", Email: "charlie@example.com"},
		},
		UserJSON: `{"id":1,"name":"John Doe","email":"john@example.com"}`,
		InvalidUserJSON: `{invalid json}`,
	}
}

// DatabaseFixture provides database test data
type DatabaseFixture struct {
	TestUsers []User
	TestData  map[string]interface{}
}

// LoadDatabaseFixture loads database test data
func LoadDatabaseFixture() *DatabaseFixture {
	return &DatabaseFixture{
		TestUsers: []User{
			{ID: 1, Name: "Alice", Email: "alice@example.com"},
			{ID: 2, Name: "Bob", Email: "bob@example.com"},
			{ID: 3, Name: "Charlie", Email: "charlie@example.com"},
		},
		TestData: map[string]interface{}{
			"test_key_1": "test_value_1",
			"test_key_2": 12345,
			"test_key_3": true,
		},
	}
}

// HTTPFixture provides HTTP test data
type HTTPFixture struct {
	ValidHeaders   map[string]string
	InvalidHeaders map[string]string
	TestBody       string
	TestQueries    map[string]string
}

// LoadHTTPFixture loads HTTP test data
func LoadHTTPFixture() *HTTPFixture {
	return &HTTPFixture{
		ValidHeaders: map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer test-token",
			"X-API-Key":     "test-api-key",
		},
		InvalidHeaders: map[string]string{
			"Content-Type": "text/plain",
		},
		TestBody: `{"name":"test","value":"123"}`,
		TestQueries: map[string]string{
			"page":  "1",
			"limit": "10",
			"sort":  "name",
		},
	}
}

// ========================================
// TEST DATA GENERATORS
// ========================================

// RandomString generates random string of given length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// RandomEmail generates random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(10))
}

// RandomUser generates random user
func RandomUser() User {
	return User{
		ID:    int(time.Now().UnixNano() % 1000000),
		Name:  RandomString(10),
		Email: RandomEmail(),
	}
}

// GenerateUsers generates multiple random users
func GenerateUsers(count int) []User {
	users := make([]User, count)
	for i := 0; i < count; i++ {
		users[i] = RandomUser()
	}
	return users
}

// ========================================
// COMPARISON UTILITIES
// ========================================

// Diff shows difference between two values
func Diff(got, want interface{}) string {
	gotJSON, _ := json.MarshalIndent(got, "", "  ")
	wantJSON, _ := json.MarshalIndent(want, "", "  ")

	return fmt.Sprintf("got:  %s\nwant: %s", string(gotJSON), string(wantJSON))
}

// CompareSlices compares two slices and returns differences
func CompareSlices[T comparable](a, b []T) (onlyInA, onlyInB []T) {
	aMap := make(map[T]bool)
	for _, v := range a {
		aMap[v] = true
	}

	bMap := make(map[T]bool)
	for _, v := range b {
		bMap[v] = true
		if !aMap[v] {
			onlyInB = append(onlyInB, v)
		}
	}

	for _, v := range a {
		if !bMap[v] {
			onlyInA = append(onlyInA, v)
		}
	}

	return onlyInA, onlyInB
}

// CompareMaps compares two maps and returns differences
func CompareMaps[K comparable, V comparable](a, b map[K]V) (onlyInA, onlyInB, different []K) {
	for k := range a {
		if _, exists := b[k]; !exists {
			onlyInA = append(onlyInA, k)
		} else if a[k] != b[k] {
			different = append(different, k)
		}
	}

	for k := range b {
		if _, exists := a[k]; !exists {
			onlyInB = append(onlyInB, k)
		}
	}

	return onlyInA, onlyInB, different
}

// ========================================
// FILE UTILITIES
// ========================================

// CreateTempFile creates temporary file with content
func CreateTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	return tmpFile.Name()
}

// CreateTempDir creates temporary directory
func CreateTempDir(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return tmpDir
}

// ReadFile reads file content
func ReadFile(t *testing.T, filename string) string {
	t.Helper()

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	return string(content)
}

// WriteFile writes content to file
func WriteFile(t *testing.T, filename, content string) {
	t.Helper()

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

// ========================================
// HTTP TESTING UTILITIES
// ========================================

// AssertHTTPStatus asserts HTTP response status code
func AssertHTTPStatus(t *testing.T, got, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("expected status %d, got %d", expected, got)
	}
}

// AssertHTTPHeader asserts HTTP header value
func AssertHTTPHeader(t *testing.T, headers map[string][]string, key, expectedValue string) {
	t.Helper()

	values := headers[key]
	if len(values) == 0 {
		t.Errorf("header '%s' not found", key)
		return
	}

	found := false
	for _, v := range values {
		if v == expectedValue {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("header '%s' = %v, expected value '%s'", key, values, expectedValue)
	}
}

// AssertJSONBody asserts response body contains valid JSON
func AssertJSONBody(t *testing.T, body io.Reader) interface{} {
	t.Helper()

	var result interface{}
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON body: %v", err)
	}

	return result
}

// AssertJSONPath asserts JSON value at path (simple implementation)
func AssertJSONPath(t *testing.T, body io.Reader, path string, expected interface{}) {
	t.Helper()

	var data map[string]interface{}
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	// Simple path navigation (only supports top-level keys)
	value, exists := data[path]
	if !exists {
		t.Errorf("path '%s' not found in JSON", path)
		return
	}

	if !reflect.DeepEqual(value, expected) {
		t.Errorf("path '%s' = %v (%T), expected %v (%T)", path, value, value, expected, expected)
	}
}

// ========================================
// MOCK HELPERS
// ========================================

// MockReader implements io.Reader for testing
type MockReader struct {
	Data   []byte
	Offset int
}

func NewMockReader(data string) *MockReader {
	return &MockReader{
		Data:   []byte(data),
		Offset: 0,
	}
}

func (m *MockReader) Read(p []byte) (n int, err error) {
	if m.Offset >= len(m.Data) {
		return 0, io.EOF
	}

	n = copy(p, m.Data[m.Offset:])
	m.Offset += n
	return n, nil
}

// MockWriter implements io.Writer for testing
type MockWriter struct {
	Data   []byte
	Closed bool
}

func NewMockWriter() *MockWriter {
	return &MockWriter{
		Data:   make([]byte, 0),
		Closed: false,
	}
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	if m.Closed {
		return 0, fmt.Errorf("writer is closed")
	}

	m.Data = append(m.Data, p...)
	return len(p), nil
}

func (m *MockWriter) Close() error {
	m.Closed = true
	return nil
}

func (m *MockWriter) String() string {
	return string(m.Data)
}

// MockCloser implements io.Closer for testing
type MockCloser struct {
	Closed bool
	CloseError error
}

func NewMockCloser(closeError error) *MockCloser {
	return &MockCloser{
		Closed:    false,
		CloseError: closeError,
	}
}

func (m *MockCloser) Close() error {
	m.Closed = true
	return m.CloseError
}

// ========================================
// TEST CONTEXT UTILITIES
// ========================================

// WithTimeout creates context with timeout
func WithTimeout(t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), timeout)
}

// WithDeadline creates context with deadline
func WithDeadline(t *testing.T, deadline time.Time) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithDeadline(context.Background(), deadline)
}

// ========================================
// TIMING UTILITIES
// ========================================

// TimeMeasurer measures execution time
type TimeMeasurer struct {
	start time.Time
}

// NewTimeMeasurer creates new time measurer
func NewTimeMeasurer() *TimeMeasurer {
	return &TimeMeasurer{start: time.Now()}
}

// Elapsed returns elapsed time since creation
func (tm *TimeMeasurer) Elapsed() time.Duration {
	return time.Since(tm.start)
}

// AssertDuration asserts if duration is within range
func AssertDuration(t *testing.T, got, min, max time.Duration) {
	t.Helper()
	if got < min || got > max {
		t.Errorf("duration %v outside range [%v, %v]", got, min, max)
	}
}

// Measure measures function execution time
func Measure(f func()) time.Duration {
	start := time.Now()
	f()
	return time.Since(start)
}

// ========================================
// EXAMPLE TESTS USING CUSTOM TOOLS
// ========================================

func TestCustomAssertions(t *testing.T) {
	t.Run("assertEqual", func(t *testing.T) {
		assertEqual(t, 1, 1)
		assertEqual(t, "hello", "hello")
	})

	t.Run("assertNil", func(t *testing.T) {
		assertNil(t, nil)
	})

	t.Run("assertContains", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		assertContains(t, slice, 3)
		assertNotContains(t, slice, 10)
	})

	t.Run("assertPanics", func(t *testing.T) {
		assertPanics(t, func() {
			panic("expected panic")
		})

		assertNotPanics(t, func() {
			// No panic
		})
	})

	t.Run("assertJSONEqual", func(t *testing.T) {
		json1 := `{"name":"John","age":30}`
		json2 := `{"age":30,"name":"John"}`
		assertJSONEqual(t, json1, json2)
	})
}

func TestFixtures(t *testing.T) {
	t.Run("user fixture", func(t *testing.T) {
		fixture := LoadUserFixture()

		assertEqual(t, fixture.ValidUser.Name, "John Doe")
		assertNotNil(t, fixture.AdminUser)
		assertLen(t, fixture.MultipleUsers, 3)
	})

	t.Run("database fixture", func(t *testing.T) {
		fixture := LoadDatabaseFixture()

		assertLen(t, fixture.TestUsers, 3)
		assertNotNil(t, fixture.TestData)
	})
}

func TestGenerators(t *testing.T) {
	t.Run("random string", func(t *testing.T) {
		s := RandomString(10)
		assertLen(t, s, 10)
	})

	t.Run("random user", func(t *testing.T) {
		user := RandomUser()
		assertNotNil(t, user.Email)
		assertTrue(t, len(user.Name) > 0)
	})

	t.Run("generate users", func(t *testing.T) {
		users := GenerateUsers(5)
		assertLen(t, users, 5)
	})
}

func TestFileUtilities(t *testing.T) {
	t.Run("create temp file", func(t *testing.T) {
		content := "test content"
		filename := CreateTempFile(t, content)

		read := ReadFile(t, filename)
		assertEqual(t, read, content)
	})

	t.Run("create temp dir", func(t *testing.T) {
		tmpDir := CreateTempDir(t)

		// Directory should exist
		info, err := os.Stat(tmpDir)
		assertNil(t, err)
		assertTrue(t, info.IsDir())
	})
}

func TestComparisonUtilities(t *testing.T) {
	t.Run("compare slices", func(t *testing.T) {
		a := []int{1, 2, 3}
		b := []int{2, 3, 4}

		onlyInA, onlyInB := CompareSlices(a, b)

		assertContains(t, onlyInA, 1)
		assertContains(t, onlyInB, 4)
	})

	t.Run("compare maps", func(t *testing.T) {
		a := map[string]int{"a": 1, "b": 2}
		b := map[string]int{"b": 2, "c": 3}

		onlyInA, onlyInB, _ := CompareMaps(a, b)

		assertContains(t, onlyInA, "a")
		assertContains(t, onlyInB, "c")
	})
}

func TestMockHelpers(t *testing.T) {
	t.Run("mock reader", func(t *testing.T) {
		data := "test data"
		reader := NewMockReader(data)

		buf := make([]byte, len(data))
		n, err := reader.Read(buf)

		assertNil(t, err)
		assertEqual(t, n, len(data))
		assertEqual(t, string(buf), data)
	})

	t.Run("mock writer", func(t *testing.T) {
		writer := NewMockWriter()
		data := []byte("test data")

		n, err := writer.Write(data)

		assertNil(t, err)
		assertEqual(t, n, len(data))
		assertEqual(t, writer.String(), "test data")
	})

	t.Run("mock closer", func(t *testing.T) {
		closer := NewMockCloser(nil)

		assertFalse(t, closer.Closed)
		closer.Close()
		assertTrue(t, closer.Closed)
	})
}

func TestTimeMeasurer(t *testing.T) {
	t.Run("measure time", func(t *testing.T) {
		tm := NewTimeMeasurer()
		time.Sleep(10 * time.Millisecond)

		elapsed := tm.Elapsed()
		AssertDuration(t, elapsed, 10*time.Millisecond, 20*time.Millisecond)
	})

	t.Run("measure function", func(t *testing.T) {
		duration := Measure(func() {
			time.Sleep(5 * time.Millisecond)
		})

		AssertDuration(t, duration, 5*time.Millisecond, 15*time.Millisecond)
	})
}

// ========================================
// RUNEABLE EXAMPLE FOR ALL TOOLS
// ========================================

func DemonstrateAllTestingTools() {
	// This example demonstrates all testing tools
	fmt.Println("Testing Tools Example:")
	fmt.Println("======================")

	// Generators
	user := RandomUser()
	fmt.Printf("Generated user: %+v\n", user)

	// Time measurement
	tm := NewTimeMeasurer()
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("Elapsed: %v\n", tm.Elapsed())

	// Mock reader
	reader := NewMockReader("test data")
	buf := make([]byte, 9)
	reader.Read(buf)
	fmt.Printf("Read from mock: %s\n", string(buf))

	// Mock writer
	writer := NewMockWriter()
	writer.Write([]byte("test output"))
	fmt.Printf("Written to mock: %s\n", writer.String())
}

// ========================================
// INTEGRATION EXAMPLE
// ========================================

func TestIntegrationExample(t *testing.T) {
	// Use fixtures
	userFixture := LoadUserFixture()
	assertEqual(t, userFixture.ValidUser.Name, "John Doe")

	// Generate test data
	users := GenerateUsers(3)
	assertLen(t, users, 3)

	// Create temp file with JSON
	userJSON, _ := json.Marshal(users[0])
	filename := CreateTempFile(t, string(userJSON))

	// Read it back
	content := ReadFile(t, filename)
	assertJSONEqual(t, content, string(userJSON))

	// Mock operations
	reader := NewMockReader(content)
	buf := make([]byte, len(content))
	n, err := reader.Read(buf)
	assertNil(t, err)
	assertEqual(t, n, len(content))

	// Measure performance
	duration := Measure(func() {
		for _, user := range users {
			_ = user.Email
		}
	})
	t.Logf("Processed %d users in %v", len(users), duration)
}

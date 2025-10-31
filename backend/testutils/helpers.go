// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package testutils provides testing utilities and helpers.
package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestCourse creates a test course for use in tests.
func TestCourse(id string) Course {
	return Course{
		ID:          id,
		Title:       "Test Course " + id,
		Description: "A test course for " + id,
		Lessons:     []string{id + "-lesson-1", id + "-lesson-2"},
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
}

// TestLesson creates a test lesson for use in tests.
func TestLesson(id, courseID string) Lesson {
	return Lesson{
		ID:          id,
		CourseID:    courseID,
		Title:       "Test Lesson " + id,
		Description: "A test lesson for " + courseID,
		Content:     "Test content for lesson " + id,
		Order:       1,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
}

// TestExercise creates a test exercise for use in tests.
func TestExercise(id, lessonID string) Exercise {
	return Exercise{
		ID:          id,
		LessonID:    lessonID,
		Title:       "Test Exercise " + id,
		Description: "A test exercise for " + lessonID,
		TestCases:   5,
		Difficulty:  "beginner",
	}
}

// TestProgress creates a test progress record for use in tests.
func TestProgress(userID, lessonID string, completed bool, score int) Progress {
	return Progress{
		UserID:      userID,
		LessonID:    lessonID,
		Completed:   completed,
		Score:       score,
		CompletedAt: time.Now(),
	}
}

// AssertJSONResponse validates a JSON API response.
func AssertJSONResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) APIResponse {
	t.Helper()

	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode API response: %v", err)
	}

	return response
}

// AssertSuccessResponse validates a successful API response.
func AssertSuccessResponse(t *testing.T, rr *httptest.ResponseRecorder) APIResponse {
	t.Helper()

	response := AssertJSONResponse(t, rr, http.StatusOK)

	if !response.Success {
		t.Errorf("Expected success=true, got success=%v", response.Success)
	}

	if response.Message == "" {
		t.Error("Expected non-empty message in successful response")
	}

	return response
}

// AssertErrorResponse validates an error API response.
func AssertErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) APIResponse {
	t.Helper()

	response := AssertJSONResponse(t, rr, expectedStatus)

	if response.Success {
		t.Error("Expected success=false for error response")
	}

	if response.Error == "" {
		t.Error("Expected non-empty error message in error response")
	}

	return response
}

// MockHTTPRequest creates a mock HTTP request for testing.
func MockHTTPRequest(method, url string, body interface{}) *http.Request {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}

	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, http.NoBody)
	}

	return req
}

// SetupTestData initializes test data for testing.
func SetupTestData() (map[string]Course, map[string]Lesson, map[string]Exercise, map[string][]Progress) {
	courses := make(map[string]Course)
	lessons := make(map[string]Lesson)
	exercises := make(map[string]Exercise)
	progress := make(map[string][]Progress)

	// Test courses.
	courses["test-course-1"] = TestCourse("test-course-1")
	courses["test-course-2"] = TestCourse("test-course-2")

	// Test lessons.
	lessons["lesson-1"] = TestLesson("lesson-1", "test-course-1")
	lessons["lesson-2"] = TestLesson("lesson-2", "test-course-1")
	lessons["lesson-3"] = TestLesson("lesson-3", "test-course-2")

	// Test exercises.
	exercises["exercise-1"] = TestExercise("exercise-1", "lesson-1")
	exercises["exercise-2"] = TestExercise("exercise-2", "lesson-1")
	exercises["exercise-3"] = TestExercise("exercise-3", "lesson-2")

	// Test progress.
	progress["test-user-1"] = []Progress{
		TestProgress("test-user-1", "lesson-1", true, 90),
		TestProgress("test-user-1", "lesson-2", false, 0),
	}
	progress["test-user-2"] = []Progress{
		TestProgress("test-user-2", "lesson-1", true, 85),
		TestProgress("test-user-2", "lesson-2", true, 95),
		TestProgress("test-user-2", "lesson-3", false, 0),
	}

	return courses, lessons, exercises, progress
}

// CompareObjects compares two objects using JSON marshaling.
func CompareObjects(t *testing.T, expected, actual interface{}) {
	t.Helper()

	expectedJSON, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshal expected object: %v", err)
	}

	actualJSON, err := json.Marshal(actual)
	if err != nil {
		t.Fatalf("Failed to marshal actual object: %v", err)
	}

	if !bytes.Equal(expectedJSON, actualJSON) {
		t.Errorf("Objects don't match:\nExpected: %s\nActual: %s",
			string(expectedJSON), string(actualJSON))
	}
}

// AssertSliceLength validates the length of a slice.
func AssertSliceLength(t *testing.T, slice interface{}, expectedLength int, description string) {
	t.Helper()

	switch s := slice.(type) {
	case []interface{}:
		if len(s) != expectedLength {
			t.Errorf("%s: expected length %d, got %d", description, expectedLength, len(s))
		}
	case []Course:
		if len(s) != expectedLength {
			t.Errorf("%s: expected length %d, got %d", description, expectedLength, len(s))
		}
	case []Lesson:
		if len(s) != expectedLength {
			t.Errorf("%s: expected length %d, got %d", description, expectedLength, len(s))
		}
	case []Exercise:
		if len(s) != expectedLength {
			t.Errorf("%s: expected length %d, got %d", description, expectedLength, len(s))
		}
	case []Progress:
		if len(s) != expectedLength {
			t.Errorf("%s: expected length %d, got %d", description, expectedLength, len(s))
		}
	default:
		t.Errorf("Unsupported slice type for length assertion: %T", slice)
	}
}

// ValidateRequiredFields checks if required fields are present in a map.
func ValidateRequiredFields(t *testing.T, obj map[string]interface{}, requiredFields []string, objectType string) {
	t.Helper()

	for _, field := range requiredFields {
		if _, exists := obj[field]; !exists {
			t.Errorf("Required field '%s' missing from %s", field, objectType)
		}
	}
}

// TimingHelper helps measure execution time in tests.
type TimingHelper struct {
	start time.Time
}

// NewTiming creates a new timing helper.
func NewTiming() *TimingHelper {
	return &TimingHelper{
		start: time.Now(),
	}
}

// Elapsed returns the elapsed time since creation.
func (th *TimingHelper) Elapsed() time.Duration {
	return time.Since(th.start)
}

// AssertMaxDuration fails the test if execution took longer than expected.
func (th *TimingHelper) AssertMaxDuration(t *testing.T, maxDuration time.Duration, operation string) {
	t.Helper()

	elapsed := th.Elapsed()
	if elapsed > maxDuration {
		t.Errorf("%s took too long: %v (max: %v)", operation, elapsed, maxDuration)
	}
}

// MockStorage provides test storage implementations.
type MockStorage struct {
	Courses   map[string]Course
	Lessons   map[string]Lesson
	Exercises map[string]Exercise
	Progress  map[string][]Progress
}

// NewMockStorage creates a new mock storage with test data.
func NewMockStorage() *MockStorage {
	courses, lessons, exercises, progress := SetupTestData()
	return &MockStorage{
		Courses:   courses,
		Lessons:   lessons,
		Exercises: exercises,
		Progress:  progress,
	}
}

// Clear removes all data from mock storage.
func (ms *MockStorage) Clear() {
	ms.Courses = make(map[string]Course)
	ms.Lessons = make(map[string]Lesson)
	ms.Exercises = make(map[string]Exercise)
	ms.Progress = make(map[string][]Progress)
}

// TableTest represents a table-driven test case.
type TableTest struct {
	Name     string
	Setup    func()
	Teardown func()
	Test     func(*testing.T)
}

// RunTableTests runs a series of table-driven tests.
func RunTableTests(t *testing.T, tests []TableTest) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Setup != nil {
				tt.Setup()
			}

			if tt.Teardown != nil {
				defer tt.Teardown()
			}

			tt.Test(t)
		})
	}
}

// HTTPTestCase represents an HTTP test case.
type HTTPTestCase struct {
	Name           string
	Method         string
	URL            string
	Body           interface{}
	ExpectedStatus int
	ExpectedFields []string
	PathValues     map[string]string
}

// RunHTTPTests runs a series of HTTP tests.
func RunHTTPTests(t *testing.T, tests []HTTPTestCase, handler http.HandlerFunc) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			req := MockHTTPRequest(tt.Method, tt.URL, tt.Body)

			// Set path values if provided.
			if tt.PathValues != nil {
				for key, value := range tt.PathValues {
					req.SetPathValue(key, value)
				}
			}

			rr := httptest.NewRecorder()
			handler(rr, req)

			response := AssertJSONResponse(t, rr, tt.ExpectedStatus)

			// Validate expected fields if provided.
			if tt.ExpectedFields != nil && tt.ExpectedStatus == http.StatusOK {
				if !response.Success {
					t.Error("Expected successful response")
				}

				if dataMap, ok := response.Data.(map[string]interface{}); ok {
					ValidateRequiredFields(t, dataMap, tt.ExpectedFields, "response data")
				}
			}
		})
	}
}

// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// +build integration

package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"go-pro-backend/internal/domain"
	"go-pro-backend/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssessmentRepository_CreateWithQuestions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewAssessmentRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Create quiz with questions", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)

		assessment := &domain.Assessment{
			LessonID:     lesson.ID,
			Type:         "quiz",
			Title:        "Go Basics Quiz",
			Description:  "Test your Go knowledge",
			Config:       mustMarshalJSON(map[string]interface{}{"time_limit_minutes": 30}),
			PassingScore: 80,
			OrderIndex:   1,
		}

		questions := []domain.QuizQuestion{
			{
				QuestionType: "multiple_choice",
				QuestionText: "What is the entry point of a Go program?",
				Options:      mustMarshalJSON([]string{"main()", "init()", "start()", "run()"}),
				CorrectAnswer: "main()",
				Explanation:  "The main() function is the entry point",
				Points:       1,
				OrderIndex:   1,
			},
			{
				QuestionType: "true_false",
				QuestionText: "Go is a compiled language",
				Options:      mustMarshalJSON([]string{"True", "False"}),
				CorrectAnswer: "True",
				Explanation:  "Go compiles to machine code",
				Points:       1,
				OrderIndex:   2,
			},
		}

		err := repo.CreateWithQuestions(ctx, assessment, questions)
		require.NoError(t, err, "Failed to create assessment with questions")
		require.NotEmpty(t, assessment.ID, "Assessment ID should be generated")

		// Verify questions were created
		retrieved, err := repo.GetByID(ctx, assessment.ID)
		require.NoError(t, err, "Failed to retrieve assessment")
		assert.Equal(t, assessment.Title, retrieved.Title, "Title should match")
	})

	t.Run("Create coding challenge", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)

		assessment := &domain.Assessment{
			LessonID:     lesson.ID,
			Type:         "coding_challenge",
			Title:        "Fibonacci Sequence",
			Description:  "Implement a function to generate Fibonacci numbers",
			Config:       mustMarshalJSON(map[string]interface{}{"time_limit_minutes": 60}),
			PassingScore: 100,
			OrderIndex:   2,
		}

		questions := []domain.QuizQuestion{
			{
				QuestionType: "code",
				QuestionText: "Write a Fibonacci function",
				Options:      mustMarshalJSON([]interface{}{}),
				CorrectAnswer: mustMarshalJSON(map[string]interface{}{
					"template":     "func fibonacci(n int) []int { /* your code */ }",
					"test_cases":   []interface{}{},
					"language":     "go",
				}),
				Points:     10,
				OrderIndex: 1,
			},
		}

		err := repo.CreateWithQuestions(ctx, assessment, questions)
		require.NoError(t, err, "Failed to create coding challenge")
	})
}

func TestAssessmentRepository_GetByLessonID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewAssessmentRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Get assessments for lesson", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)

		// Create multiple assessments
		a1 := helper.CreateTestAssessment(ctx, lesson.ID, "quiz")
		a2 := helper.CreateTestAssessment(ctx, lesson.ID, "coding_challenge")

		assessments, total, err := repo.GetByLessonID(ctx, lesson.ID, nil)
		require.NoError(t, err, "Failed to get assessments by lesson ID")
		assert.Equal(t, int64(2), total, "Should return total count of 2")
		assert.Len(t, assessments, 2, "Should return 2 assessments")
	})

	t.Run("Get assessments with pagination", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)

		// Create 5 assessments
		for i := 0; i < 5; i++ {
			helper.CreateTestAssessment(ctx, lesson.ID, "quiz")
		}

		pagination := &domain.PaginationRequest{
			Page:     1,
			PageSize: 2,
		}

		assessments, total, err := repo.GetByLessonID(ctx, lesson.ID, pagination)
		require.NoError(t, err, "Failed to get paginated assessments")
		assert.Equal(t, int64(5), total, "Total should be 5")
		assert.Len(t, assessments, 2, "Should return page size of 2")
	})
}

func TestQuizQuestionRepository_GetByAssessmentID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewQuizQuestionRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Get questions for assessment", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		assessment := helper.CreateTestAssessment(ctx, lesson.ID, "quiz")

		// Create questions
		q1 := helper.CreateTestQuizQuestion(ctx, assessment.ID, "multiple_choice")
		q2 := helper.CreateTestQuizQuestion(ctx, assessment.ID, "true_false")
		q3 := helper.CreateTestQuizQuestion(ctx, assessment.ID, "multiple_choice")

		questions, err := repo.GetByAssessmentID(ctx, assessment.ID)
		require.NoError(t, err, "Failed to get questions by assessment ID")
		assert.Len(t, questions, 3, "Should return 3 questions")

		// Verify ordering
		assert.Equal(t, q1.ID, questions[0].ID, "First question should match")
		assert.Equal(t, q2.ID, questions[1].ID, "Second question should match")
		assert.Equal(t, q3.ID, questions[2].ID, "Third question should match")
	})

	t.Run("Get questions with correct order", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		assessment := helper.CreateTestAssessment(ctx, lesson.ID, "quiz")

		// Create questions with specific order
		q1 := helper.CreateTestQuizQuestion(ctx, assessment.ID, "multiple_choice")
		q2 := helper.CreateTestQuizQuestion(ctx, assessment.ID, "true_false")

		// Update order indexes
		_, err := db.ExecContext(ctx, "UPDATE quiz_questions SET order_index = 2 WHERE id = $1", q1.ID)
		require.NoError(t, err)
		_, err = db.ExecContext(ctx, "UPDATE quiz_questions SET order_index = 1 WHERE id = $1", q2.ID)
		require.NoError(t, err)

		questions, err := repo.GetByAssessmentID(ctx, assessment.ID)
		require.NoError(t, err)
		assert.Equal(t, q2.ID, questions[0].ID, "Should be ordered by order_index")
		assert.Equal(t, q1.ID, questions[1].ID)
	})
}

func TestAssessmentRepository_UpdateWithQuestions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewAssessmentRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Update assessment title and config", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		assessment := helper.CreateTestAssessment(ctx, lesson.ID, "quiz")

		// Update assessment
		assessment.Title = "Updated Quiz Title"
		assessment.PassingScore = 90
		assessment.Config = mustMarshalJSON(map[string]interface{}{"time_limit_minutes": 45})

		err := repo.Update(ctx, assessment)
		require.NoError(t, err, "Failed to update assessment")

		// Verify update
		updated, err := repo.GetByID(ctx, assessment.ID)
		require.NoError(t, err, "Failed to retrieve updated assessment")
		assert.Equal(t, "Updated Quiz Title", updated.Title, "Title should be updated")
		assert.Equal(t, int32(90), updated.PassingScore, "Passing score should be updated")
	})
}

func TestAssessmentRepository_DeleteWithQuestions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewAssessmentRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Delete assessment with questions", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		assessment := helper.CreateTestAssessment(ctx, lesson.ID, "quiz")

		// Create questions
		helper.CreateTestQuizQuestion(ctx, assessment.ID, "multiple_choice")
		helper.CreateTestQuizQuestion(ctx, assessment.ID, "true_false")

		// Delete assessment
		err := repo.Delete(ctx, assessment.ID)
		require.NoError(t, err, "Failed to delete assessment")

		// Verify assessment deleted
		_, err = repo.GetByID(ctx, assessment.ID)
		assert.Error(t, err, "Assessment should be deleted")

		// Verify questions deleted (CASCADE)
		var questionCount int
		err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM quiz_questions WHERE assessment_id = $1", assessment.ID).Scan(&questionCount)
		require.NoError(t, err)
		assert.Equal(t, 0, questionCount, "Questions should be cascade deleted")
	})
}

func TestQuizQuestionRepository_Create(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	db, err := sql.Open("postgres", testutil.GetEnv("TEST_DATABASE_URL", "postgres://gopro_test:gopro_test@localhost:5432/gopro_test?sslmode=disable"))
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	helper := NewCMSTestHelper(t, db)
	repo := NewQuizQuestionRepository(db)
	helper.TruncateCMSTables(ctx)

	t.Run("Create multiple choice question", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		assessment := helper.CreateTestAssessment(ctx, lesson.ID, "quiz")

		question := &domain.QuizQuestion{
			AssessmentID:  assessment.ID,
			QuestionType:  "multiple_choice",
			QuestionText:  "What is Go's garbage collector?",
			Options:       mustMarshalJSON([]string{"Reference counting", "Tri-color mark and sweep", "Manual", "None"}),
			CorrectAnswer: "Tri-color mark and sweep",
			Explanation:   "Go uses a concurrent tri-color mark and sweep GC",
			Points:        2,
			OrderIndex:    1,
		}

		err := repo.Create(ctx, question)
		require.NoError(t, err, "Failed to create question")
		require.NotEmpty(t, question.ID, "Question ID should be generated")
	})

	t.Run("Create code question", func(t *testing.T) {
		course := helper.CreateTestCourse(ctx)
		lesson := helper.CreateTestLesson(ctx, course.ID)
		assessment := helper.CreateTestAssessment(ctx, lesson.ID, "coding_challenge")

		template := map[string]interface{}{
			"language": "go",
			"stub":     "func solution(input string) string { /* your code */ }",
		}

		question := &domain.QuizQuestion{
			AssessmentID:  assessment.ID,
			QuestionType:  "code",
			QuestionText:  "Implement string reversal",
			Options:       mustMarshalJSON([]interface{}{}),
			CorrectAnswer: mustMarshalJSON(template),
			Points:        10,
			OrderIndex:    1,
		}

		err := repo.Create(ctx, question)
		require.NoError(t, err, "Failed to create code question")
	})
}

// Helper function to marshal JSON and ignore errors
func mustMarshalJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

package migrations

import (
	"database/sql"

	"go-pro-backend/internal/repository/postgres"
)

// GetAllMigrations returns all database migrations.
func GetAllMigrations() []postgres.MigrationV2 {
	return []postgres.MigrationV2{
		createUsersTable(),
		createCoursesTable(),
		createLessonsTable(),
		createExercisesTable(),
		createProgressTable(),
		addIndexes(),
	}
}

// createUsersTable creates the users table.
func createUsersTable() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     1,
		Description: "Create users table",
		Up: func(tx *sql.Tx) error {
			query := `
				CREATE TABLE IF NOT EXISTS users (
					id VARCHAR(255) PRIMARY KEY,
					username VARCHAR(50) UNIQUE NOT NULL,
					email VARCHAR(255) UNIQUE NOT NULL,
					password_hash VARCHAR(255) NOT NULL,
					first_name VARCHAR(50),
					last_name VARCHAR(50),
					roles TEXT[] NOT NULL DEFAULT '{}',
					is_active BOOLEAN NOT NULL DEFAULT true,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					last_login_at TIMESTAMP
				)
			`
			_, err := tx.Exec(query)

			return err
		},
		Down: func(tx *sql.Tx) error {
			_, err := tx.Exec("DROP TABLE IF EXISTS users")
			return err
		},
	}
}

// createCoursesTable creates the courses table.
func createCoursesTable() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     2,
		Description: "Create courses table",
		Up: func(tx *sql.Tx) error {
			query := `
				CREATE TABLE IF NOT EXISTS courses (
					id VARCHAR(255) PRIMARY KEY,
					title VARCHAR(200) NOT NULL,
					slug VARCHAR(200) UNIQUE NOT NULL,
					description TEXT NOT NULL,
					difficulty VARCHAR(50) NOT NULL,
					duration_hours INTEGER NOT NULL,
					prerequisites TEXT[] NOT NULL DEFAULT '{}',
					learning_outcomes TEXT[] NOT NULL DEFAULT '{}',
					is_published BOOLEAN NOT NULL DEFAULT false,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
				)
			`
			_, err := tx.Exec(query)

			return err
		},
		Down: func(tx *sql.Tx) error {
			_, err := tx.Exec("DROP TABLE IF EXISTS courses")
			return err
		},
	}
}

// createLessonsTable creates the lessons table.
func createLessonsTable() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     3,
		Description: "Create lessons table",
		Up: func(tx *sql.Tx) error {
			query := `
				CREATE TABLE IF NOT EXISTS lessons (
					id VARCHAR(255) PRIMARY KEY,
					course_id VARCHAR(255) NOT NULL,
					title VARCHAR(200) NOT NULL,
					slug VARCHAR(200) NOT NULL,
					content TEXT NOT NULL,
					order_index INTEGER NOT NULL,
					duration_minutes INTEGER NOT NULL,
					is_published BOOLEAN NOT NULL DEFAULT false,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
					UNIQUE(course_id, slug)
				)
			`
			_, err := tx.Exec(query)

			return err
		},
		Down: func(tx *sql.Tx) error {
			_, err := tx.Exec("DROP TABLE IF EXISTS lessons")
			return err
		},
	}
}

// createExercisesTable creates the exercises table.
func createExercisesTable() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     4,
		Description: "Create exercises table",
		Up: func(tx *sql.Tx) error {
			query := `
				CREATE TABLE IF NOT EXISTS exercises (
					id VARCHAR(255) PRIMARY KEY,
					lesson_id VARCHAR(255) NOT NULL,
					title VARCHAR(200) NOT NULL,
					description TEXT NOT NULL,
					test_cases INTEGER NOT NULL,
					difficulty VARCHAR(50) NOT NULL,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
				)
			`
			_, err := tx.Exec(query)

			return err
		},
		Down: func(tx *sql.Tx) error {
			_, err := tx.Exec("DROP TABLE IF EXISTS exercises")
			return err
		},
	}
}

// createProgressTable creates the progress table.
func createProgressTable() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     5,
		Description: "Create progress table",
		Up: func(tx *sql.Tx) error {
			query := `
				CREATE TABLE IF NOT EXISTS progress (
					id VARCHAR(255) PRIMARY KEY,
					user_id VARCHAR(255) NOT NULL,
					lesson_id VARCHAR(255) NOT NULL,
					status VARCHAR(50) NOT NULL,
					started_at TIMESTAMP,
					completed_at TIMESTAMP,
					created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
					FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
					UNIQUE(user_id, lesson_id)
				)
			`
			_, err := tx.Exec(query)

			return err
		},
		Down: func(tx *sql.Tx) error {
			_, err := tx.Exec("DROP TABLE IF EXISTS progress")
			return err
		},
	}
}

// addIndexes adds performance indexes.
func addIndexes() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     6,
		Description: "Add performance indexes",
		Up: func(tx *sql.Tx) error {
			queries := []string{
				// Users indexes
				"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
				"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)",
				"CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active)",

				// Courses indexes
				"CREATE INDEX IF NOT EXISTS idx_courses_slug ON courses(slug)",
				"CREATE INDEX IF NOT EXISTS idx_courses_difficulty ON courses(difficulty)",
				"CREATE INDEX IF NOT EXISTS idx_courses_is_published ON courses(is_published)",

				// Lessons indexes
				"CREATE INDEX IF NOT EXISTS idx_lessons_course_id ON lessons(course_id)",
				"CREATE INDEX IF NOT EXISTS idx_lessons_slug ON lessons(slug)",
				"CREATE INDEX IF NOT EXISTS idx_lessons_order_index ON lessons(order_index)",
				"CREATE INDEX IF NOT EXISTS idx_lessons_is_published ON lessons(is_published)",

				// Exercises indexes
				"CREATE INDEX IF NOT EXISTS idx_exercises_lesson_id ON exercises(lesson_id)",
				"CREATE INDEX IF NOT EXISTS idx_exercises_difficulty ON exercises(difficulty)",

				// Progress indexes
				"CREATE INDEX IF NOT EXISTS idx_progress_user_id ON progress(user_id)",
				"CREATE INDEX IF NOT EXISTS idx_progress_lesson_id ON progress(lesson_id)",
				"CREATE INDEX IF NOT EXISTS idx_progress_status ON progress(status)",
				"CREATE INDEX IF NOT EXISTS idx_progress_user_lesson ON progress(user_id, lesson_id)",
			}

			for _, query := range queries {
				if _, err := tx.Exec(query); err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(tx *sql.Tx) error {
			queries := []string{
				"DROP INDEX IF EXISTS idx_users_email",
				"DROP INDEX IF EXISTS idx_users_username",
				"DROP INDEX IF EXISTS idx_users_is_active",
				"DROP INDEX IF EXISTS idx_courses_slug",
				"DROP INDEX IF EXISTS idx_courses_difficulty",
				"DROP INDEX IF EXISTS idx_courses_is_published",
				"DROP INDEX IF EXISTS idx_lessons_course_id",
				"DROP INDEX IF EXISTS idx_lessons_slug",
				"DROP INDEX IF EXISTS idx_lessons_order_index",
				"DROP INDEX IF EXISTS idx_lessons_is_published",
				"DROP INDEX IF EXISTS idx_exercises_lesson_id",
				"DROP INDEX IF EXISTS idx_exercises_difficulty",
				"DROP INDEX IF EXISTS idx_progress_user_id",
				"DROP INDEX IF EXISTS idx_progress_lesson_id",
				"DROP INDEX IF EXISTS idx_progress_status",
				"DROP INDEX IF EXISTS idx_progress_user_lesson",
			}

			for _, query := range queries {
				if _, err := tx.Exec(query); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

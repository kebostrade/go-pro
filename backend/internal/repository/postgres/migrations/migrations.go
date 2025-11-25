// GO-PRO Learning Platform Backend
// Copyright (c) 2025 GO-PRO Team
// Licensed under MIT License

// Package migrations provides database migration definitions for the GO-PRO Learning Platform.
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
		extendLessonsTable(),       // Version 7: Add detailed content fields
		seedLessonsData(),          // Version 8: Populate with 20 lessons
		addPerformanceIndexes(),    // Version 9: Add performance optimization indexes
		updateUsersTableForFirebase(), // Version 10: Add Firebase authentication fields
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

// addPerformanceIndexes adds performance optimization indexes.
// Version 9: Optimizes dashboard queries and curriculum ordering.
func addPerformanceIndexes() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     9,
		Description: "Add performance optimization indexes",
		Up: func(tx *sql.Tx) error {
			queries := []string{
				// Progress table composite index for dashboard queries
				// Optimizes: SELECT * FROM progress WHERE user_id = ? AND status = ? ORDER BY updated_at DESC
				"CREATE INDEX IF NOT EXISTS idx_progress_user_status_updated ON progress(user_id, status, updated_at DESC)",

				// Progress table covering index for efficient user progress lookups
				// Optimizes: SELECT user_id, lesson_id, status FROM progress WHERE user_id = ?
				"CREATE INDEX IF NOT EXISTS idx_progress_user_covering ON progress(user_id) INCLUDE (lesson_id, status, completed_at)",

				// Lessons table composite index for curriculum ordering
				// Optimizes: SELECT * FROM lessons WHERE course_id = ? ORDER BY order_index ASC
				"CREATE INDEX IF NOT EXISTS idx_lessons_course_order ON lessons(course_id, order_index ASC)",

				// Lessons table covering index for curriculum list views
				// Optimizes: SELECT id, title, slug, difficulty FROM lessons WHERE is_published = true
				"CREATE INDEX IF NOT EXISTS idx_lessons_published_covering ON lessons(is_published) INCLUDE (id, title, slug, difficulty, duration_minutes) WHERE is_published = true",

				// Progress table partial index for active lessons (performance optimization)
				// Optimizes: SELECT * FROM progress WHERE status = 'in_progress'
				"CREATE INDEX IF NOT EXISTS idx_progress_in_progress ON progress(user_id, updated_at) WHERE status = 'in_progress'",

				// Courses table covering index for published course listings
				// Optimizes: SELECT id, title, slug, difficulty FROM courses WHERE is_published = true
				"CREATE INDEX IF NOT EXISTS idx_courses_published_covering ON courses(is_published) INCLUDE (id, title, slug, difficulty) WHERE is_published = true",
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
				"DROP INDEX IF EXISTS idx_progress_user_status_updated",
				"DROP INDEX IF EXISTS idx_progress_user_covering",
				"DROP INDEX IF EXISTS idx_lessons_course_order",
				"DROP INDEX IF EXISTS idx_lessons_published_covering",
				"DROP INDEX IF EXISTS idx_progress_in_progress",
				"DROP INDEX IF EXISTS idx_courses_published_covering",
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

// updateUsersTableForFirebase updates users table to support Firebase authentication.
// Version 10: Adds Firebase-specific fields and updates schema to match domain model.
func updateUsersTableForFirebase() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     10,
		Description: "Update users table for Firebase authentication",
		Up: func(tx *sql.Tx) error {
			// 1. Add new Firebase-specific columns
			alterQueries := []string{
				// Add Firebase UID column (unique identifier from Firebase)
				"ALTER TABLE users ADD COLUMN IF NOT EXISTS firebase_uid VARCHAR(128)",

				// Add display_name from Firebase profile
				"ALTER TABLE users ADD COLUMN IF NOT EXISTS display_name VARCHAR(255)",

				// Add photo_url from Firebase profile
				"ALTER TABLE users ADD COLUMN IF NOT EXISTS photo_url TEXT",

				// Add single role column (replaces roles array)
				"ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'student'",
			}

			for _, query := range alterQueries {
				if _, err := tx.Exec(query); err != nil {
					return err
				}
			}

			// 2. Migrate existing roles array to single role (if data exists)
			// Use the first role in the array, default to 'student' if empty
			migrationQuery := `
				UPDATE users
				SET role = CASE
					WHEN array_length(roles, 1) > 0 THEN roles[1]
					ELSE 'student'
				END
				WHERE role IS NULL OR role = ''
			`
			if _, err := tx.Exec(migrationQuery); err != nil {
				return err
			}

			// 3. Generate temporary firebase_uid for existing users (UUID format)
			// In production, these should be replaced with real Firebase UIDs during user migration
			tempUidQuery := `
				UPDATE users
				SET firebase_uid = 'temp_' || id
				WHERE firebase_uid IS NULL OR firebase_uid = ''
			`
			if _, err := tx.Exec(tempUidQuery); err != nil {
				return err
			}

			// 4. Add constraints after data migration
			constraintQueries := []string{
				// Make firebase_uid unique and not null
				"ALTER TABLE users ALTER COLUMN firebase_uid SET NOT NULL",
				"ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS uq_users_firebase_uid UNIQUE (firebase_uid)",

				// Add role validation constraint
				"ALTER TABLE users ADD CONSTRAINT IF NOT EXISTS chk_users_role CHECK (role IN ('student', 'admin'))",

				// Make role not null
				"ALTER TABLE users ALTER COLUMN role SET NOT NULL",
			}

			for _, query := range constraintQueries {
				if _, err := tx.Exec(query); err != nil {
					return err
				}
			}

			// 5. Create indexes for Firebase fields
			indexQueries := []string{
				"CREATE INDEX IF NOT EXISTS idx_users_firebase_uid ON users(firebase_uid)",
				"CREATE INDEX IF NOT EXISTS idx_users_role ON users(role) WHERE is_active = TRUE",
			}

			for _, query := range indexQueries {
				if _, err := tx.Exec(query); err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(tx *sql.Tx) error {
			// Drop indexes
			dropIndexQueries := []string{
				"DROP INDEX IF EXISTS idx_users_firebase_uid",
				"DROP INDEX IF EXISTS idx_users_role",
			}

			for _, query := range dropIndexQueries {
				if _, err := tx.Exec(query); err != nil {
					return err
				}
			}

			// Drop constraints
			dropConstraintQueries := []string{
				"ALTER TABLE users DROP CONSTRAINT IF EXISTS uq_users_firebase_uid",
				"ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_role",
			}

			for _, query := range dropConstraintQueries {
				if _, err := tx.Exec(query); err != nil {
					return err
				}
			}

			// Drop columns
			dropColumnQuery := `
				ALTER TABLE users
				DROP COLUMN IF EXISTS firebase_uid,
				DROP COLUMN IF EXISTS display_name,
				DROP COLUMN IF EXISTS photo_url,
				DROP COLUMN IF EXISTS role
			`
			_, err := tx.Exec(dropColumnQuery)
			return err
		},
	}
}

// extendLessonsTable adds detailed content fields to lessons table.
func extendLessonsTable() postgres.MigrationV2 {
	return postgres.MigrationV2{
		Version:     7,
		Description: "Extend lessons table with detailed content fields",
		Up: func(tx *sql.Tx) error {
			query := `
				ALTER TABLE lessons
				ADD COLUMN IF NOT EXISTS description TEXT DEFAULT '',
				ADD COLUMN IF NOT EXISTS difficulty VARCHAR(50) DEFAULT 'beginner',
				ADD COLUMN IF NOT EXISTS phase VARCHAR(50) DEFAULT 'Foundations',
				ADD COLUMN IF NOT EXISTS objectives JSONB DEFAULT '[]'::jsonb,
				ADD COLUMN IF NOT EXISTS theory TEXT DEFAULT '',
				ADD COLUMN IF NOT EXISTS code_example TEXT DEFAULT '',
				ADD COLUMN IF NOT EXISTS solution TEXT DEFAULT '',
				ADD COLUMN IF NOT EXISTS exercises JSONB DEFAULT '[]'::jsonb,
				ADD COLUMN IF NOT EXISTS next_lesson_id VARCHAR(255),
				ADD COLUMN IF NOT EXISTS prev_lesson_id VARCHAR(255)
			`
			if _, err := tx.Exec(query); err != nil {
				return err
			}

			fkQueries := []string{
				`ALTER TABLE lessons ADD CONSTRAINT IF NOT EXISTS fk_lessons_next_lesson FOREIGN KEY (next_lesson_id) REFERENCES lessons(id) ON DELETE SET NULL`,
				`ALTER TABLE lessons ADD CONSTRAINT IF NOT EXISTS fk_lessons_prev_lesson FOREIGN KEY (prev_lesson_id) REFERENCES lessons(id) ON DELETE SET NULL`,
			}

			for _, fkQuery := range fkQueries {
				if _, err := tx.Exec(fkQuery); err != nil {
					return err
				}
			}

			indexQueries := []string{
				"CREATE INDEX IF NOT EXISTS idx_lessons_difficulty ON lessons(difficulty)",
				"CREATE INDEX IF NOT EXISTS idx_lessons_phase ON lessons(phase)",
				"CREATE INDEX IF NOT EXISTS idx_lessons_next_lesson ON lessons(next_lesson_id)",
				"CREATE INDEX IF NOT EXISTS idx_lessons_prev_lesson ON lessons(prev_lesson_id)",
				"CREATE INDEX IF NOT EXISTS idx_lessons_objectives ON lessons USING GIN (objectives)",
				"CREATE INDEX IF NOT EXISTS idx_lessons_exercises ON lessons USING GIN (exercises)",
			}

			for _, indexQuery := range indexQueries {
				if _, err := tx.Exec(indexQuery); err != nil {
					return err
				}
			}

			return nil
		},
		Down: func(tx *sql.Tx) error {
			indexQueries := []string{
				"DROP INDEX IF EXISTS idx_lessons_difficulty",
				"DROP INDEX IF EXISTS idx_lessons_phase",
				"DROP INDEX IF EXISTS idx_lessons_next_lesson",
				"DROP INDEX IF EXISTS idx_lessons_prev_lesson",
				"DROP INDEX IF EXISTS idx_lessons_objectives",
				"DROP INDEX IF EXISTS idx_lessons_exercises",
			}

			for _, indexQuery := range indexQueries {
				if _, err := tx.Exec(indexQuery); err != nil {
					return err
				}
			}

			fkQueries := []string{
				"ALTER TABLE lessons DROP CONSTRAINT IF EXISTS fk_lessons_next_lesson",
				"ALTER TABLE lessons DROP CONSTRAINT IF EXISTS fk_lessons_prev_lesson",
			}

			for _, fkQuery := range fkQueries {
				if _, err := tx.Exec(fkQuery); err != nil {
					return err
				}
			}

			query := `
				ALTER TABLE lessons
				DROP COLUMN IF EXISTS description,
				DROP COLUMN IF EXISTS difficulty,
				DROP COLUMN IF EXISTS phase,
				DROP COLUMN IF EXISTS objectives,
				DROP COLUMN IF EXISTS theory,
				DROP COLUMN IF EXISTS code_example,
				DROP COLUMN IF EXISTS solution,
				DROP COLUMN IF EXISTS exercises,
				DROP COLUMN IF EXISTS next_lesson_id,
				DROP COLUMN IF EXISTS prev_lesson_id
			`
			_, err := tx.Exec(query)
			return err
		},
	}
}

-- GO-PRO Learning Platform Backend
-- Copyright (c) 2025 GO-PRO Team
-- Licensed under MIT License

-- Rollback: Drop assessments table

DROP TRIGGER IF EXISTS trigger_update_assessments_updated_at ON assessments;
DROP FUNCTION IF EXISTS update_assessments_updated_at();
DROP INDEX IF EXISTS idx_assessments_created_at;
DROP INDEX IF EXISTS idx_assessments_type;
DROP INDEX IF EXISTS idx_assessments_lesson_id;
DROP TABLE IF EXISTS assessments;

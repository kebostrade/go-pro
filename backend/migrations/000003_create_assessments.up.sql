-- GO-PRO Learning Platform Backend
-- Copyright (c) 2025 GO-PRO Team
-- Licensed under MIT License

-- Migration: Create assessments table
-- Description: Stores assessments (quizzes, coding exercises, projects) linked to lessons

CREATE TABLE IF NOT EXISTS assessments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('quiz', 'coding_exercise', 'project')),
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    passing_score INT NOT NULL DEFAULT 70 CHECK (passing_score BETWEEN 0 AND 100),
    time_limit_minutes INT CHECK (time_limit_minutes BETWEEN 1 AND 180),
    order_index INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT assessments_lesson_id_type_unique UNIQUE (lesson_id, type, order_index)
);

-- Indexes for performance
CREATE INDEX idx_assessments_lesson_id ON assessments(lesson_id);
CREATE INDEX idx_assessments_type ON assessments(type);
CREATE INDEX idx_assessments_created_at ON assessments(created_at DESC);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_assessments_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_assessments_updated_at
    BEFORE UPDATE ON assessments
    FOR EACH ROW
    EXECUTE FUNCTION update_assessments_updated_at();

-- Comments for documentation
COMMENT ON TABLE assessments IS 'Stores assessment configurations for quizzes, coding exercises, and projects';
COMMENT ON COLUMN assessments.config IS 'Flexible JSONB config for assessment-specific settings';
COMMENT ON COLUMN assessments.passing_score IS 'Minimum score (0-100) required to pass';
COMMENT ON COLUMN assessments.time_limit_minutes IS 'Optional time limit for timed assessments';

-- GO-PRO Learning Platform Backend
-- Copyright (c) 2025 GO-PRO Team
-- Licensed under MIT License

-- Migration: Create submissions and submission_comments tables
-- Description: Stores student assessment submissions and inline comments

CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content JSONB NOT NULL DEFAULT '{}',
    score INT CHECK (score BETWEEN 0 AND 100),
    feedback TEXT,
    graded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    graded_at TIMESTAMP,
    submitted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'submitted' CHECK (status IN ('submitted', 'graded', 'returned')),
    CONSTRAINT submissions_user_assessment_unique UNIQUE (user_id, assessment_id)
);

CREATE TABLE IF NOT EXISTS submission_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment_text TEXT NOT NULL,
    line_number INT CHECK (line_number > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_submissions_assessment_id ON submissions(assessment_id);
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_status ON submissions(status);
CREATE INDEX idx_submissions_user_assessment ON submissions(user_id, assessment_id);
CREATE INDEX idx_submissions_submitted_at ON submissions(submitted_at DESC);
CREATE INDEX idx_submission_comments_submission_id ON submission_comments(submission_id);
CREATE INDEX idx_submission_comments_author_id ON submission_comments(author_id);

-- Comments for documentation
COMMENT ON TABLE submissions IS 'Stores student assessment submissions with flexible content';
COMMENT ON COLUMN submissions.content IS 'JSONB content for quiz answers, code, or project links';
COMMENT ON COLUMN submissions.status IS 'Workflow status: submitted -> graded -> returned';
COMMENT ON TABLE submission_comments IS 'Stores inline comments on submissions (e.g., code review feedback)';
COMMENT ON COLUMN submission_comments.line_number IS 'Optional line number for code-based comments';

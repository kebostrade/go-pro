-- GO-PRO Learning Platform Backend
-- Copyright (c) 2025 GO-PRO Team
-- Licensed under MIT License

-- Migration: Create peer_reviews table
-- Description: Stores peer review assignments and completions for project submissions

CREATE TABLE IF NOT EXISTS peer_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rubric_scores JSONB NOT NULL DEFAULT '{}',
    feedback TEXT NOT NULL DEFAULT '',
    inline_comments JSONB NOT NULL DEFAULT '[]',
    is_anonymous BOOLEAN NOT NULL DEFAULT TRUE,
    karma_points INT DEFAULT 0 CHECK (karma_points BETWEEN 0 AND 5),
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'late')),
    submitted_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deadline TIMESTAMP,
    CONSTRAINT peer_reviews_submission_reviewer_unique UNIQUE (submission_id, reviewer_id)
);

-- Indexes for performance
CREATE INDEX idx_peer_reviews_submission_id ON peer_reviews(submission_id);
CREATE INDEX idx_peer_reviews_reviewer_id ON peer_reviews(reviewer_id);
CREATE INDEX idx_peer_reviews_status ON peer_reviews(status);
CREATE INDEX idx_peer_reviews_created_at ON peer_reviews(created_at DESC);
CREATE INDEX idx_peer_reviews_deadline ON peer_reviews(deadline) WHERE deadline IS NOT NULL;
CREATE INDEX idx_peer_reviews_submission_status ON peer_reviews(submission_id, status);

-- Comments for documentation
COMMENT ON TABLE peer_reviews IS 'Stores peer review assignments and completions';
COMMENT ON COLUMN peer_reviews.rubric_scores IS 'JSONB mapping rubric criterion IDs to scores';
COMMENT ON COLUMN peer_reviews.inline_comments IS 'JSONB array of inline comments with line numbers';
COMMENT ON COLUMN peer_reviews.is_anonymous IS 'Whether reviewer identity is hidden from student';
COMMENT ON COLUMN peer_reviews.karma_points IS 'Points awarded by reviewee for helpfulness (1-5 stars)';
COMMENT ON COLUMN peer_reviews.status IS 'Review status: pending -> completed/late';
COMMENT ON COLUMN peer_reviews.deadline IS 'Optional deadline for review completion';

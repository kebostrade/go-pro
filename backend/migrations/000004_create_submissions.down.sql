-- GO-PRO Learning Platform Backend
-- Copyright (c) 2025 GO-PRO Team
-- Licensed under MIT License

-- Rollback: Drop submissions and submission_comments tables

DROP INDEX IF EXISTS idx_submission_comments_author_id;
DROP INDEX IF EXISTS idx_submission_comments_submission_id;
DROP INDEX IF EXISTS idx_submissions_submitted_at;
DROP INDEX IF EXISTS idx_submissions_user_assessment;
DROP INDEX IF EXISTS idx_submissions_status;
DROP INDEX IF EXISTS idx_submissions_user_id;
DROP INDEX IF EXISTS idx_submissions_assessment_id;
DROP TABLE IF EXISTS submission_comments;
DROP TABLE IF EXISTS submissions;

-- GO-PRO Learning Platform Backend
-- Copyright (c) 2025 GO-PRO Team
-- Licensed under MIT License

-- Migration: Rollback peer_reviews table
-- Description: Removes peer review table and indexes

DROP INDEX IF EXISTS idx_peer_reviews_submission_status;
DROP INDEX IF EXISTS idx_peer_reviews_deadline;
DROP INDEX IF EXISTS idx_peer_reviews_created_at;
DROP INDEX IF EXISTS idx_peer_reviews_status;
DROP INDEX IF EXISTS idx_peer_reviews_reviewer_id;
DROP INDEX IF EXISTS idx_peer_reviews_submission_id;

DROP TABLE IF EXISTS peer_reviews;

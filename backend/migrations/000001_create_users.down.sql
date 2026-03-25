-- Drop users table
-- Migration: 000001_create_users
-- Direction: DOWN

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_users_created_at ON users;
DROP INDEX IF EXISTS idx_users_is_active ON users;
DROP INDEX IF EXISTS idx_users_email ON users;
DROP TABLE IF EXISTS users;

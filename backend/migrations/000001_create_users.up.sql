-- Create users table for authentication and user management
-- Migration: 000001_create_users
-- Direction: UP

IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'users')
BEGIN
    CREATE TABLE users (
        id NVARCHAR(64) PRIMARY KEY,
        email NVARCHAR(254) UNIQUE NOT NULL,
        password_hash NVARCHAR(256) NOT NULL,
        roles NVARCHAR(MAX) NOT NULL DEFAULT 'user',
        created_at DATETIMEOFFSET NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        updated_at DATETIMEOFFSET NOT NULL DEFAULT SYSDATETIMEOFFSET(),
        last_login_at DATETIMEOFFSET NULL,
        is_active BIT NOT NULL DEFAULT 1,
        is_verified BIT NOT NULL DEFAULT 0
    );
END

-- Create indexes for common queries
IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_users_email' AND object_id = OBJECT_ID('users'))
    CREATE INDEX idx_users_email ON users(email);

IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_users_is_active' AND object_id = OBJECT_ID('users'))
    CREATE INDEX idx_users_is_active ON users(is_active);

IF NOT EXISTS (SELECT * FROM sys.indexes WHERE name = 'idx_users_created_at' AND object_id = OBJECT_ID('users'))
    CREATE INDEX idx_users_created_at ON users(created_at);

-- Create updated_at trigger
IF EXISTS (SELECT * FROM sys.triggers WHERE name = 'update_users_updated_at')
    DROP TRIGGER update_users_updated_at;
GO

CREATE TRIGGER update_users_updated_at
ON users
FOR UPDATE
AS
BEGIN
    SET NOCOUNT ON;
    UPDATE users
    SET updated_at = SYSDATETIMEOFFSET()
    FROM inserted
    WHERE users.id = inserted.id;
END;
GO

# Migration V10: Firebase Authentication Support

## Overview
This migration updates the existing `users` table to support Firebase Authentication integration by adding Firebase-specific fields and updating the schema to match the domain model.

## Version
- **Version**: 10
- **Description**: Update users table for Firebase authentication
- **File**: `internal/repository/postgres/migrations/migrations.go`
- **Function**: `updateUsersTableForFirebase()`

## Changes Applied

### 1. New Columns Added
| Column | Type | Description |
|--------|------|-------------|
| `firebase_uid` | VARCHAR(128) NOT NULL UNIQUE | Firebase User ID (unique identifier from Firebase Auth) |
| `display_name` | VARCHAR(255) | Full name from Firebase user profile |
| `photo_url` | TEXT | Profile picture URL from Firebase |
| `role` | VARCHAR(20) NOT NULL | Single role field ('student' or 'admin') |

### 2. Data Migration
- **Roles Array → Single Role**: Migrates existing `roles[]` array to single `role` field
  - Uses first role from array if exists
  - Defaults to 'student' if roles array is empty
- **Temporary Firebase UIDs**: Generates temporary firebase_uid for existing users
  - Format: `temp_{user_id}`
  - **Important**: These should be replaced with real Firebase UIDs during production migration

### 3. Constraints Added
| Constraint | Type | Description |
|------------|------|-------------|
| `uq_users_firebase_uid` | UNIQUE | Ensures firebase_uid uniqueness |
| `chk_users_role` | CHECK | Validates role is either 'student' or 'admin' |
| NOT NULL on `firebase_uid` | NOT NULL | firebase_uid is required |
| NOT NULL on `role` | NOT NULL | role is required |

### 4. Indexes Created
| Index | Target | Purpose |
|-------|--------|---------|
| `idx_users_firebase_uid` | firebase_uid | Fast lookups by Firebase UID |
| `idx_users_role` | role (WHERE is_active = TRUE) | Filtered index for active users by role |

## Domain Model Alignment

The migration aligns the database schema with the domain model:

```go
type User struct {
    ID           string     // Existing
    FirebaseUID  string     // NEW: Firebase User ID
    Username     string     // Existing
    Email        string     // Existing
    DisplayName  string     // NEW: From Firebase profile
    PhotoURL     string     // NEW: From Firebase profile
    PasswordHash string     // Existing (for potential future use)
    FirstName    string     // Existing
    LastName     string     // Existing
    Role         UserRole   // NEW: Single role (student/admin)
    Roles        []string   // Legacy: kept for backward compatibility
    IsActive     bool       // Existing
    CreatedAt    time.Time  // Existing
    UpdatedAt    time.Time  // Existing
    LastLoginAt  *time.Time // Existing
}
```

## Migration Execution

### Up Migration
1. Adds new columns with nullable constraints initially
2. Migrates existing data (roles array → single role)
3. Generates temporary Firebase UIDs for existing users
4. Applies NOT NULL and CHECK constraints after data migration
5. Creates performance indexes

### Down Migration (Rollback)
1. Drops performance indexes
2. Drops constraints
3. Removes new columns (firebase_uid, display_name, photo_url, role)
4. **Note**: Original roles array is preserved during rollback

## Production Deployment Notes

### Pre-Migration Checklist
- [ ] Backup database
- [ ] Review existing users table data
- [ ] Plan Firebase UID migration strategy for existing users

### Post-Migration Actions
1. **Replace Temporary Firebase UIDs**: Update all `temp_*` firebase_uid values with real Firebase UIDs
   ```sql
   UPDATE users
   SET firebase_uid = 'actual_firebase_uid_from_firebase_auth'
   WHERE firebase_uid LIKE 'temp_%';
   ```

2. **Verify Data Integrity**:
   ```sql
   -- Check for temporary UIDs still in database
   SELECT COUNT(*) FROM users WHERE firebase_uid LIKE 'temp_%';

   -- Verify role constraints
   SELECT COUNT(*), role FROM users GROUP BY role;

   -- Check firebase_uid uniqueness
   SELECT firebase_uid, COUNT(*)
   FROM users
   GROUP BY firebase_uid
   HAVING COUNT(*) > 1;
   ```

3. **Update Application Code**: Ensure UserRepository queries use new fields correctly

## Testing

### Unit Tests
```bash
cd backend
go test ./internal/repository/postgres/... -v
```

### Integration Tests
```bash
# Start test database
docker-compose -f docker-compose.test.yml up -d

# Run migration tests
go test ./internal/repository/postgres/migrations/... -tags=integration -v
```

### Manual Verification
```sql
-- Check schema
\d users

-- Verify indexes
\di users*

-- Check constraints
SELECT conname, pg_get_constraintdef(oid)
FROM pg_constraint
WHERE conrelid = 'users'::regclass;
```

## Rollback Procedure

If rollback is needed:

```go
// In Go code using MigrationManager
migrationManager := postgres.NewMigrationManager(db, logger)
migrationManager.RegisterMultiple(migrations.GetAllMigrations())
err := migrationManager.Down(ctx, 1) // Rollback 1 migration
```

Or manually:
```sql
-- Drop indexes
DROP INDEX IF EXISTS idx_users_firebase_uid;
DROP INDEX IF EXISTS idx_users_role;

-- Drop constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS uq_users_firebase_uid;
ALTER TABLE users DROP CONSTRAINT IF EXISTS chk_users_role;

-- Drop columns
ALTER TABLE users
    DROP COLUMN IF EXISTS firebase_uid,
    DROP COLUMN IF EXISTS display_name,
    DROP COLUMN IF EXISTS photo_url,
    DROP COLUMN IF EXISTS role;
```

## Related Files
- Domain model: `internal/domain/models.go` (User struct)
- Repository: `internal/repository/postgres/user.go`
- Migration manager: `internal/repository/postgres/migration.go`
- All migrations: `internal/repository/postgres/migrations/migrations.go`

## Migration Status

Check migration status:
```go
statuses, err := migrationManager.Status(ctx)
for _, status := range statuses {
    fmt.Printf("Version %d: %s - Applied: %v\n",
        status.Version,
        status.Description,
        status.Applied)
}
```

## References
- Firebase Authentication Docs: https://firebase.google.com/docs/auth
- Go domain model: `internal/domain/models.go`
- Initial schema: `scripts/init-db.sql`

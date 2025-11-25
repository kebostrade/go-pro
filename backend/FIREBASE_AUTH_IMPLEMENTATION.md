# Firebase Authentication Implementation Status

**Created**: 2025-01-24
**Issue**: #65
**Status**: 🟡 In Progress (60% Complete)

## Overview
Integration of Firebase Authentication with Email/Password, Google OAuth, and GitHub OAuth. Backend uses Firebase Admin SDK for token verification with hybrid storage (Firebase Auth + PostgreSQL).

## ✅ Completed Phases

### Phase 1: Dependencies & Configuration ✅
- **Firebase Admin SDK**: Added `firebase.google.com/go/v4` dependency
- **Environment Configuration**: Updated `.env.example` with Firebase settings:
  ```env
  FIREBASE_PROJECT_ID=your-firebase-project-id
  FIREBASE_CREDENTIALS_PATH=./config/firebase-admin-sdk.json
  ```
- **Directory Structure**: Created `config/` directory for Firebase Admin SDK credentials
- **Security**: Added `.gitignore` rules for Firebase credentials

### Phase 2: Domain Layer ✅
- **User Model Enhancement** (`internal/domain/models.go`):
  - Added `FirebaseUID string` - Firebase user ID
  - Added `DisplayName string` - Full name from Firebase
  - Added `PhotoURL string` - Profile picture from Firebase
  - Added `Role UserRole` - Two-role system (student/admin)
  - Kept legacy fields for backward compatibility

- **User Role System**:
  ```go
  type UserRole string
  const (
      RoleStudent UserRole = "student"
      RoleAdmin   UserRole = "admin"
  )
  ```

- **Auth DTOs**:
  - `VerifyTokenRequest` - Firebase ID token verification
  - `VerifyTokenResponse` - Token verification result with user data
  - `FirebaseClaims` - Custom claims extraction from Firebase tokens
  - `UserProfileResponse` - Safe user profile data for API responses
  - `UpdateUserRequest` - User profile updates
  - `UpdateUserRoleRequest` - Admin-only role management

### Phase 3: Repository Layer ✅
- **UserRepository Interface** (`internal/repository/interfaces.go`):
  ```go
  type UserRepository interface {
      Create(ctx context.Context, user *domain.User) error
      GetByID(ctx context.Context, id string) (*domain.User, error)
      GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
      GetByEmail(ctx context.Context, email string) (*domain.User, error)
      GetAll(ctx context.Context, pagination *domain.PaginationRequest) ([]*domain.User, int64, error)
      Update(ctx context.Context, user *domain.User) error
      UpdateLastLogin(ctx context.Context, userID string) error
      Delete(ctx context.Context, id string) error
  }
  ```

- **In-Memory Implementation** (`internal/repository/memory_simple.go`):
  - Full CRUD operations
  - Three-way indexing: by ID, Firebase UID, and email
  - Thread-safe with sync.RWMutex
  - Pagination support

- **PostgreSQL Implementation** (`internal/repository/postgres/user.go`):
  - Full CRUD operations with SQL queries
  - Proper error handling with custom API errors
  - NULL-safe handling for optional fields
  - Username generation from email
  - Integrated into `postgres.Repositories` struct

## 🔄 Remaining Phases

### Phase 4: Service Layer (Auth Service) - NEXT
**Location**: `internal/service/auth.go`

**Required Implementation**:
```go
type AuthService interface {
    // Initialize Firebase Admin SDK
    InitializeFirebase(ctx context.Context) error

    // Verify Firebase ID token and return user claims
    VerifyFirebaseToken(ctx context.Context, idToken string) (*domain.FirebaseClaims, error)

    // Get or create user from Firebase token (sync to backend)
    GetOrCreateUser(ctx context.Context, firebaseUID, email, displayName, photoURL string) (*domain.User, error)

    // Get user profile
    GetUserProfile(ctx context.Context, userID string) (*domain.UserProfileResponse, error)

    // Update user role (admin only)
    UpdateUserRole(ctx context.Context, userID string, role domain.UserRole) error

    // Update last login timestamp
    UpdateLastLogin(ctx context.Context, userID string) error
}
```

**Key Features**:
- Firebase Admin SDK initialization from credentials file
- Token verification extracting UID, email, name, picture
- User synchronization: create if new, update if exists
- Role management for admin operations

### Phase 5: Middleware Layer (JWT Verification)
**Location**: `internal/middleware/auth.go`

**Required Middleware**:
```go
// AuthRequired verifies Firebase token and adds user to context
func AuthRequired(authService service.AuthService) Middleware

// AdminRequired requires authenticated admin user
func AdminRequired(authService service.AuthService) Middleware

// GetUserFromContext retrieves user from request context
func GetUserFromContext(ctx context.Context) *domain.User
```

**Flow**:
1. Extract `Authorization: Bearer <token>` header
2. Verify Firebase ID token via AuthService
3. Get/create user in backend database
4. Add user to request context
5. Continue to handler

### Phase 6: Handler Layer (Auth & Admin Endpoints)
**Location**: `internal/handler/auth.go`, `internal/handler/admin.go`

**Auth Endpoints**:
```
POST /api/v1/auth/verify    - Verify Firebase token, sync user to backend
GET  /api/v1/auth/me        - Get current user profile
PUT  /api/v1/auth/me        - Update user profile (display name, photo)
```

**Admin Endpoints**:
```
GET    /api/v1/admin/users           - List all users (paginated)
GET    /api/v1/admin/users/{id}      - Get user details
PUT    /api/v1/admin/users/{id}/role - Update user role (admin only)
DELETE /api/v1/admin/users/{id}      - Delete user (admin only)
```

**Protected Routes** - Apply `AuthRequired()` middleware to:
- `POST /api/v1/exercises/{id}/submit`
- `GET /api/v1/progress/{userId}`
- `POST /api/v1/progress/{userId}/lesson/{lessonId}`

**Admin Routes** - Apply `AdminRequired()` middleware to:
- `POST /api/v1/courses`
- `PUT /api/v1/courses/{id}`
- `DELETE /api/v1/courses/{id}`
- All admin endpoints

### Phase 7: Database Migration (Users Table)
**Location**: `internal/repository/postgres/migrations/003_create_users_table.sql`

**SQL Migration**:
```sql
CREATE TABLE IF NOT EXISTS gopro.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firebase_uid VARCHAR(128) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    photo_url TEXT,
    password_hash TEXT, -- For legacy support
    role VARCHAR(20) NOT NULL DEFAULT 'student',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_firebase_uid ON gopro.users(firebase_uid);
CREATE INDEX idx_users_email ON gopro.users(email);
CREATE INDEX idx_users_role ON gopro.users(role) WHERE is_active = TRUE;
```

**Migration Integration**:
- Add to `postgres.RunMigrations()` in `repositories.go`
- Run automatically on server startup or via CLI command

### Phase 8: Frontend Integration Updates
**Location**: `frontend/src/lib/api.ts`, `frontend/src/contexts/auth-context.tsx`

**API Client Updates**:
```typescript
// Get Firebase ID token from auth context
const idToken = await user.getIdToken();

// Add to all authenticated requests
headers: {
  'Authorization': `Bearer ${idToken}`,
  'Content-Type': 'application/json'
}

// Handle 401 responses (expired token)
if (response.status === 401) {
  // Refresh token or redirect to login
}
```

**Auth Context Updates**:
```typescript
// After Firebase login, sync with backend
const verifyAndSync = async (user: User) => {
  const idToken = await user.getIdToken();
  const response = await api.post('/auth/verify', { id_token: idToken });
  // Store backend user ID for progress tracking
  setBackendUser(response.data.user);
};
```

**Progress Sync Strategy**:
- Use backend API as source of truth
- Maintain Firestore as cache for offline support
- Sync on login and periodically

### Phase 9: Testing & Documentation
**Testing**:
- Unit tests for AuthService
- Unit tests for UserRepository (memory & PostgreSQL)
- Integration tests for auth endpoints
- Middleware tests for token verification
- E2E tests for complete auth flow

**Documentation**:
- API documentation for auth endpoints
- Firebase setup guide
- Environment configuration guide
- Deployment checklist
- Security best practices

## Setup Instructions

### Backend Setup
1. **Install Dependencies**:
   ```bash
   cd backend
   go mod tidy
   ```

2. **Configure Firebase Admin SDK**:
   - Go to Firebase Console → Project Settings → Service Accounts
   - Click "Generate New Private Key"
   - Save as `backend/config/firebase-admin-sdk.json`

3. **Update Environment**:
   ```bash
   cp .env.example .env
   # Edit .env with your Firebase project ID
   FIREBASE_PROJECT_ID=your-actual-project-id
   FIREBASE_CREDENTIALS_PATH=./config/firebase-admin-sdk.json
   ```

4. **Run Database Migration** (after Phase 7):
   ```bash
   go run ./cmd/server --migrate
   ```

### Frontend Setup (Already Done ✅)
Frontend Firebase Auth is already fully implemented:
- Email/Password authentication
- Google OAuth
- GitHub OAuth (requires GitHub OAuth app setup)
- User profile management
- Progress tracking in Firestore

**GitHub OAuth Setup**:
1. Create OAuth app at github.com/settings/developers
2. Add Client ID/Secret to Firebase Console → Authentication → Sign-in method → GitHub

## Architecture Decisions

### Why Hybrid Storage (Firebase + PostgreSQL)?
- **Firebase Auth**: Handles authentication complexity (OAuth, email verification, password reset)
- **PostgreSQL**: Stores application data (user metadata, progress, courses)
- **Benefits**: Best of both worlds - robust auth + full backend control

### Why Firebase Admin SDK (Not Just JWT)?
- **User Management**: Can verify, revoke tokens, manage users from backend
- **Token Verification**: Official, secure method with automatic key rotation
- **Future Features**: Easy to add custom claims, user deletion, etc.

### Why Two Roles Only?
- **Simplicity**: Student vs Admin covers 99% of use cases
- **Scalable**: Easy to extend to multi-role system later if needed
- **Clear Permissions**: Easier to reason about who can do what

## Security Considerations

### Token Security
- **HTTPS Only**: Never send Firebase tokens over HTTP
- **Token Expiration**: Firebase tokens valid for 1 hour
- **Refresh Strategy**: Frontend should refresh tokens proactively

### Credentials Security
- **Never Commit**: Firebase Admin SDK credentials never in version control
- **Environment Variables**: Use secret management in production
- **Access Control**: Limit Firebase project permissions

### Database Security
- **Indexed Lookups**: Fast queries by Firebase UID and email
- **No Password Storage**: Firebase handles passwords, backend never sees them
- **Role Validation**: Always verify user role before admin operations

## Next Steps (Priority Order)

1. **Implement AuthService** (Phase 4) - Core authentication logic
2. **Create Auth Middleware** (Phase 5) - Protect routes
3. **Add Auth Endpoints** (Phase 6) - API for token verification
4. **Run Database Migration** (Phase 7) - Create users table
5. **Update Frontend** (Phase 8) - Send tokens to backend
6. **Write Tests** (Phase 9) - Ensure everything works

## Questions & Decisions Needed

1. **First User as Admin**: Should the first registered user automatically become admin?
   - **Recommendation**: Yes, or use environment variable `ADMIN_EMAILS` for initial admins

2. **Email Verification**: Require verified email before backend access?
   - **Recommendation**: Yes for production, optional for development

3. **Anonymous Auth**: Allow anonymous users to try lessons?
   - **Recommendation**: Yes, convert to registered user later

4. **Session Management**: Backend sessions or stateless JWT only?
   - **Recommendation**: Stateless JWT (Firebase tokens) - simpler, scalable

5. **Progress Migration**: Migrate existing Firestore progress to backend?
   - **Recommendation**: Yes, write migration script to sync on first login

## Useful Commands

```bash
# Build backend
cd backend && go build ./cmd/server

# Run backend (after all phases)
./server

# Run tests
go test ./...

# Run with hot reload (after Phase 4+)
air

# Check for linting issues
golangci-lint run

# Run security scan
gosec ./...
```

## Resources

- **Firebase Admin Go SDK**: https://firebase.google.com/docs/admin/setup
- **Firebase Auth**: https://firebase.google.com/docs/auth
- **Go PostgreSQL Driver**: https://github.com/lib/pq
- **Issue Tracker**: https://github.com/DimaJoyti/go-pro/issues/65

## Contributors
- **Backend Implementation**: Claude Code
- **Frontend Implementation**: Already completed
- **Database Design**: Clean architecture with repository pattern

---

**Last Updated**: 2025-01-24
**Next Update**: After Phase 4 completion

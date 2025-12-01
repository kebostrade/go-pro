# GO-PRO Architecture Guide

Complete technical architecture of the GO-PRO learning platform.

## System Overview

GO-PRO is a **dual-purpose system**:
1. **Learning Platform**: Interactive web-based course with exercises and progress tracking
2. **Real-world Example**: Production AI agent framework for financial services

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Next.js)                        │
│         Learning Dashboard + Code Editor Interface           │
└────────────────────┬────────────────────────────────────────┘
                     │ HTTP/REST
┌────────────────────▼────────────────────────────────────────┐
│              Backend API (Go REST API)                       │
│  ├─ Course Management    ├─ Exercise Validation             │
│  ├─ Progress Tracking    ├─ User Management                 │
│  └─ Code Execution       └─ Analytics                       │
└────────────────────┬────────────────────────────────────────┘
                     │
        ┌────────────┴────────────┐
        │                         │
┌───────▼────────┐      ┌────────▼──────────┐
│  PostgreSQL    │      │  Redis Cache      │
│  (Persistence) │      │  (Sessions/Data)  │
└────────────────┘      └───────────────────┘
```

## Module Structure

### Standardized on Go 1.23

Each major component is a **separate Go module** with its own `go.mod`:

```
go-pro/
├── backend/                       # go.mod (Learning Platform)
│   ├── cmd/server/
│   ├── internal/
│   │   ├── config/
│   │   ├── domain/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── repository/
│   │   └── service/
│   └── pkg/
│
├── course/                        # go.mod (Course Content)
│   ├── lessons/
│   ├── code/
│   └── projects/
│
├── basic/                         # go.mod (Learning Examples)
│   ├── examples/
│   ├── exercises/
│   └── projects/
│
├── services/                      # Microservices (experimental)
│   ├── ai-agent-platform/         # go.mod (Production AI Framework)
│   │   ├── internal/agent/
│   │   ├── internal/llm/
│   │   ├── internal/tools/
│   │   └── pkg/types/
│   │
│   ├── api-gateway/               # go.mod (API Gateway)
│   └── shared/                    # go.mod (Shared Libraries)
│
└── frontend/                      # package.json (Next.js)
    ├── app/
    ├── components/
    ├── lib/
    └── contexts/
```

**Important**: Always `cd` to the correct module directory before running Go commands.

## Backend Architecture (Clean Architecture)

### Layer 1: Handler Layer (HTTP Interface)
```
handler/
├── auth.go          # Authentication endpoints
├── lesson.go        # Lesson management
├── exercise.go      # Exercise submission
└── progress.go      # Progress tracking
```

**Responsibility**: Parse HTTP requests, call services, return HTTP responses

**Example**:
```go
// handler/exercise.go
func (h *Handler) SubmitExercise(w http.ResponseWriter, r *http.Request) {
    // Parse request
    // Call service
    // Return response
}
```

### Layer 2: Service Layer (Business Logic)
```
service/
├── lesson.go        # Lesson operations
├── exercise.go      # Exercise logic
└── progress.go      # Progress calculation
```

**Responsibility**: Implement business rules independent of HTTP

**Example**:
```go
// service/exercise.go
type ExerciseService struct {
    repo ExerciseRepository
}

func (s *ExerciseService) Submit(ctx context.Context, sub Submission) error {
    // Validate
    // Test code
    // Record progress
    // Return result
}
```

### Layer 3: Repository Layer (Data Access)
```
repository/
├── interfaces.go       # All repository contracts
├── memory_simple.go    # In-memory implementation
└── postgres/           # PostgreSQL implementation
    ├── exercise.go
    ├── lesson.go
    └── migrations/
```

**Responsibility**: Data persistence operations

**Key Pattern**: Repository interfaces allow switching implementations:

```go
// repository/interfaces.go
type ExerciseRepository interface {
    Save(ctx context.Context, exercise *Exercise) error
    FindByID(ctx context.Context, id string) (*Exercise, error)
    FindByLesson(ctx context.Context, lessonID string) ([]Exercise, error)
}

// Different implementations can be used:
// - memory_simple.go: In-memory (tests, development)
// - postgres/exercise.go: PostgreSQL (production)
```

### Layer 4: Domain Layer (Business Entities)
```
domain/
├── models.go        # Domain entities
└── errors.go        # Domain-specific errors
```

**Responsibility**: Define business entities and rules

```go
// domain/models.go
type Exercise struct {
    ID        string
    Title     string
    Content   string
    TestCode  string
    Difficulty Level
}

type Submission struct {
    ExerciseID string
    UserID     string
    Code       string
    Status     SubmissionStatus
    Score      int
}
```

## Frontend Architecture (Next.js 15)

### App Router Structure
```
frontend/
├── app/                   # Next.js 15 App Router
│   ├── layout.tsx        # Root layout
│   ├── page.tsx          # Home page
│   ├── courses/          # Course pages
│   │   └── [courseId]/lesson/[lessonId]/page.tsx
│   ├── exercises/        # Exercise pages
│   │   └── [exerciseId]/page.tsx
│   └── dashboard/        # User dashboard
│
├── components/           # Reusable components
│   ├── CodeEditor.tsx
│   ├── LessonCard.tsx
│   ├── ProgressBar.tsx
│   └── ui/              # Base UI components
│
├── contexts/            # React contexts
│   ├── auth-context.tsx
│   └── theme-context.tsx
│
├── lib/                 # Utilities
│   ├── api.ts          # API client
│   └── auth.ts         # Auth utilities
│
└── styles/             # CSS & Tailwind
```

### Key Components

**CodeEditor** - Monaco Editor wrapper
```tsx
<CodeEditor
  language="go"
  value={code}
  onChange={setCode}
  theme="vs-dark"
/>
```

**ProgressTracker** - Real-time progress updates
```tsx
<ProgressTracker
  lessonId={lessonId}
  userId={userId}
  onComplete={handleComplete}
/>
```

### State Management

Uses **React Context** for global state:
```tsx
// contexts/auth-context.tsx
<AuthProvider>
  <App />
</AuthProvider>

// Access anywhere:
const { user, login, logout } = useAuth()
```

## Data Flow

### Exercise Submission Flow

```
1. User submits code in CodeEditor
   ↓
2. Frontend calls: POST /api/v1/exercises/{id}/submit
   ↓
3. Backend Handler receives request
   ├─ Validates submission
   ├─ Calls ExerciseService.Submit()
   ├─ Service validates code syntax
   ├─ Service compiles and runs tests
   ├─ Records result in Repository
   └─ Returns response
   ↓
4. Frontend receives score/feedback
   ↓
5. Updates UI and progress tracking
```

### Progress Tracking Flow

```
1. Service saves exercise result
   ↓
2. Progress Service calculates metrics:
   - Lessons completed
   - Exercises solved
   - Average score
   - Time spent
   ↓
3. Frontend fetches progress:
   GET /api/v1/progress/{userId}
   ↓
4. Dashboard displays charts and stats
```

## Database Schema (PostgreSQL)

### Core Tables

```sql
-- Users
users (id, email, name, created_at)

-- Courses
courses (id, title, description, level)

-- Lessons
lessons (id, course_id, title, content, order)

-- Exercises
exercises (id, lesson_id, title, content, test_code, difficulty)

-- Submissions
submissions (id, user_id, exercise_id, code, status, score, submitted_at)

-- Progress
progress (id, user_id, lesson_id, completed_at, score)
```

### Relationships

```
User 1--M Submission
User 1--M Progress
Course 1--M Lesson
Lesson 1--M Exercise
Exercise 1--M Submission
```

## API Structure

### RESTful Endpoints

```
GET  /api/v1/health              # Health check
GET  /api/v1/courses             # List courses
GET  /api/v1/courses/{id}        # Get course
GET  /api/v1/courses/{id}/lessons # Lessons in course
GET  /api/v1/exercises/{id}      # Get exercise
POST /api/v1/exercises/{id}/submit # Submit solution
GET  /api/v1/progress/{userId}   # Get progress
POST /api/v1/progress/{userId}/lesson/{lessonId} # Update progress
```

### Response Format

```json
{
  "success": true,
  "data": {
    // Response payload
  },
  "error": null,
  "timestamp": "2025-01-01T00:00:00Z"
}
```

### Error Handling

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": "INVALID_CODE",
    "message": "Syntax error in submitted code",
    "details": "..."
  }
}
```

## Middleware Pipeline

```
Request
   ↓
1. Authentication Middleware
   └─ Validates JWT token
   ↓
2. CORS Middleware
   └─ Validates origin
   ↓
3. Logging Middleware
   └─ Records request details
   ↓
4. Rate Limiting Middleware
   └─ Prevents abuse
   ↓
Handler → Service → Repository → Database
   ↓
Response
```

## Code Execution Strategy

### Safe Code Execution

Exercise submissions are executed in **isolated environment**:

```
User Code → Compile → Run Tests → Capture Output → Return Result
```

**Security measures**:
- Code compiled before execution
- Tests run with timeout limits
- Output captured and sanitized
- Errors caught and reported

## Testing Architecture

### Backend Testing

```
tests/
├── unit/        # Service logic tests
├── integration/ # API endpoint tests
└── e2e/        # Full workflow tests
```

**Patterns**:
- Table-driven tests
- Mock repositories
- Test fixtures

```go
func TestExerciseSubmit(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    int // expected score
        wantErr bool
    }{
        {"valid solution", "correct code", 100, false},
        {"invalid syntax", "bad code", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

### Frontend Testing

```
tests/
├── components/  # Component tests
├── pages/       # Page tests
└── utils/       # Utility tests
```

Uses **Jest** and **React Testing Library**

## Performance Considerations

### Caching Strategy

- **Redis**: Cache lesson content, exercise metadata
- **In-Memory**: Cache user sessions
- **CDN**: Static assets (frontend)

### Database Optimization

- Indexes on frequently queried columns
- Connection pooling
- Query optimization

### Code Execution Optimization

- Compile once, run multiple test cases
- Timeout limits prevent infinite loops
- Resource limits prevent resource exhaustion

## Security Architecture

### Authentication

- JWT tokens for API authentication
- Token refresh mechanism
- Secure token storage

### Authorization

- Role-based access control (RBAC)
- User can only access own progress
- Admin endpoints protected

### Input Validation

- Validate all user inputs
- Sanitize code submissions
- SQL injection prevention via parameterized queries

### Data Protection

- Passwords hashed with bcrypt
- Sensitive data encrypted at rest
- HTTPS for all communication

## Deployment Architecture

### Backend Deployment

```
Go Application
  ↓
Docker Container
  ↓
Kubernetes Cluster (or Cloud Run)
  ↓
PostgreSQL (Cloud SQL / RDS)
```

### Frontend Deployment

```
Next.js Application
  ↓
Build → Static files + API routes
  ↓
Vercel / Netlify / Static host
  ↓
CDN for assets
```

### Environment Parity

- Development: Local with in-memory DB
- Staging: Cloud with test DB
- Production: Cloud with encrypted data

## Extension Points

### Add New Module

1. Create module directory: `cd go-pro && mkdir services/my-service`
2. Initialize: `cd services/my-service && go mod init`
3. Implement service
4. Register in API gateway

### Add New Lesson

1. Create directory: `course/code/lesson-XX/`
2. Add content: `lessons/lesson-XX/README.md`
3. Add exercises: `code/lesson-XX/exercises/`
4. Add tests: `code/lesson-XX/*.test.go`
5. Update syllabus: `course/syllabus.md`

### Add New Feature to API

1. Create handler: `backend/internal/handler/feature.go`
2. Create service: `backend/internal/service/feature.go`
3. Create repository: `backend/internal/repository/feature.go`
4. Create domain: Update `backend/internal/domain/models.go`
5. Register route in router setup
6. Add tests at each layer

## Related Documentation

- [Module Guide](MODULE_GUIDE.md) - Working with individual modules
- [Backend API](../backend/README.md) - API reference
- [Testing Guide](TESTING_GUIDE.md) - Testing strategies
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues

---

**Architecture evolved from**: Clean Architecture (Robert C. Martin), Domain-Driven Design, RESTful API design

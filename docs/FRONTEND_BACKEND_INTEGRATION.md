# Frontend-Backend Integration Guide

## Overview

This guide explains how the GO-PRO frontend (Next.js) integrates with the backend API (Go).

## Architecture

```
┌─────────────────┐         ┌─────────────────┐
│   Frontend      │         │    Backend      │
│   (Next.js)     │ ◄─────► │    (Go API)     │
│   Port: 3000    │  HTTP   │   Port: 8080    │
└─────────────────┘         └─────────────────┘
```

## Backend API Endpoints

The backend provides the following REST API endpoints:

### Health & Status
- `GET /api/v1/health` - Health check

### Course Management
- `GET /api/v1/courses` - List all courses (paginated)
- `GET /api/v1/courses/{id}` - Get course details
- `POST /api/v1/courses` - Create new course (admin)
- `PUT /api/v1/courses/{id}` - Update course (admin)
- `DELETE /api/v1/courses/{id}` - Delete course (admin)

### Lessons
- `GET /api/v1/courses/{courseId}/lessons` - Get course lessons
- `GET /api/v1/lessons/{id}` - Get lesson details
- `POST /api/v1/lessons` - Create lesson (admin)
- `PUT /api/v1/lessons/{id}` - Update lesson (admin)
- `DELETE /api/v1/lessons/{id}` - Delete lesson (admin)

### Exercises
- `GET /api/v1/exercises/{id}` - Get exercise details
- `POST /api/v1/exercises/{id}/submit` - Submit exercise solution

### Progress Tracking
- `GET /api/v1/progress/{userId}` - Get user progress
- `POST /api/v1/progress/{userId}/lesson/{lessonId}` - Update lesson progress

### Curriculum
- `GET /api/v1/curriculum` - Get full curriculum structure
- `GET /api/v1/curriculum/lesson/{id}` - Get detailed lesson information

## Frontend API Client

The frontend uses a centralized API client located at `frontend/src/lib/api.ts`:

```typescript
import { api } from '@/lib/api';

// Example: Fetch curriculum
const curriculum = await api.getCurriculum();

// Example: Get lesson details
const lesson = await api.getLessonDetail(1);

// Example: Update progress
await api.updateProgress('user-123', 'lesson-1', {
  completed: true,
  score: 95
});
```

## Environment Configuration

### Backend (.env)
```bash
# Server
PORT=8080
HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=gopro
DB_USER=postgres
DB_PASSWORD=your-password

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
```

### Frontend (.env.local)
```bash
# Backend API
NEXT_PUBLIC_API_URL=http://localhost:8080

# Firebase (for auth)
NEXT_PUBLIC_FIREBASE_API_KEY=your-key
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=your-domain
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id
```

## Running the Full Stack

### Option 1: Development Mode

**Terminal 1 - Backend:**
```bash
cd backend
make dev
# Backend runs on http://localhost:8080
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
# Frontend runs on http://localhost:3000
```

### Option 2: Docker Compose

```bash
make docker-dev
# Backend: http://localhost:8080
# Frontend: http://localhost:3000
# PostgreSQL: localhost:5432
# Redis: localhost:6379
```

## API Response Format

All API responses follow this structure:

```json
{
  "success": true,
  "data": { ... },
  "message": "Operation successful",
  "request_id": "uuid",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

Error responses:

```json
{
  "success": false,
  "error": {
    "type": "validation_error",
    "message": "Invalid input",
    "details": {
      "field": "error message"
    }
  },
  "request_id": "uuid",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

## Testing the Integration

1. **Start the backend:**
   ```bash
   cd backend && make dev
   ```

2. **Test API endpoints:**
   ```bash
   curl http://localhost:8080/api/v1/health
   curl http://localhost:8080/api/v1/courses
   ```

3. **Start the frontend:**
   ```bash
   cd frontend && npm run dev
   ```

4. **Access the application:**
   - Frontend: http://localhost:3000
   - Backend API Docs: http://localhost:8080

## Troubleshooting

### CORS Issues
If you see CORS errors, ensure the backend's `CORS_ALLOWED_ORIGINS` includes your frontend URL.

### Connection Refused
- Verify backend is running on port 8080
- Check `NEXT_PUBLIC_API_URL` in frontend `.env.local`

### 404 Errors
- Ensure you're using the correct API endpoint paths
- Check backend logs for routing issues

## Next Steps

- [ ] Set up authentication with Firebase
- [ ] Implement real-time progress updates
- [ ] Add WebSocket support for live coding
- [ ] Configure production deployment


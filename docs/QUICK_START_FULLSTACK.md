# GO-PRO Full Stack Quick Start Guide

## 🚀 Get Started in 5 Minutes

This guide will help you get the complete GO-PRO learning platform running locally.

## Prerequisites

- **Go 1.23** - [Download](https://go.dev/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **Docker & Docker Compose** (optional) - [Download](https://www.docker.com/)
- **PostgreSQL** (if not using Docker)
- **Redis** (if not using Docker)

## Option 1: Quick Start with Docker (Recommended)

### 1. Clone and Setup

```bash
git clone https://github.com/DimaJoyti/go-pro.git
cd go-pro
```

### 2. Start Everything

```bash
make docker-dev
```

This starts:
- ✅ Backend API on http://localhost:8080
- ✅ Frontend on http://localhost:3000
- ✅ PostgreSQL on localhost:5432
- ✅ Redis on localhost:6379
- ✅ Adminer (DB UI) on http://localhost:8081

### 3. Access the Platform

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **API Docs**: http://localhost:8080

## Option 2: Manual Setup

### 1. Backend Setup

```bash
# Navigate to backend
cd backend

# Copy environment file
cp .env.example .env
# Edit .env with your database credentials

# Install dependencies
go mod download

# Run database migrations (if needed)
# psql -U postgres -d gopro -f scripts/init-db.sql

# Start backend
make dev
# Backend runs on http://localhost:8080
```

### 2. Frontend Setup (New Terminal)

```bash
# Navigate to frontend
cd frontend

# Install dependencies
bun install

# Copy environment file
cp .env.example .env.local
# Edit .env.local:
# NEXT_PUBLIC_API_URL=http://localhost:8080

# Start frontend
bun run dev
# Frontend runs on http://localhost:3000
```

### 3. Verify Setup

```bash
# Test backend
curl http://localhost:8080/api/v1/health

# Test frontend
open http://localhost:3000
```

## 🎯 What You Get

### Backend Features
- ✅ RESTful API with Go 1.23
- ✅ Course management
- ✅ Lesson system
- ✅ Exercise submission
- ✅ Progress tracking
- ✅ PostgreSQL database
- ✅ Redis caching

### Frontend Features
- ✅ Next.js 15 with React 19
- ✅ TypeScript
- ✅ Tailwind CSS
- ✅ Interactive code editor (Monaco)
- ✅ Firebase authentication
- ✅ Responsive design
- ✅ Real-time progress tracking

## 📚 Available Endpoints

### Backend API

```bash
# Health check
GET /api/v1/health

# Courses
GET /api/v1/courses
GET /api/v1/courses/{id}
GET /api/v1/courses/{courseId}/lessons

# Lessons
GET /api/v1/lessons/{id}

# Exercises
GET /api/v1/exercises/{id}
POST /api/v1/exercises/{id}/submit

# Progress
GET /api/v1/progress/{userId}
POST /api/v1/progress/{userId}/lesson/{lessonId}

# Curriculum
GET /api/v1/curriculum
GET /api/v1/curriculum/lesson/{id}
```

## 🧪 Testing

### Test Backend

```bash
cd backend
go test ./...
```

### Test Frontend

```bash
cd frontend
bun run lint
bun run build
```

### Test Projects

```bash
# Test all examples
./scripts/test-all-examples.sh

# Test specific project
cd basic/projects/design-patterns
go test ./...
```

## 🛠️ Development Workflow

### Backend Development

```bash
cd backend
make dev          # Run with hot reload
make test         # Run tests
make lint         # Run linter
make build        # Build binary
```

### Frontend Development

```bash
cd frontend
bun run dev       # Development server
bun run build     # Production build
bun run lint      # Lint code
```

## 📖 Learning Path

1. **Start with Basics**: `basic/examples/`
2. **Practice Exercises**: `basic/exercises/`
3. **Build Projects**: `basic/projects/`
4. **Advanced Topics**: `advanced/`
5. **Real-world Apps**: `services/`

## 🔧 Troubleshooting

### Backend won't start
- Check PostgreSQL is running
- Verify `.env` configuration
- Check port 8080 is available

### Frontend won't start
- Run `bun install` again
- Check `.env.local` has correct API URL
- Verify port 3000 is available

### CORS errors
- Ensure backend `CORS_ALLOWED_ORIGINS` includes `http://localhost:3000`
- Check frontend is using correct API URL

### Database connection failed
- Verify PostgreSQL is running
- Check database credentials in `.env`
- Ensure database exists: `createdb gopro`

## 📚 Documentation

- **Full Integration Guide**: [FRONTEND_BACKEND_INTEGRATION.md](./FRONTEND_BACKEND_INTEGRATION.md)
- **Development Guide**: [CLAUDE.md](./CLAUDE.md)
- **Upgrade Summary**: [UPGRADE_SUMMARY_2025.md](./UPGRADE_SUMMARY_2025.md)
- **Learning Paths**: [LEARNING_PATHS.md](./LEARNING_PATHS.md)
- **Projects Guide**: [PROJECTS.md](./PROJECTS.md)

## 🎓 Next Steps

1. ✅ Complete the setup above
2. 📖 Read the [Learning Paths](./LEARNING_PATHS.md)
3. 💻 Start with [Basic Examples](./basic/examples/)
4. 🏗️ Build your first [Project](./basic/projects/)
5. 🚀 Deploy to production

## 🤝 Need Help?

- Check [CLAUDE.md](./CLAUDE.md) for detailed development guide
- Review [FRONTEND_BACKEND_INTEGRATION.md](./FRONTEND_BACKEND_INTEGRATION.md) for API details
- See [TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md) for common issues

---

**Happy Learning! 🎉**


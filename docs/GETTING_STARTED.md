# Getting Started with GO-PRO

Complete guide for developers to get up and running with the GO-PRO learning platform.

## Prerequisites

- **Go 1.23+** - [Download](https://go.dev/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **Git** - [Download](https://git-scm.com/)
- **Make** (optional but recommended)

## Quick Start (5 minutes)

### Clone Repository
```bash
git clone https://github.com/DimaJoyti/go-pro.git
cd go-pro
```

### Option 1: Use Script (Recommended)
```bash
# Start both backend and frontend
./scripts/start-dev.sh

# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# API Docs: http://localhost:8080/api/v1
```

### Option 2: Manual Start
```bash
# Terminal 1: Start Backend
cd backend
go mod tidy
go run ./cmd/server

# Terminal 2: Start Frontend
cd frontend
npm install
npm run dev
```

## Project Structure Overview

```
go-pro/
├── backend/           # Go REST API server
├── frontend/          # Next.js learning dashboard
├── course/            # Lesson content & exercises
├── basic/             # Examples, exercises, projects
├── docs/              # Documentation
└── CLAUDE.md          # Developer guidance
```

## Verify Installation

### Backend Health Check
```bash
# Terminal 1 started backend? Try:
curl http://localhost:8080/api/v1/health

# Expected response:
# {"status":"ok"}
```

### Frontend Access
```bash
# Open in browser:
# http://localhost:3000
```

## Development Workflow

### Working with Backend

```bash
cd backend

# Run tests
go test ./...

# Run with hot reload (requires air)
make dev

# Build for production
make build

# Linting & formatting
make lint
make fmt
```

### Working with Frontend

```bash
cd frontend

# Development server
npm run dev

# Build for production
npm run build

# Run production build
npm start

# Linting
npm run lint
```

### Working with Course Content

```bash
cd course

# Read lessons
cat lessons/lesson-01/README.md

# Run exercises
cd code/lesson-01
go test ./...

# Run solutions
go run solutions/main.go
```

## Environment Configuration

### Backend (.env)

Create `backend/.env`:
```env
PORT=8080
ENV=development
DATABASE_URL=postgres://user:pass@localhost/goprodb
JWT_SECRET=your-secret-key
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### Frontend (.env.local)

Create `frontend/.env.local`:
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_ENV=development
```

## Common Commands

| Task | Command |
|------|---------|
| Start full stack | `./scripts/start-dev.sh` |
| Run backend tests | `cd backend && go test ./...` |
| Run frontend tests | `cd frontend && npm test` |
| Build backend | `cd backend && go build ./cmd/server` |
| Build frontend | `cd frontend && npm run build` |
| Run linting | `cd backend && make lint` |
| Format code | `cd backend && make fmt` |
| Check types | `cd frontend && npm run type-check` |

## Troubleshooting

### Backend won't start
```bash
# Check if port 8080 is free
lsof -i :8080

# Check Go version
go version  # Should be 1.23+

# Check dependencies
cd backend && go mod tidy && go mod verify
```

### Frontend won't start
```bash
# Check Node version
node --version  # Should be 18+

# Clear npm cache
npm cache clean --force

# Reinstall dependencies
rm -rf node_modules package-lock.json
npm install
```

### API connection errors
```bash
# Verify backend is running
curl http://localhost:8080/api/v1/health

# Check frontend env variables
cat frontend/.env.local

# Check CORS settings
grep CORS backend/.env
```

## Next Steps

1. **Read the Architecture Guide**: [ARCHITECTURE.md](ARCHITECTURE.md)
2. **Explore Module Guide**: [MODULE_GUIDE.md](MODULE_GUIDE.md)
3. **Start a Lesson**: [Course README](../course/README.md)
4. **Run the Tests**: `cd backend && go test ./...`
5. **Try a Project**: [Projects Guide](PROJECTS.md)

## Getting Help

- 📖 **Docs**: See individual module READMEs
- 🧪 **Tests**: `go test -v` for detailed output
- 🔍 **Logs**: Check terminal output for error messages
- 📚 **Examples**: Browse `basic/examples/` directory

## Development Tips

### Use Makefile Commands
```bash
cd backend
make help  # Show all available commands
```

### Run with Verbose Output
```bash
# Go tests
go test -v ./...

# Frontend with debug logging
NEXT_DEBUG=* npm run dev
```

### Hot Reload Backend
```bash
cd backend
# Install air if needed: go install github.com/cosmtrek/air@latest
air
```

### Database Management (if using PostgreSQL)
```bash
# Migrations are handled automatically
# Check migration status in backend logs
```

---

**Ready to code?** Start with [Quick Start Guide](QUICK_START_GUIDE.md) →

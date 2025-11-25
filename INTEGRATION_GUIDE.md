# GO-PRO Integration Guide

Complete guide for running the GO-PRO full-stack application with backend and frontend integration.

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18 or higher
- **npm** or **yarn**

### One-Command Startup

```bash
./scripts/start-dev.sh
```

This script will:
1. Check and create environment files if needed
2. Build and start the backend server
3. Install frontend dependencies (if needed)
4. Start the frontend development server
5. Display server URLs and status

### Manual Startup

If you prefer to start services individually:

#### 1. Start Backend

```bash
cd backend
go run ./cmd/server
```

Backend will be available at: `http://localhost:8080`

#### 2. Start Frontend

```bash
cd frontend
npm install  # First time only
npm run dev
```

Frontend will be available at: `http://localhost:3000`

## 🔧 Configuration

### Backend Configuration

The backend uses environment variables defined in `backend/.env`:

```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001,http://127.0.0.1:3000

# Logging
LOG_LEVEL=debug
LOG_FORMAT=json

# Development Mode
DEV_MODE=true
```

**Key Settings:**
- `SERVER_PORT`: Backend server port (default: 8080)
- `CORS_ALLOWED_ORIGINS`: Comma-separated list of allowed frontend origins
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

### Frontend Configuration

The frontend uses environment variables defined in `frontend/.env.local`:

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1

# Environment
NODE_ENV=development
NEXT_PUBLIC_ENV=development
```

**Key Settings:**
- `NEXT_PUBLIC_API_URL`: Backend API base URL

## 🧪 Testing Integration

### Run Integration Tests

```bash
./scripts/test-integration.sh
```

This will test:
- ✅ Backend server health
- ✅ Frontend server availability
- ✅ API endpoints functionality
- ✅ CORS configuration
- ✅ Backend-Frontend connectivity

### Manual Testing

#### Test Backend API

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get courses
curl http://localhost:8080/api/v1/courses

# Test CORS
curl -H "Origin: http://localhost:3000" \
     -H "Access-Control-Request-Method: GET" \
     -X OPTIONS \
     http://localhost:8080/api/v1/courses -v
```

#### Test Frontend

Open your browser and navigate to:
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080/api/v1

## 📡 API Endpoints

### Backend API (http://localhost:8080/api/v1)

#### Health & Status
- `GET /health` - Health check endpoint

#### Courses
- `GET /courses` - List all courses
- `GET /courses/{id}` - Get course by ID
- `POST /courses` - Create new course
- `PUT /courses/{id}` - Update course
- `DELETE /courses/{id}` - Delete course

#### Lessons
- `GET /courses/{courseId}/lessons` - Get course lessons
- `GET /lessons/{id}` - Get lesson by ID
- `POST /lessons` - Create new lesson
- `PUT /lessons/{id}` - Update lesson
- `DELETE /lessons/{id}` - Delete lesson

#### Progress
- `GET /progress/{userId}` - Get user progress
- `POST /progress` - Update user progress

## 🔐 CORS Configuration

The backend is configured to allow requests from the frontend:

**Allowed Origins:**
- `http://localhost:3000` (default frontend)
- `http://localhost:3001` (alternative port)
- `http://127.0.0.1:3000` (localhost alternative)

**Allowed Methods:**
- GET, POST, PUT, PATCH, DELETE, OPTIONS

**Allowed Headers:**
- Accept, Content-Type, Authorization, X-Request-ID, X-CSRF-Token

**Exposed Headers:**
- X-Request-ID, X-Total-Count, X-Page, X-Page-Size

## 🗄️ Data Storage

### Development Mode

By default, the application runs with **in-memory storage**:
- No database required
- Data is lost when server restarts
- Perfect for development and testing

### Production Mode (Optional)

For persistent storage, you can configure:

#### PostgreSQL

```env
# backend/.env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=gopro_user
DB_PASSWORD=gopro_password
DB_NAME=gopro_dev
```

#### Redis (Optional - for caching)

```env
# backend/.env
REDIS_HOST=localhost
REDIS_PORT=6379
```

## 🐛 Troubleshooting

### Backend Issues

**Port already in use:**
```bash
# Find process using port 8080
lsof -i :8080
# Kill the process
kill -9 <PID>
```

**CORS errors:**
- Check `CORS_ALLOWED_ORIGINS` in `backend/.env`
- Ensure frontend origin is included
- Restart backend after changes

**Build errors:**
```bash
cd backend
go mod tidy
go build ./cmd/server
```

### Frontend Issues

**Port already in use:**
```bash
# Find process using port 3000
lsof -i :3000
# Kill the process
kill -9 <PID>
```

**API connection errors:**
- Verify `NEXT_PUBLIC_API_URL` in `frontend/.env.local`
- Ensure backend is running
- Check browser console for CORS errors

**Dependency issues:**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Common Issues

**Cannot connect to backend:**
1. Verify backend is running: `curl http://localhost:8080/api/v1/health`
2. Check backend logs: `logs/backend.log`
3. Verify CORS configuration

**Frontend not loading:**
1. Check frontend logs: `logs/frontend.log`
2. Clear browser cache
3. Try incognito/private mode

## 📊 Monitoring

### View Logs

```bash
# Backend logs
tail -f logs/backend.log

# Frontend logs
tail -f logs/frontend.log
```

### Health Checks

```bash
# Backend health
curl http://localhost:8080/api/v1/health | jq

# Frontend health
curl http://localhost:3000
```

## 🚢 Production Deployment

For production deployment, see:
- [Backend Deployment Guide](backend/README.md)
- [Frontend Deployment Guide](frontend/README.md)
- [Docker Deployment](docker/README.md)
- [Kubernetes Deployment](k8s/README.md)

## 📚 Additional Resources

- [Backend Documentation](backend/README.md)
- [Frontend Documentation](frontend/README.md)
- [API Documentation](docs/API.md)
- [Development Guide](docs/DEVELOPMENT.md)

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines.

## 📝 License

This project is licensed under the MIT License.


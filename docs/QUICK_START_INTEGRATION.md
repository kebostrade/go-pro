# 🚀 GO-PRO Quick Start - Full Stack Integration

Get the complete GO-PRO learning platform running in under 2 minutes!

## ⚡ Super Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/DimaJoyti/go-pro.git
cd go-pro

# 2. Start everything
./scripts/start-dev.sh

# 3. Open your browser
# Frontend: http://localhost:3000
# Backend API: http://localhost:8080/api/v1
```

That's it! 🎉

## 🧪 Verify Everything Works

```bash
# Run integration tests
./scripts/test-integration.sh
```

Expected output:
```
✓ Backend Server: Running at http://localhost:8080
✓ Frontend Server: Running at http://localhost:3000
✓ API Endpoints: Working correctly
✓ CORS: Configured properly
✓ Integration: Backend and Frontend connected

All integration tests passed! 🎉
```

## 📦 What's Running?

### Backend (Port 8080)
- **Health Check:** http://localhost:8080/api/v1/health
- **Courses API:** http://localhost:8080/api/v1/courses
- **Lessons API:** http://localhost:8080/api/v1/lessons
- **Progress API:** http://localhost:8080/api/v1/progress

### Frontend (Port 3000)
- **Home Page:** http://localhost:3000
- **Dashboard:** http://localhost:3000/dashboard
- **Courses:** http://localhost:3000/courses
- **Practice:** http://localhost:3000/practice

## 🛠️ Alternative: Manual Start

### Terminal 1 - Backend
```bash
cd backend
go run ./cmd/server
```

### Terminal 2 - Frontend
```bash
cd frontend
npm install  # First time only
npm run dev
```

## 🎯 Using Make Commands

```bash
# Start development environment
make start-dev

# Run integration tests
make test-integration

# Build backend
make build

# Run backend tests
make test
```

## 📝 Configuration Files

All configuration is automatic! But if you need to customize:

### Backend: `backend/.env`
```env
SERVER_PORT=8080
CORS_ALLOWED_ORIGINS=http://localhost:3000
LOG_LEVEL=debug
```

### Frontend: `frontend/.env.local`
```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## 🔍 Test API Endpoints

```bash
# Health check
curl http://localhost:8080/api/v1/health | jq

# Get all courses
curl http://localhost:8080/api/v1/courses | jq

# Test CORS
curl -H "Origin: http://localhost:3000" \
     -X OPTIONS \
     http://localhost:8080/api/v1/courses -v
```

## 🐛 Troubleshooting

### Port Already in Use

**Backend (8080):**
```bash
lsof -i :8080
kill -9 <PID>
```

**Frontend (3000):**
```bash
lsof -i :3000
kill -9 <PID>
```

### Backend Not Starting
```bash
cd backend
go mod tidy
go build ./cmd/server
./bin/server
```

### Frontend Not Starting
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run dev
```

### CORS Errors
1. Check `backend/.env` has `CORS_ALLOWED_ORIGINS=http://localhost:3000`
2. Restart backend server
3. Clear browser cache

## 📚 Next Steps

1. **Explore the Frontend:** http://localhost:3000
2. **Try the API:** http://localhost:8080/api/v1
3. **Read the Docs:** [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)
4. **Start Learning:** Check out the courses and lessons!

## 🎓 Learning Resources

- **Full Integration Guide:** [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md)
- **Backend Documentation:** [backend/README.md](backend/README.md)
- **Frontend Documentation:** [frontend/README.md](frontend/README.md)
- **API Documentation:** [docs/API.md](docs/API.md)

## 💡 Pro Tips

1. **Use the integration test script** to verify everything is working
2. **Check logs** in `logs/backend.log` and `logs/frontend.log`
3. **Use browser DevTools** to inspect API calls
4. **Enable debug logging** in `backend/.env` with `LOG_LEVEL=debug`

## 🎉 You're Ready!

The GO-PRO platform is now running and ready for you to start learning Go!

Happy coding! 🚀


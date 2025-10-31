## GO-PRO

A concise project prompt for a full-stack Golang learning platform.

### Overview
GO-PRO is a hands-on learning platform for mastering Go. It combines:
- A REST API backend written in Go for course content, exercises, and progress tracking
- A modern frontend (Next.js) for the learner experience
- A curated curriculum with lessons, exercises, and solutions you can run locally

### Goals
- Learn core Go concepts by writing and testing real code
- Explore idiomatic Go APIs and server patterns
- Practice with unit tests and simple API integrations
- Deploy a Next.js frontend (optionally to Cloudflare) to wrap the learning experience

## Tutorials
Working with WebSockets in Go
Working with Microservices in Go
Working with Design Patterns in Go
Working with Concurrency in Go
Web Architecture With Golang
Doker Kubernetes and Teraform with go
Kafka and RabbitMQ with go
Ethical Hacking with Go

### Repository structure (high-level)
- `backend/` — Go HTTP API (sample course data, exercises, progress)
- `frontend/` — Next.js app (React) with Cloudflare-ready config
- `course/` — Lessons, exercises, and solutions
	- `code/lesson-01/exercises` — Start here; run tests and implement
- `projects/` — Larger, end-to-end practice projects

---

## How to run locally

Prerequisites:
- Go 1.22+ (the server uses the new `http.ServeMux` method routing patterns)
- Node.js 18+ (Node 20+ recommended) and npm

### 1) Start the Go API

```
cd backend
go run .
```

What you get:
- API base: `http://localhost:8080`
- Health check: `http://localhost:8080/api/v1/health`
- Inline API docs: open `http://localhost:8080/`

Quick test (optional):
```
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/courses
```

Run backend tests (if any):
```
cd backend
go test ./...
```

### 2) Start the Frontend (Next.js)

```
cd frontend
npm install
npm run dev
```

What you get:
- Frontend dev server: `http://localhost:3000`

Optional Cloudflare preview/deploy (requires configuration in `wrangler.jsonc`):
```
npm run preview   # build + preview via OpenNext Cloudflare
npm run deploy    # build + deploy via OpenNext Cloudflare
```

---

## Learn by doing

- Open `course/lessons/lesson-01/README.md` for the first lesson overview.
- Implement exercises in `course/code/lesson-01/exercises/` and run tests:
```
cd course/code/lesson-01/exercises
go test
```
- Compare your work with `course/code/lesson-01/solutions/` after you try.

---

## Roadmap (suggested next steps)
- Persist data in a real database (Postgres, SQLite, etc.) instead of in-memory
- Add user accounts and authenticated progress tracking
- Build an exercise runner that executes user code and returns results
- Expand lessons (concurrency, interfaces, generics, testing, microservices)

## Contributing
PRs and improvements to lessons, examples, and tests are welcome. Keep changes small and include a short note on why/what changed.

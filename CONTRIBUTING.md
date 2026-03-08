# Contributing to GO-PRO Learning Platform

Thank you for your interest in contributing to GO-PRO! This document provides guidelines and instructions for contributing to the project.

## 🎯 Ways to Contribute

There are many ways you can contribute to GO-PRO:

- 📝 **Improve Documentation**: Fix typos, clarify explanations, add examples
- 🐛 **Report Bugs**: Submit detailed bug reports with reproduction steps
- ✨ **Suggest Features**: Propose new lessons, exercises, or platform features
- 💻 **Submit Code**: Fix bugs, implement features, or add new lessons
- 🧪 **Write Tests**: Improve test coverage and quality
- 🎨 **Improve UI/UX**: Enhance the frontend dashboard
- 📚 **Create Content**: Add new lessons, exercises, or projects

## 🚀 Getting Started

### 1. Fork and Clone

```bash
# Fork the repository on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/go-pro.git
cd go-pro

# Add upstream remote
git remote add upstream https://github.com/DimaJoyti/go-pro.git
```

### 2. Set Up Development Environment

```bash
# Install dependencies and set up the project
./start.sh --setup

# Or manually:
cd backend && go mod download && cd ..
cd frontend && bun install && cd ..
```

### 3. Create a Branch

```bash
# Create a new branch for your work
git checkout -b feature/your-feature-name

# Or for bug fixes:
git checkout -b fix/bug-description
```

## 📋 Development Guidelines

### Code Style

#### Go Code
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` and `goimports` for formatting
- Run `make lint` before committing
- Write meaningful comments for exported functions
- Keep functions small and focused

```bash
# Format your code
make fmt

# Run linter
make lint

# Fix linting issues automatically
make lint-fix
```

#### TypeScript/React Code
- Follow the existing code style
- Use TypeScript for type safety
- Use functional components with hooks
- Keep components small and reusable

```bash
cd frontend
bun run lint
```

### Testing

All code changes should include appropriate tests:

#### Backend Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific tests
cd backend
go test -v ./internal/handlers/...
```

#### Frontend Tests
```bash
cd frontend
bun test
```

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(lessons): add lesson 21 on advanced patterns
fix(api): correct error handling in progress endpoint
docs(readme): update installation instructions
test(handlers): add tests for course handler
```

## 🎓 Adding New Lessons

### Lesson Structure

Each lesson should follow this structure:

```
course/
├── lessons/
│   └── lesson-XX/
│       ├── README.md          # Lesson content
│       ├── objectives.md      # Learning objectives
│       └── resources.md       # Additional resources
└── code/
    └── lesson-XX/
        ├── main.go           # Runnable examples
        ├── exercises/        # Practice problems
        │   ├── exercise1.go
        │   └── exercise1_test.go
        └── solutions/        # Reference solutions
            └── exercise1.go
```

### Lesson Content Guidelines

1. **Clear Objectives**: Start with clear learning objectives
2. **Progressive Difficulty**: Build on previous lessons
3. **Practical Examples**: Include real-world examples
4. **Hands-on Exercises**: Provide coding challenges
5. **Comprehensive Tests**: Include automated tests
6. **Additional Resources**: Link to relevant documentation

### Example Lesson Template

```markdown
# Lesson XX: Topic Name

## 🎯 Learning Objectives

By the end of this lesson, you will be able to:
- Objective 1
- Objective 2
- Objective 3

## 📚 Introduction

Brief introduction to the topic...

## 💡 Key Concepts

### Concept 1
Explanation with code examples...

### Concept 2
Explanation with code examples...

## 🔨 Hands-on Practice

Try the exercises in `code/lesson-XX/exercises/`

## ✅ Summary

Key takeaways...

## 📖 Additional Resources

- [Resource 1](link)
- [Resource 2](link)
```

## 🐛 Reporting Bugs

When reporting bugs, please include:

1. **Description**: Clear description of the issue
2. **Steps to Reproduce**: Detailed steps to reproduce the bug
3. **Expected Behavior**: What you expected to happen
4. **Actual Behavior**: What actually happened
5. **Environment**: OS, Go version, Node.js version
6. **Screenshots**: If applicable

Use the bug report template when creating an issue.

## ✨ Suggesting Features

When suggesting features, please include:

1. **Problem Statement**: What problem does this solve?
2. **Proposed Solution**: How should it work?
3. **Alternatives**: Other solutions you've considered
4. **Additional Context**: Any other relevant information

## 🔍 Code Review Process

1. **Submit PR**: Create a pull request with your changes
2. **CI Checks**: Ensure all CI checks pass
3. **Review**: Wait for maintainer review
4. **Address Feedback**: Make requested changes
5. **Approval**: Get approval from maintainers
6. **Merge**: Your PR will be merged!

### PR Checklist

Before submitting a PR, ensure:

- [ ] Code follows project style guidelines
- [ ] All tests pass (`make test`)
- [ ] New tests added for new features
- [ ] Documentation updated if needed
- [ ] Commit messages follow conventions
- [ ] No merge conflicts with main branch
- [ ] PR description clearly explains changes

## 🏗️ Project Structure

Understanding the project structure:

```
go-pro/
├── backend/              # Go backend API
│   ├── cmd/             # Application entry points
│   ├── internal/        # Private application code
│   ├── pkg/             # Public libraries
│   └── test/            # Test utilities
├── frontend/            # Next.js frontend
│   ├── app/            # Next.js app directory
│   ├── components/     # React components
│   └── lib/            # Utility functions
├── course/             # Course content
│   ├── lessons/        # Lesson materials
│   ├── code/          # Exercises and solutions
│   └── projects/      # Hands-on projects
└── docs/              # Additional documentation
```

## 🤝 Community Guidelines

- Be respectful and inclusive
- Help others learn and grow
- Provide constructive feedback
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md)

## 📞 Getting Help

- **Questions**: Open a discussion on GitHub
- **Bugs**: Create an issue with the bug template
- **Features**: Create an issue with the feature template
- **Chat**: Join our community chat (coming soon)

## 📜 License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing to GO-PRO! Your efforts help make Go learning accessible to everyone. 🚀


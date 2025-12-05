# Quick Start Guide

## 🚀 Get Started in 30 Seconds

### Option 1: Interactive Runner (Recommended)
```bash
cd basic
go run cmd/runner/main.go
```

Then select:
- `1-12` for examples
- `p1` for Calculator project
- `p2` for Todo List project
- `a` to run all examples
- `q` to quit

### Option 2: Run Individual Examples
```bash
cd basic/examples/01_hello
go run main.go
```

### Option 3: Test Everything
```bash
./test-all.sh
```

---

## 📚 Learning Path

### Day 1: Basics (Examples 1-4)
```bash
cd basic
go run cmd/runner/main.go
# Select: 1, 2, 3, 4
```

**Topics:**
- Hello World
- Variables & Constants
- Functions
- Pointers

### Day 2: Data Structures (Examples 5-7)
```bash
# Select: 5, 6, 7
```

**Topics:**
- Arrays & Slices
- Control Flow
- Maps

### Day 3: Advanced Basics (Examples 8-10)
```bash
# Select: 8, 9, 10
```

**Topics:**
- Structs
- Interfaces
- Error Handling

### Day 4: Concurrency & Advanced (Examples 11-12)
```bash
# Select: 11, 12
```

**Topics:**
- Goroutines & Channels
- Generics, Reflection, Context

### Day 5: Practice Exercises
```bash
cd basic/exercises/01_basics
go run fizzbuzz.go              # Try yourself
go run fizzbuzz_solution.go     # Check solution
```

### Day 6-7: Build Projects
```bash
cd basic/projects/calculator
go run main.go

cd basic/projects/todo_list
go run main.go
```

---

## 🎯 Quick Commands

### Examples
| Command | Description |
|---------|-------------|
| `cd basic && go run cmd/runner/main.go` | Interactive runner |
| `cd basic/examples/01_hello && go run main.go` | Run specific example |
| `cd basic && ./test-examples.sh` | Test all examples |

### Exercises
| Command | Description |
|---------|-------------|
| `cd basic/exercises/01_basics && go run fizzbuzz.go` | Try exercise |
| `cd basic/exercises/01_basics && go run fizzbuzz_solution.go` | View solution |

### Projects
| Command | Description |
|---------|-------------|
| `cd basic/projects/calculator && go run main.go` | Run calculator |
| `cd basic/projects/todo_list && go run main.go` | Run todo list |

### Testing
| Command | Description |
|---------|-------------|
| `./test-all.sh` | Test everything (18 tests) |
| `cd basic && ./test-examples.sh` | Test examples only |

---

## 📖 What's Available

### ✅ 12 Examples
1. Hello World - Basic syntax
2. Variables - Data types and constants
3. Functions - Function patterns
4. Pointers - Memory management
5. Arrays & Slices - Collections
6. Control Flow - Conditionals and loops
7. Maps - Key-value storage
8. Structs - Custom types
9. Interfaces - Polymorphism
10. Errors - Error handling
11. Concurrency - Goroutines & channels
12. Advanced - Generics, reflection, JSON

### ✅ 4 Exercise Sets
- **Basic:** FizzBuzz, Reverse String
- **Intermediate:** URL Shortener
- **Advanced:** Web Crawler

### ✅ 2 Complete Projects
- **Calculator:** Interactive CLI calculator
- **Todo List:** Task management app

---

## 🧪 Verify Installation

```bash
# Test everything works
./test-all.sh

# Expected output:
# ✓ Passed: 18
# ✗ Failed: 0
# 🎉 All tests passed!
```

---

## 💡 Tips

### For Beginners
1. Start with the interactive runner
2. Read the code in each example
3. Modify examples and re-run them
4. Try exercises before looking at solutions

### For Practice
1. Complete all exercises
2. Build the projects from scratch
3. Add features to existing projects
4. Create your own projects

### For Reference
1. Use examples as code snippets
2. Check solutions when stuck
3. Refer to project code for patterns

---

## 🆘 Troubleshooting

### Example doesn't run
```bash
# Make sure you're in the right directory
cd basic/examples/01_hello
go run main.go
```

### Interactive runner not working
```bash
# Check Go is installed
go version

# Run from basic directory
cd basic
go run cmd/runner/main.go
```

### Tests failing
```bash
# Make sure all files are present
ls basic/examples/
ls basic/exercises/
ls basic/projects/

# Run test script
./test-all.sh
```

---

## 📊 Progress Tracker

Track your learning progress:

### Examples
- [ ] 01 - Hello World
- [ ] 02 - Variables
- [ ] 03 - Functions
- [ ] 04 - Pointers
- [ ] 05 - Arrays & Slices
- [ ] 06 - Control Flow
- [ ] 07 - Maps
- [ ] 08 - Structs
- [ ] 09 - Interfaces
- [ ] 10 - Errors
- [ ] 11 - Concurrency
- [ ] 12 - Advanced

### Exercises
- [ ] FizzBuzz
- [ ] Reverse String
- [ ] URL Shortener
- [ ] Web Crawler

### Projects
- [ ] Calculator
- [ ] Todo List

---

## 🎓 Next Steps

After completing all examples, exercises, and projects:

1. **Build Your Own Projects**
   - File organizer
   - HTTP server
   - REST API
   - Database app

2. **Learn More**
   - Go documentation: https://go.dev/doc/
   - Go by Example: https://gobyexample.com/
   - Effective Go: https://go.dev/doc/effective_go

3. **Practice**
   - LeetCode in Go
   - Build CLI tools
   - Contribute to open source

---

## ✨ Summary

- **18 working examples, exercises, and projects**
- **100% test coverage**
- **Interactive learning environment**
- **Complete documentation**

**Ready to start? Run:**
```bash
cd basic && go run cmd/runner/main.go
```

Happy coding! 🚀


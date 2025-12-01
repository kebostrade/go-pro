# Snake Game Tutorial - Complete Index

Welcome to the comprehensive Snake Game tutorial in Go with Ebiten! This document guides you through all available resources.

## 📂 Project Contents

### Main Files

#### 📖 Documentation Files
| File | Purpose | Best For |
|------|---------|----------|
| **README.md** | Project overview and quick start | Getting started, project overview |
| **TUTORIAL.md** | Step-by-step game development guide | Learning game development concepts |
| **QUICK_REFERENCE.md** | Cheat sheet for common tasks | Quick lookups, copy-paste code |
| **ADVANCED_FEATURES.md** | Feature implementation guides | Extending the game |
| **INDEX.md** | This file - navigation guide | Finding what you need |

#### 💻 Code Files
| File | Purpose | Complexity |
|------|---------|-----------|
| **main.go** | Complete, runnable game | Advanced (full game) |
| **step1_window.go** | Basic Ebiten window setup | Beginner |
| **step2_render.go** | Draw the snake on screen | Beginner |
| **step3_movement.go** | Make snake move automatically | Beginner |
| **step4_input.go** | Keyboard controls (arrow keys) | Intermediate |
| **step5_collision.go** | Collision detection and food | Intermediate |
| **step6_gameover.go** | Game states and UI text | Intermediate |
| **step7_difficulty.go** | Difficulty progression | Intermediate |

## 🎯 Learning Paths

### Path 1: Quick Start (15 minutes)
1. Read **README.md** - Overview
2. Run `go run main.go` - See the game
3. Play around, modify colors/speeds
4. ✅ Done! You have a working Snake game

### Path 2: Step-by-Step Learning (2-3 hours)
1. **README.md** - Understand project structure
2. **step1_window.go** - Learn Ebiten basics
3. **step2_render.go** - Learn grid system and drawing
4. **step3_movement.go** - Learn game timing
5. **step4_input.go** - Learn input handling
6. **step5_collision.go** - Learn collision detection
7. **step6_gameover.go** - Learn game state management
8. **step7_difficulty.go** - Learn difficulty scaling
9. **main.go** - See the complete solution
10. ✅ Understand complete game development

### Path 3: Deep Dive with Extensions (4-6 hours)
1. Complete Path 2 (Step-by-Step Learning)
2. **ADVANCED_FEATURES.md** - Choose features:
   - Use image assets
   - Add sound effects
   - Implement pause
   - Save high scores
   - Create game modes
   - Add obstacles
   - Create power-ups
   - Add particle effects
3. **QUICK_REFERENCE.md** - Reference as needed
4. ✅ Build an advanced version with custom features

### Path 4: Reference (On-demand)
- **QUICK_REFERENCE.md** - When you need quick answers
- **TUTORIAL.md** - When you need deep explanation
- **ADVANCED_FEATURES.md** - When you want to add features

## 📚 What Each File Teaches

### step1_window.go
**Concepts**: Game interface, window setup, basic game loop
**Key Learning**: How Ebiten games are structured
**Time**: 10 minutes
```
Main concepts:
- ebiten.Game interface
- Update() method
- Draw() method
- Layout() method
- Game loop basics
```

### step2_render.go
**Concepts**: Grid system, data structures, drawing
**Key Learning**: How to represent and draw game objects
**Time**: 15 minutes
```
Main concepts:
- Point struct
- Grid-based coordinates
- Snake as slice of Points
- Drawing colored rectangles
- Head vs body rendering
```

### step3_movement.go
**Concepts**: Game timing, tick system, direction vectors
**Key Learning**: How to make objects move smoothly
**Time**: 20 minutes
```
Main concepts:
- tickCount and tickSpeed
- Direction vectors
- Slice prepend/append operations
- Moving without affecting game speed
```

### step4_input.go
**Concepts**: Input handling, input validation, input buffering
**Key Learning**: How to handle player input responsively
**Time**: 15 minutes
```
Main concepts:
- ebiten.IsKeyPressed()
- Preventing invalid moves
- Input buffering (nextDir pattern)
- Smooth directional changes
```

### step5_collision.go
**Concepts**: Collision detection, game mechanics, state changes
**Key Learning**: How to detect collisions and implement game mechanics
**Time**: 25 minutes
```
Main concepts:
- Point equality comparison
- Wall collision detection
- Self-collision detection
- Food collision and growth
- Deterministic spawning
```

### step6_gameover.go
**Concepts**: Game states, state management, UI rendering
**Key Learning**: How to manage different game states
**Time**: 20 minutes
```
Main concepts:
- gameStarted and gameOver flags
- State-specific logic
- inpututil.IsKeyJustPressed()
- Text rendering with DebugPrintAt
- State reset
```

### step7_difficulty.go
**Concepts**: Difficulty progression, game balance
**Key Learning**: How to scale game difficulty dynamically
**Time**: 15 minutes
```
Main concepts:
- Dynamic speed adjustment
- Score-based difficulty
- Minimum/maximum difficulty bounds
- Game balance tuning
```

### main.go
**Concepts**: Complete integration of all concepts
**Key Learning**: Production-ready game structure
**Time**: Study only, don't run multiple copies
```
Main concepts:
- All previous concepts integrated
- Optimized rendering
- Complete game loop
- All features working together
```

## 🚀 Running the Game

### Option 1: Run Completed Game
```bash
cd snake-game-tutorial
go run main.go
```
Click window → Start with Space → Play!

### Option 2: Learn Step by Step
```bash
go run step1_window.go  # Just a window
go run step2_render.go  # Render snake
go run step3_movement.go # Snake moves
go run step4_input.go    # Control snake
go run step5_collision.go # Food and collision
go run step6_gameover.go  # Game states
go run step7_difficulty.go # Progressive difficulty
```

### Option 3: Build Executable
```bash
go build -o snake-game
./snake-game
```

## 🎮 Controls

- **Space**: Start game / Restart after game over
- **Arrow Keys**: Move the snake
- **R**: Restart (when game is over)

## 📖 Documentation Quick Reference

### Understanding Concepts?
→ Read **TUTORIAL.md**

### Need Code Examples?
→ Check **QUICK_REFERENCE.md**

### Want to Add Features?
→ See **ADVANCED_FEATURES.md**

### Getting Started?
→ Read **README.md**

### Need to Find Something?
→ You're reading it! (INDEX.md)

## 🧬 Code Structure Overview

```
Game Interface Implementation:
├── Update() - Game logic per frame
├── Draw(screen) - Render graphics
└── Layout() - Screen dimensions

Game State:
├── Snake data (positions)
├── Food data (position)
├── Game flags (started, over, paused)
├── Score tracking
└── Movement parameters

Game Loop:
├── Handle input
├── Update positions
├── Check collisions
├── Update UI
└── Render to screen
```

## 🎯 Key Concepts Checklist

### Beginner
- [ ] Understand game loop (Update/Draw)
- [ ] Create game structs
- [ ] Draw pixels to screen
- [ ] Understand grid system

### Intermediate
- [ ] Handle keyboard input
- [ ] Implement collision detection
- [ ] Manage game state
- [ ] Create time-based movement

### Advanced
- [ ] Scale difficulty dynamically
- [ ] Add game features
- [ ] Optimize performance
- [ ] Handle edge cases

## 🔧 Modification Ideas

### Easy (< 30 minutes)
- Change colors
- Adjust speed (`tickSpeed`)
- Change grid size
- Modify starting snake size
- Add debug text

### Medium (1-2 hours)
- Add pause functionality
- Save high scores
- Add obstacles
- Implement multiple difficulties
- Create game modes

### Hard (2-4 hours)
- Add sound effects
- Use image assets
- Create power-ups
- Implement particle effects
- Add 2-player mode

## 📊 Learning Progression

```
Beginner
  ↓
step1 (Window) - 10 min
step2 (Render) - 15 min
  ↓
Intermediate
  ↓
step3 (Movement) - 20 min
step4 (Input) - 15 min
step5 (Collision) - 25 min
step6 (UI) - 20 min
step7 (Difficulty) - 15 min
  ↓
Advanced
  ↓
main.go (Complete) - Study
ADVANCED_FEATURES.md - Extend
  ↓
Master
```

## 🎓 Topics Covered

### Core Game Development
- Game loops and frame timing
- Game state management
- Input handling
- Collision detection
- Difficulty scaling

### Go Programming
- Structs and interfaces
- Slice operations
- Error handling
- Function receivers
- Type assertions

### Ebiten Framework
- Window and screen management
- Graphics rendering
- Input handling
- Game interface implementation
- Timing and frame rate

## 🐛 Debugging Tips

### Check if it compiles
```bash
go build
```

### Check for syntax errors
```bash
go run main.go
```

### Add debug prints
```go
log.Println("Debug:", value)
```

### Draw debug boxes
```go
drawCell(screen, x, y, color.RGBA{255, 0, 0, 255})
```

## 📱 Platform Support

This game runs on:
- ✅ Windows
- ✅ macOS
- ✅ Linux
- ✅ Web (with modifications)
- ✅ Mobile (with modifications)

## 🎉 Next Steps After Tutorial

1. **Customize**: Add your own features
2. **Share**: Show it to friends or GitHub
3. **Publish**: Create an executable
4. **Challenge**: Try harder game mechanics
5. **Learn**: Study the Ebiten source code
6. **Create**: Make your own game

## 🤝 Getting Help

### Stuck on a concept?
→ Read TUTORIAL.md section on that topic

### Need API documentation?
→ Visit https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2

### Want to see examples?
→ Check https://github.com/hajimehoshi/ebiten/tree/main/examples

### Have a question?
→ Ask in the TUTORIAL.md FAQ section

## 📊 Project Statistics

```
Total Code Lines: ~1,500
Documentation Lines: ~3,000
Step Files: 7
Complete Files: 1
Concepts Covered: 20+
Features Taught: 15+
Time to Complete: 2-3 hours
Complexity Range: Beginner → Advanced
```

## ✅ Success Criteria

You've successfully completed this tutorial when you can:

- [ ] Run the complete game without errors
- [ ] Understand how each step builds on previous ones
- [ ] Modify game parameters (colors, speed, size)
- [ ] Add new features using the patterns learned
- [ ] Build and run an executable
- [ ] Explain the game loop to someone else

## 🎮 Play the Game

```bash
# Quick way to start playing
cd /path/to/snake-game-tutorial
go run main.go

# Wait for window to appear
# Press SPACE to start
# Use arrow keys to move
# Eat food to grow
# Don't hit walls or yourself!
```

---

**Happy learning! Now pick a learning path above and get started.** 🚀

Start with README.md if you're new, or jump to main.go if you want to play right away!

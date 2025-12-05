# Getting Started - Quick Setup Guide

Get the Snake game running in 5 minutes!

## ✅ Prerequisites Checklist

- [ ] Go 1.20 or later installed
- [ ] Terminal/Command prompt available
- [ ] Internet connection (for first dependency download)
- [ ] ~200MB free disk space

### Check Your Go Version

```bash
go version
# Should output: go version go1.23 (or later) ...
```

If you need to install Go, visit: https://golang.org/doc/install

## 🚀 Setup Steps (5 minutes)

### Step 1: Navigate to Project
```bash
cd /path/to/snake-game-tutorial
```

### Step 2: Download Dependencies
```bash
go mod download
```

This will download Ebiten and its dependencies (~200MB). This only needs to happen once.

### Step 3: Run the Game
```bash
go run main.go
```

**That's it!** A window should appear with the Snake game.

### Step 4: Play
- Press **SPACE** to start
- Use **Arrow Keys** to move
- Try to eat the red food
- Don't hit walls or yourself!
- Press **R** to restart after game over

## 🎮 First-Time Gameplay (2 minutes)

1. **Start**: Press SPACE
2. **Move**: Use arrow keys to move the snake
3. **Eat**: Move to the red square (food)
4. **Grow**: Each food eaten makes the snake longer
5. **Speed Up**: Game gets faster as you score
6. **Game Over**: Hit a wall or yourself
7. **Restart**: Press R to play again

## 📖 Learning the Code (Next Steps)

### Option A: Just Want to Play?
```bash
go run main.go
# Play and have fun!
```

### Option B: Want to Learn?
Read files in order:
1. **README.md** - Understand the project
2. **step1_window.go** - Learn basics
3. **step2_render.go** - Learn drawing
4. ... continue through steps ...
5. **main.go** - See complete game

### Option C: Want to Modify?
1. Open **main.go** in your favorite editor
2. Find the `const` section
3. Try changing values:

```go
// Try changing these:
cellSize   = 30    // Make cells bigger
tickSpeed  = 4     // Make snake faster
screenWidth = 800  // Make game area bigger
```

Then run: `go run main.go`

## 🐛 Troubleshooting

### Problem: "cannot find package github.com/hajimehoshi/ebiten"

**Solution 1: Run mod download first**
```bash
go mod download
go run main.go
```

**Solution 2: On Linux, install system libraries**
```bash
sudo apt-get install libgl1-mesa-dev libxrandr-dev libxcursor-dev libxi-dev libxinerama-dev libxext-dev
go mod download
go run main.go
```

### Problem: "permission denied" on macOS
```bash
# Allow the app
xattr -d com.apple.quarantine ./snake-game
```

### Problem: Window appears but immediately closes
This usually means the game panicked. Check:
```bash
go run main.go 2>&1 | head -20
```

Look for error messages.

### Problem: Game runs slowly
Try:
```bash
# Reduce tickSpeed in main.go
tickSpeed: 4  // Instead of 8
```

Lower numbers = faster snake

### Problem: Input is laggy
This is normal if you're pressing keys too fast. Try:
- Press one direction at a time
- Release before pressing next
- The snake buffers input automatically

### Problem: "cannot locate" error on Windows
Make sure you're in the correct directory:
```bash
cd "C:\path\to\snake-game-tutorial"
go run main.go
```

## 💻 Building an Executable

Want to share the game? Create an executable:

```bash
# Windows
go build -o snake-game.exe
# Run it
snake-game.exe

# macOS
go build -o snake-game
./snake-game

# Linux
go build -o snake-game
./snake-game
```

## 🎯 Common First Modifications

### Change Snake Speed
In **main.go**, find:
```go
tickSpeed: 8,  // Move every 8 frames
```

Try:
```go
tickSpeed: 4,  // Twice as fast
tickSpeed: 12, // Slower
```

### Change Game Grid Size
Find:
```go
const (
    gridWidth  = 20
    gridHeight = 20
```

Try:
```go
const (
    gridWidth  = 30  // Bigger grid
    gridHeight = 30
```

### Change Colors
Find in the `Draw` function:
```go
color.RGBA{R: 0, G: 255, B: 0, A: 255}  // Green snake
color.RGBA{R: 255, G: 0, B: 0, A: 255}  // Red food
```

Try different RGB values:
```go
// Blue snake
color.RGBA{R: 0, G: 0, B: 255, A: 255}

// Orange food
color.RGBA{R: 255, G: 165, B: 0, A: 255}
```

### Add Starting Score
Find:
```go
score: 0,
```

Change to:
```go
score: 100,  // Start with 100 points
```

## 📚 Next Learning Steps

### After Running Successfully:
1. Read **README.md** for complete overview
2. Check **QUICK_REFERENCE.md** for code snippets
3. Follow **TUTORIAL.md** to understand concepts
4. See **ADVANCED_FEATURES.md** to add features

### Ideas to Try:
- [ ] Change the grid size
- [ ] Modify colors to your preference
- [ ] Speed up or slow down the snake
- [ ] Read the step files to understand how it works
- [ ] Modify step1-7 files and run them individually
- [ ] Add your own feature (see ADVANCED_FEATURES.md)

## 🎓 Understanding the Code

### Where is the main loop?
File: **main.go**
- `Update()` - Called every frame for logic
- `Draw()` - Called every frame for graphics

### Where is the snake?
File: **main.go**
- `GameState.snake` - A slice of Points

### Where are the controls?
File: **main.go**
- `handleInput()` function

### Where is collision detection?
File: **main.go**
- `moveSnake()` function

## 📊 What Each Step teaches

| File | Time | Teaches |
|------|------|---------|
| step1_window.go | 10 min | Game interface basics |
| step2_render.go | 15 min | Drawing and coordinates |
| step3_movement.go | 20 min | Game timing and movement |
| step4_input.go | 15 min | Keyboard input |
| step5_collision.go | 25 min | Collision detection |
| step6_gameover.go | 20 min | Game states |
| step7_difficulty.go | 15 min | Dynamic difficulty |

## ✨ Quick Tips

1. **Want to run a specific step?**
   ```bash
   go run step3_movement.go  # Watch snake move
   go run step4_input.go     # Add controls
   ```

2. **Want to modify and test?**
   - Edit a file
   - Run: `go run filename.go`
   - See changes immediately!

3. **Want to see the complete game?**
   ```bash
   go run main.go
   ```

4. **Want to build a standalone game?**
   ```bash
   go build -o my-snake-game
   ./my-snake-game
   ```

## 🎮 Game Tips

- **Early Game**: Move smoothly in one direction
- **Mid Game**: Plan ahead for snake position
- **Late Game**: Space becomes valuable, move carefully!
- **High Score**: Focus on eating food efficiently
- **Survival**: Avoid walls at all costs

## 🚀 You're Ready!

Now just run:
```bash
go run main.go
```

And enjoy! 🎮

### Questions?
- Read **README.md** for overview
- Check **TUTORIAL.md** for concepts
- See **QUICK_REFERENCE.md** for code
- Check **ADVANCED_FEATURES.md** for ideas
- Review **INDEX.md** for navigation

---

**Happy Gaming!** 🐍

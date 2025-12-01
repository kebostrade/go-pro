# Snake Game Tutorial in Go with Ebiten

A comprehensive step-by-step tutorial for building a classic Snake game using Go and the [Ebiten](https://ebiten.org/) game engine.

## 📚 Project Structure

```
snake-game-tutorial/
├── main.go                  # Complete, production-ready solution
├── step1_window.go          # Step 1: Basic window and game loop
├── step2_render.go          # Step 2: Rendering the snake
├── step3_movement.go        # Step 3: Snake movement
├── step4_input.go           # Step 4: Keyboard input handling
├── step5_collision.go       # Step 5: Collision detection and food
├── step6_gameover.go        # Step 6: Game over and UI text
├── step7_difficulty.go      # Step 7: Difficulty progression
├── TUTORIAL.md              # Complete tutorial with explanations
├── ADVANCED_FEATURES.md     # Ideas for extending the game
└── go.mod                   # Module definition
```

## 🚀 Quick Start

### Prerequisites

- Go 1.23 or later
- Basic understanding of Go syntax
- Familiarity with game loops (optional)

### Installation

1. **Clone or download this project:**
   ```bash
   cd snake-game-tutorial
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the complete game:**
   ```bash
   go run main.go
   ```

### Linux Dependencies

If you're on Linux, you may need additional system dependencies:

```bash
# Ubuntu/Debian
sudo apt-get install libgl1-mesa-dev libxrandr-dev libxcursor-dev libxi-dev libxinerama-lib libxext-dev

# Fedora
sudo dnf install mesa-libGL-devel libxrandr-devel libxcursor-devel libxi-devel libxinerama-devel libxext-devel
```

## 🎮 How to Play

- **Space**: Start the game
- **Arrow Keys**: Move the snake
- **R**: Restart after game over

## 📖 Learning Path

Follow these steps to learn how the game is built:

### Step 1: Basic Window and Game Loop
```bash
go run step1_window.go
```
- Create an Ebiten game window
- Implement the Game interface
- Understand the update-draw loop

### Step 2: Rendering the Snake
```bash
go run step2_render.go
```
- Grid-based coordinate system
- Data structures (Point, snake slice)
- Drawing colored rectangles

### Step 3: Snake Movement
```bash
go run step3_movement.go
```
- Game tick system
- Direction vectors
- Slice operations (prepend head, remove tail)

### Step 4: Input Handling
```bash
go run step4_input.go
```
- Keyboard input detection
- Input buffering
- Preventing invalid turns

### Step 5: Collision Detection and Food
```bash
go run step5_collision.go
```
- Detecting wall collisions
- Self-collision checks
- Food mechanics and spawning
- Growing the snake

### Step 6: Game Over and UI
```bash
go run step6_gameover.go
```
- Game state management
- Text rendering
- Game start and restart logic

### Step 7: Difficulty Progression
```bash
go run step7_difficulty.go
```
- Increasing game speed
- Scaling difficulty based on score
- Game balance tuning

## 🎯 Complete Solution

Run the full, polished game:

```bash
go run main.go
```

The complete game includes all features from steps 1-7 plus optimizations.

## 📚 Key Concepts Covered

### Game Architecture
- **Game Loop**: Update (logic) → Draw (graphics) → Repeat
- **Game State**: Centralized data structure tracking all game data
- **State Management**: Handling different game states (start, playing, over)

### Input Handling
- Keyboard detection with `ebiten.IsKeyPressed()`
- One-time key events with `inpututil.IsKeyJustPressed()`
- Input buffering for smooth directional changes

### Collision Detection
- Point-to-point collision (snake head vs food)
- Boundary checking (snake vs walls)
- Self-collision detection (snake vs body)

### Game Mechanics
- Grid-based movement system
- Snake growth on food consumption
- Difficulty scaling based on score

## 🧩 Code Patterns

### Game State Structure
```go
type GameState struct {
    snake       []Point    // All body segments
    food        Point      // Food position
    direction   Point      // Current movement direction
    score       int        // Player score
    gameOver    bool       // Game state flag
    gameStarted bool       // Has game started?
    tickCount   int        // Frame counter
    tickSpeed   int        // Frames between moves
}
```

### Game Loop Pattern
```go
func (g *Game) Update() error {
    // Handle input
    g.handleInput()

    // Update game logic
    g.tickCount++
    if g.tickCount >= g.tickSpeed {
        g.moveSnake()
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Clear screen
    screen.Fill(color.Black)

    // Draw game objects
    // Draw UI text
}
```

### Collision Detection Pattern
```go
// Check wall collision
if newHead.X < 0 || newHead.X >= gridWidth {
    gameState.gameOver = true
}

// Check self collision
for _, segment := range snake {
    if newHead == segment {
        gameState.gameOver = true
    }
}

// Check food collision
if newHead == food {
    score += 10
    spawnFood()
}
```

## 🚀 Building an Executable

### Create a standalone executable:
```bash
go build -o snake-game
./snake-game
```

### Cross-compile for other platforms:
```bash
# macOS
GOOS=darwin GOARCH=amd64 go build -o snake-game-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o snake-game.exe

# Linux ARM64 (Raspberry Pi)
GOOS=linux GOARCH=arm64 go build -o snake-game-arm64
```

## 🎨 Customization Ideas

### Visual Improvements
- Add game backgrounds
- Use images instead of colored rectangles
- Add animations for eating/dying
- Improve color scheme

### Gameplay Features
- Multiple game modes (timed, endless, survival)
- Obstacles on the map
- Different food types with different effects
- Power-ups (speed boost, shield, etc.)

### Technical Enhancements
- Save/load game state
- Replay system
- Leaderboard with high scores
- Audio effects and background music
- Pause functionality

See `ADVANCED_FEATURES.md` for detailed implementation guides.

## 🐛 Troubleshooting

### Ebiten Installation Issues

**Problem**: "cannot find -lGL" on Linux
```bash
# Solution: Install graphics libraries
sudo apt-get install libgl1-mesa-dev libxrandr-dev
```

**Problem**: "no such file or directory" on Windows
```bash
# Solution: Install C compiler (MinGW or MSVC)
# Check Ebiten docs for Windows setup
```

### Game Runs Slowly

**Solution 1**: Reduce tickSpeed
```go
state.tickSpeed = 4  // Move every 4 frames instead of 8
```

**Solution 2**: Optimize drawing
```go
// Use images instead of Set()
// Pre-calculate positions instead of calculating every frame
```

### Input Lag

**Solution**: Use input buffering
```go
state.nextDir = newDirection  // Buffer input
// Apply on next move step
```

## 📖 Additional Resources

- **Ebiten Official Documentation**: https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2
- **Ebiten GitHub Examples**: https://github.com/hajimehoshi/ebiten/tree/main/examples
- **Go Game Development**: https://www.golang-book.com/
- **Game Development Patterns**: https://gameprogrammingpatterns.com/

## 🎓 Learning Goals

After completing this tutorial, you'll understand:

✅ How to structure a game in Go
✅ Event-driven programming with Ebiten
✅ Game loop fundamentals
✅ Collision detection algorithms
✅ State management in games
✅ User input handling
✅ Performance optimization basics
✅ Game design principles

## 💡 Next Steps

1. **Master the basics**: Run each step file and modify it
2. **Extend the game**: Add features from ADVANCED_FEATURES.md
3. **Optimize the code**: Refactor for performance and readability
4. **Create your own**: Build a new game using the patterns learned
5. **Share your work**: Publish on GitHub or itch.io

## 📝 License

This tutorial is provided as educational material. Feel free to use, modify, and distribute for learning purposes.

## 🤝 Contributing

Found an issue or have an improvement? Consider:
- Testing on different platforms
- Suggesting new features
- Improving explanations
- Adding more advanced examples

## ❓ FAQ

**Q: Is Ebiten only for 2D games?**
A: Primarily, but it can handle 3D with additional libraries.

**Q: Can I publish a game made with Ebiten?**
A: Yes! It's free and open source. Many games use Ebiten.

**Q: How do I add sound effects?**
A: Use the `github.com/hajimehoshi/ebiten/v2/audio` package.

**Q: What about multiplayer?**
A: You'd need to add networking. Consider libraries like `gorilla/websocket`.

**Q: Can I use images and sprites?**
A: Yes, use `ebiten.NewImageFromFile()` or create image assets.

---

Happy coding! 🎮

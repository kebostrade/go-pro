# Snake Game in Go with Ebiten - Complete Tutorial

Build a classic Snake game from scratch using Go and the Ebiten game engine. This tutorial progresses through 7 steps, each adding new functionality.

## Table of Contents

1. [Project Setup](#project-setup)
2. [Step 1: Basic Window and Game Loop](#step-1-basic-window-and-game-loop)
3. [Step 2: Render the Snake](#step-2-render-the-snake)
4. [Step 3: Snake Movement](#step-3-snake-movement)
5. [Step 4: Input Handling](#step-4-input-handling)
6. [Step 5: Food and Collision](#step-5-food-and-collision)
7. [Step 6: Score and Game Over](#step-6-score-and-game-over)
8. [Step 7: Difficulty and Polish](#step-7-difficulty-and-polish)
9. [Extending the Game](#extending-the-game)

---

## Project Setup

### Install Ebiten

```bash
# Install Ebiten
go get github.com/hajimehoshi/ebiten/v2

# For Linux, you may need additional dependencies
# Ubuntu/Debian:
sudo apt-get install libgl1-mesa-dev libxrandr-dev libxcursor-dev libxi-dev libxinerama-dev libxext-dev

# Install required tools
go install github.com/hajimehoshi/ebiten/v2/cmd/ebitencmd@latest
```

### Create Project Structure

```
snake-game-tutorial/
├── go.mod
├── main.go              # Completed solution
├── TUTORIAL.md          # This file
├── step1_window.go      # Step 1: Basic window
├── step2_render.go      # Step 2: Rendering
├── step3_movement.go    # Step 3: Movement
├── step4_input.go       # Step 4: Input
├── step5_collision.go   # Step 5: Collision
├── step6_gameover.go    # Step 6: Game Over
└── step7_difficulty.go  # Step 7: Difficulty
```

---

## Step 1: Basic Window and Game Loop

**Goal**: Create a window and game loop using Ebiten.

**Concepts**:
- Ebiten's `Game` interface
- `Update()` method - called every frame
- `Draw()` method - render graphics
- `Layout()` method - define screen dimensions

**File**: `step1_window.go`

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

type Game struct{}

// Update is called every frame for game logic
func (g *Game) Update() error {
	return nil
}

// Draw is called every frame to render graphics
func (g *Game) Draw(screen *ebiten.Image) {
	// Fill background with black
	screen.Fill(ebiten.ColorBlack)
}

// Layout defines screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}

	ebiten.SetWindowTitle("Snake Game - Step 1")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

**Key Points**:
- Implement the `ebiten.Game` interface with three methods
- `Update()` is called 60 times per second (default)
- Game runs in a loop automatically

---

## Step 2: Render the Snake

**Goal**: Draw the snake on the grid.

**Concepts**:
- Grid system (cells)
- Data structure for snake body
- Drawing filled rectangles

**File**: `step2_render.go`

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

const (
	gridWidth  = 20
	gridHeight = 20
	cellSize   = 20

	screenWidth  = gridWidth * cellSize
	screenHeight = gridHeight * cellSize
)

// Point represents a grid position
type Point struct {
	X, Y int
}

type GameState struct {
	snake []Point // Snake body
}

type Game struct {
	state *GameState
}

func NewGame() *Game {
	return &Game{
		state: &GameState{
			snake: []Point{
				{X: 10, Y: 10},
				{X: 9, Y: 10},
				{X: 8, Y: 10},
			},
		},
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ebiten.ColorBlack)

	// Draw snake
	for i, segment := range g.state.snake {
		x := segment.X * cellSize
		y := segment.Y * cellSize

		// Head in bright green, body in darker green
		var c color.Color
		if i == 0 {
			c = color.RGBA{0, 255, 0, 255} // Head
		} else {
			c = color.RGBA{0, 200, 0, 255} // Body
		}

		drawCell(screen, x, y, c)
	}
}

// drawCell draws a single grid cell
func drawCell(screen *ebiten.Image, x, y int, c color.Color) {
	for yy := y; yy < y+cellSize; yy++ {
		for xx := x; xx < x+cellSize; xx++ {
			screen.Set(xx, yy, c)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()

	ebiten.SetWindowTitle("Snake Game - Step 2")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
```

**Key Points**:
- Snake is stored as a slice of `Point` structs
- Head is at index 0
- Grid system allows easy coordinate mapping

---

## Step 3: Snake Movement

**Goal**: Make the snake move continuously in one direction.

**Concepts**:
- Frame timing/ticks
- Direction vector
- Updating snake position

**File**: `step3_movement.go`

```go
type GameState struct {
	snake      []Point
	direction  Point // Direction to move
	tickCount  int   // Frames until next move
	tickSpeed  int   // How many frames between moves
}

func NewGame() *Game {
	return &Game{
		state: &GameState{
			snake: []Point{
				{X: 10, Y: 10},
				{X: 9, Y: 10},
				{X: 8, Y: 10},
			},
			direction: Point{X: 1, Y: 0}, // Moving right
			tickSpeed: 8, // Move every 8 frames
		},
	}
}

func (g *Game) Update() error {
	state := g.state

	// Count frames
	state.tickCount++

	// Move every tickSpeed frames
	if state.tickCount >= state.tickSpeed {
		state.tickCount = 0
		g.moveSnake()
	}

	return nil
}

func (g *Game) moveSnake() {
	state := g.state

	// Calculate new head position
	head := state.snake[0]
	newHead := Point{
		X: head.X + state.direction.X,
		Y: head.Y + state.direction.Y,
	}

	// Add new head to front
	state.snake = append([]Point{newHead}, state.snake...)

	// Remove tail (snake moves forward)
	state.snake = state.snake[:len(state.snake)-1]
}
```

**Key Points**:
- Use `tickSpeed` to slow down movement (not every frame)
- Prepend new head to snake slice
- Remove tail so snake length stays constant
- Direction is a Point with X and Y components

---

## Step 4: Input Handling

**Goal**: Allow player to control snake direction with arrow keys.

**Concepts**:
- Input detection
- Direction buffering
- Preventing 180-degree turns

**File**: `step4_input.go`

```go
package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Update() error {
	state := g.state

	// Handle input
	g.handleInput()

	// Movement tick
	state.tickCount++
	if state.tickCount >= state.tickSpeed {
		state.tickCount = 0
		g.moveSnake()
	}

	return nil
}

func (g *Game) handleInput() {
	state := g.state

	// Only allow direction change if not opposite
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && state.direction.Y == 0 {
		state.direction = Point{X: 0, Y: -1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && state.direction.Y == 0 {
		state.direction = Point{X: 0, Y: 1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && state.direction.X == 0 {
		state.direction = Point{X: -1, Y: 0}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && state.direction.X == 0 {
		state.direction = Point{X: 1, Y: 0}
	}
}
```

**Key Points**:
- Check that new direction isn't opposite of current
- This prevents turning 180° into yourself
- Use `ebiten.IsKeyPressed()` for continuous input
- Use `inpututil.IsKeyJustPressed()` for single key presses

---

## Step 5: Food and Collision

**Goal**: Add food spawning and collision detection.

**Concepts**:
- Collision detection
- Growing snake on food
- Boundary checking

**File**: `step5_collision.go`

```go
type GameState struct {
	snake     []Point
	food      Point
	direction Point
	score     int
	// ... other fields
}

func (g *Game) moveSnake() {
	state := g.state

	head := state.snake[0]
	newHead := Point{
		X: head.X + state.direction.X,
		Y: head.Y + state.direction.Y,
	}

	// Check wall collision
	if newHead.X < 0 || newHead.X >= gridWidth ||
	   newHead.Y < 0 || newHead.Y >= gridHeight {
		state.gameOver = true
		return
	}

	// Check self collision
	for _, segment := range state.snake {
		if newHead == segment {
			state.gameOver = true
			return
		}
	}

	// Add new head
	state.snake = append([]Point{newHead}, state.snake...)

	// Check food collision
	if newHead == state.food {
		state.score += 10
		g.spawnFood()
		// Don't remove tail - snake grows!
	} else {
		// Remove tail if no food eaten
		state.snake = state.snake[:len(state.snake)-1]
	}
}

func (g *Game) spawnFood() {
	// Spawn food at a random unoccupied position
	state := g.state
	for {
		x := (state.score/10 + 7) % gridWidth
		y := (state.score/10 + 3) % gridHeight
		candidate := Point{X: x, Y: y}

		// Check if position is occupied
		occupied := false
		for _, segment := range state.snake {
			if segment == candidate {
				occupied = true
				break
			}
		}

		if !occupied {
			state.food = candidate
			return
		}
	}
}
```

**Key Points**:
- Compare points for collision using `==` operator
- Wall collision: check coordinate bounds
- Self collision: check all snake segments
- Food collision: grow snake by not removing tail

---

## Step 6: Score and Game Over

**Goal**: Track score and handle game over state.

**Concepts**:
- Game states (running, game over)
- Displaying text
- Game restart

**File**: `step6_gameover.go`

```go
type GameState struct {
	snake       []Point
	food        Point
	direction   Point
	score       int
	gameOver    bool
	gameStarted bool
	// ... other fields
}

func (g *Game) Update() error {
	state := g.state

	// Start game with Space
	if !state.gameStarted && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.gameStarted = true
		return nil
	}

	if state.gameStarted && !state.gameOver {
		g.handleInput()

		state.tickCount++
		if state.tickCount >= state.tickSpeed {
			state.tickCount = 0
			g.moveSnake()
		}
	}

	// Restart with R
	if state.gameOver && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.state = NewGame().state
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(ebiten.ColorBlack)

	// Draw game objects
	for i, segment := range g.state.snake {
		x := segment.X * cellSize
		y := segment.Y * cellSize

		var c color.Color
		if i == 0 {
			c = color.RGBA{0, 255, 0, 255}
		} else {
			c = color.RGBA{0, 200, 0, 255}
		}

		drawCell(screen, x, y, c)
	}

	// Draw food
	foodX := g.state.food.X * cellSize
	foodY := g.state.food.Y * cellSize
	drawCell(screen, foodX, foodY, color.RGBA{255, 0, 0, 255})

	// Draw UI
	if !g.state.gameStarted {
		ebitenutil.DebugPrintAt(screen, "Press SPACE to start", 10, 10)
	} else if g.state.gameOver {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("GAME OVER - Score: %d", g.state.score), 10, 10)
		ebitenutil.DebugPrintAt(screen, "Press R to restart", 10, 30)
	} else {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("Score: %d", g.state.score), 10, 10)
	}
}
```

**Key Points**:
- Use boolean flags for game states
- Use `inpututil.IsKeyJustPressed()` for one-time events
- Display text with `ebitenutil.DebugPrintAt()`

---

## Step 7: Difficulty and Polish

**Goal**: Increase difficulty and add game polish.

**Concepts**:
- Progressive difficulty
- Speed scaling
- Better rendering

**File**: `step7_difficulty.go`

```go
func (g *Game) moveSnake() {
	state := g.state

	// ... collision checks ...

	if newHead == state.food {
		state.score += 10
		g.spawnFood()

		// Increase difficulty: speed up every 50 points
		if state.score%50 == 0 && state.tickSpeed > 4 {
			state.tickSpeed--
		}
	} else {
		state.snake = state.snake[:len(state.snake)-1]
	}
}
```

**Additional Polish**:
- Better colors and rendering
- Sound effects (use `audio` package)
- Different game modes
- High score tracking
- Pause functionality

---

## Extending the Game

### Feature Ideas

1. **Better Rendering**
   ```go
   // Use ebitenutil for images instead of Set()
   const cellImage = `
   iVBORw0KGgoAAAANSUhEUgAAAA...`
   // Or load PNG files
   ```

2. **Sound Effects**
   ```go
   import "github.com/hajimehoshi/ebiten/v2/audio"

   // Load and play sounds on eat, crash, etc.
   ```

3. **Pause Feature**
   ```go
   if ebiten.IsKeyPressed(ebiten.KeyP) {
       state.paused = !state.paused
   }
   ```

4. **Multiple Difficulty Levels**
   ```go
   type Difficulty int
   const (
       Easy Difficulty = iota
       Medium
       Hard
   )
   ```

5. **High Score Persistence**
   ```go
   import "encoding/json"
   // Save scores to file
   ```

6. **Smooth Rendering with Images**
   ```go
   // Use Vector drawing or image libraries
   // instead of Set() for better performance
   ```

---

## Running the Game

### From the Tutorial Files

```bash
# Run completed solution
go run main.go

# Run step by step
go run step1_window.go
go run step2_render.go
# ... etc
```

### Build Executable

```bash
go build -o snake-game
./snake-game
```

---

## Troubleshooting

**Ebiten won't install**:
- Ensure you have required C dependencies
- On Linux: `libgl1-mesa-dev libxrandr-dev libxcursor-dev libxi-dev libxinerama-dev libxext-dev`

**Game runs slowly**:
- Reduce `tickSpeed` (lower = faster)
- Optimize drawing (use images instead of Set())

**Input feels laggy**:
- Increase frame rate with `ebiten.SetMaxTPS()`
- Use input buffering (nextDir pattern)

**Food spawns inside snake**:
- Add better spawn algorithm
- Use randomization (math/rand)

---

## Learning Resources

- [Ebiten Official Docs](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2)
- [Ebiten GitHub Examples](https://github.com/hajimehoshi/ebiten/tree/main/examples)
- [Go Game Programming](https://github.com/ebiten/tutorial)

---

## Next Steps

- Add multiple difficulty levels
- Implement high score system
- Create different game modes
- Add AI opponent
- Publish as standalone executable

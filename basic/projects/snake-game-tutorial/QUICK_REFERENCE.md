# Snake Game - Quick Reference Guide

A handy cheat sheet for game development concepts and Ebiten API.

## 🎮 Game Loop

```go
// Ebiten handles this automatically
for {
    Update()   // Game logic every frame
    Draw()     // Graphics every frame
}
```

## 🎯 Core Ebiten API

### Window Setup
```go
ebiten.SetWindowTitle("Title")
ebiten.SetWindowSize(800, 600)
ebiten.SetWindowResizable(true)
ebiten.SetMaxTPS(60)  // Frames per second
```

### Screen Operations
```go
// Clear screen
screen.Fill(ebiten.ColorBlack)

// Set single pixel
screen.Set(x, y, ebiten.Color{R: 255, G: 0, B: 0, A: 255})

// Draw image
opts := &ebiten.DrawImageOptions{}
opts.GeoM.Translate(x, y)
screen.DrawImage(image, opts)
```

### Input
```go
// Continuous input
if ebiten.IsKeyPressed(ebiten.KeyArrowUp) { }

// Single key press
if inpututil.IsKeyJustPressed(ebiten.KeySpace) { }

// Mouse input
x, y := ebiten.CursorPosition()
if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) { }
```

### Text
```go
import "github.com/hajimehoshi/ebiten/v2/ebitenutil"

ebitenutil.DebugPrint(screen, "text")
ebitenutil.DebugPrintAt(screen, "text", x, y)
```

## 📊 Game State Pattern

```go
type GameState struct {
    // Game data
    player   Player
    enemies  []Enemy
    items    []Item

    // State flags
    paused   bool
    gameOver bool

    // Counters
    score    int
    frameCount int
}

type Game struct {
    state *GameState
}

func (g *Game) Update() error {
    // Update logic
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Render graphics
}

func (g *Game) Layout(w, h int) (int, int) {
    return screenWidth, screenHeight
}
```

## 🎨 Colors

```go
// Named colors
ebiten.ColorBlack
ebiten.ColorWhite

// RGBA colors
color.RGBA{R: 255, G: 0, B: 0, A: 255}  // Red

// Color constants
const Red = color.RGBA{R: 255, G: 0, B: 0, A: 255}
```

## 🔢 Common Key Codes

```go
// Arrow keys
ebiten.KeyArrowUp
ebiten.KeyArrowDown
ebiten.KeyArrowLeft
ebiten.KeyArrowRight

// Common keys
ebiten.KeySpace
ebiten.KeyEnter
ebiten.KeyEscape

// WASD
ebiten.KeyW
ebiten.KeyA
ebiten.KeyS
ebiten.KeyD

// Number/Letter keys
ebiten.Key0 through ebiten.Key9
ebiten.KeyA through ebiten.KeyZ
```

## 📐 Grid Coordinate System

```go
// Convert grid to pixels
pixelX := gridX * cellSize
pixelY := gridY * cellSize

// Convert pixels to grid
gridX := pixelX / cellSize
gridY := pixelY / cellSize

// Direction vectors
Up    = Point{X: 0, Y: -1}
Down  = Point{X: 0, Y: 1}
Left  = Point{X: -1, Y: 0}
Right = Point{X: 1, Y: 0}
```

## 🔄 Slice Operations

```go
// Append
slice = append(slice, item)

// Prepend
slice = append([]Item{item}, slice...)

// Remove at index
slice = append(slice[:i], slice[i+1:]...)

// Remove last
slice = slice[:len(slice)-1]

// Remove first
slice = slice[1:]

// Iterate
for i, item := range slice {
    // i is index, item is value
}
```

## ⚙️ Frame Timing

```go
const tickSpeed = 8  // Move every 8 frames

tickCount++
if tickCount >= tickSpeed {
    tickCount = 0
    // Move snake
}

// At 60 FPS:
// tickSpeed 8 = 7.5 updates/sec
// tickSpeed 6 = 10 updates/sec
// tickSpeed 4 = 15 updates/sec
```

## 💥 Collision Detection

```go
// Point-to-point
if head == food {
    // Collision!
}

// Point in list
for _, point := range points {
    if head == point {
        // Hit!
    }
}

// Boundary check
if x < 0 || x >= gridWidth {
    // Out of bounds!
}

// Rectangle (AABB)
if x > rect.X && x < rect.X + rect.Width &&
   y > rect.Y && y < rect.Y + rect.Height {
    // Inside rectangle!
}
```

## 📦 Common Data Structures

```go
// Point in 2D space
type Point struct {
    X, Y int
}

// Rectangle
type Rect struct {
    X, Y, Width, Height int
}

// Generic entity
type Entity struct {
    Position Point
    Velocity Point
    Size     int
    Color    color.Color
}

// Game state
type GameState struct {
    Score int
    Lives int
    Level int
    Paused bool
}
```

## 🎯 Design Patterns

### Tick-Based Movement
```go
tickCount++
if tickCount >= tickSpeed {
    tickCount = 0
    move()
}
```

### Input Buffering
```go
// Store desired direction
nextDir = newInput()

// Apply on next tick
if isValidMove(nextDir) {
    currentDir = nextDir
}
```

### State Machine
```go
const (
    StateMenu = iota
    StatePlaying
    StatePaused
    StateGameOver
)

var currentState = StateMenu

switch currentState {
case StateMenu:
    // Handle menu
case StatePlaying:
    // Handle game
case StatePaused:
    // Handle pause
}
```

### Object Pool
```go
type Pool struct {
    objects []Object
}

func (p *Pool) Get() Object {
    if len(p.objects) > 0 {
        obj := p.objects[len(p.objects)-1]
        p.objects = p.objects[:len(p.objects)-1]
        return obj
    }
    return NewObject()
}
```

## 🧮 Math Utilities

```go
// Absolute value
import "math"
abs := math.Abs(float64(x))

// Distance
dist := math.Sqrt(float64(dx*dx + dy*dy))

// Clamp value
if x < min { x = min }
if x > max { x = max }

// Random number
import "math/rand"
rand.Intn(100)      // 0-99
rand.Float64() * 100 // 0.0-100.0
```

## 📝 Common Patterns in Snake

### Check Direction Validity
```go
// Prevent 180-degree turns
isOpposite := currentDir.X == -newDir.X &&
             currentDir.Y == -newDir.Y

if !isOpposite {
    currentDir = newDir
}
```

### Spawn at Empty Position
```go
for {
    x := rand.Intn(gridWidth)
    y := rand.Intn(gridHeight)
    pos := Point{X: x, Y: y}

    if !isOccupied(pos, snake) {
        return pos
    }
}
```

### Grow Snake
```go
// Add head
snake = append([]Point{newHead}, snake...)

if ate_food {
    // Don't remove tail - snake grows
} else {
    // Remove tail - snake moves
    snake = snake[:len(snake)-1]
}
```

### Move to New Position
```go
newHead := Point{
    X: head.X + direction.X,
    Y: head.Y + direction.Y,
}

// Wrapping (optional)
newHead.X = (newHead.X + gridWidth) % gridWidth
newHead.Y = (newHead.Y + gridHeight) % gridHeight
```

## 🐛 Debug Tricks

```go
// Print to console
log.Println("Debug:", value)

// Show FPS
fmt.Sprintf("FPS: %v", ebiten.ActualFPS())

// Draw debug info
ebitenutil.DebugPrintAt(screen,
    fmt.Sprintf("X: %d Y: %d", x, y), 10, 10)

// Highlight position
drawRect(screen, x, y, size, size, color.RGBA{255, 0, 0, 255})
```

## ⚡ Performance Tips

```go
// Don't create new objects every frame
image := ebiten.NewImage(...)  // BAD - create once
screen.DrawImage(image, opts)   // GOOD

// Cache calculations
cellPixels := cellSize * cellSize  // Calculate once

// Use slices efficiently
items := make([]Item, 0, 100)  // Pre-allocate capacity

// Avoid unnecessary allocations
slice := append([]Item{item}, slice...)  // Creates new slice
slice = append([]Item{item}, slice...)   // Better with pre-allocated

// Use pointers for large structures
player := &Player{}  // Cheaper to pass around
```

## 📚 File I/O

```go
import "os"

// Read file
data, err := os.ReadFile("file.txt")
if err != nil { log.Fatal(err) }

// Write file
err := os.WriteFile("file.txt", []byte(data), 0644)
if err != nil { log.Fatal(err) }

// Check file exists
_, err := os.Stat("file.txt")
if os.IsNotExist(err) { /* doesn't exist */ }
```

## 🎬 Animation Timing

```go
const animationSpeed = 4  // Change every 4 frames
animFrame := (frameCount / animationSpeed) % numFrames

// Pulsing effect
intensity := math.Sin(float64(frameCount) * 0.1) * 127 + 128
```

## 🎵 Loading Resources

```go
// Image
img, _ := ebiten.NewImageFromFile("sprite.png")

// Sound
f, _ := os.Open("sound.wav")
defer f.Close()
context := audio.NewContext(44100)
player, _ := context.NewPlayer(f)

// JSON
var data struct {
    Name string
    Score int
}
json.Unmarshal(jsonBytes, &data)
```

## 🏗️ Project Structure

```
game/
├── main.go                 # Entry point
├── game.go                 # Game struct and main logic
├── entities.go             # Player, Enemy, Item structs
├── collision.go            # Collision detection
├── rendering.go            # Draw functions
├── input.go                # Input handling
├── assets/                 # Images, sounds, etc
└── go.mod                  # Module definition
```

## ✅ Testing Checklist

- [ ] Collision detection works correctly
- [ ] Input is responsive
- [ ] Game loop runs smoothly (60 FPS)
- [ ] Score updates correctly
- [ ] Game over triggers properly
- [ ] Can restart game
- [ ] No memory leaks (check FPS over time)
- [ ] Performance acceptable on target hardware

---

**Remember**: Start simple, test often, optimize later!

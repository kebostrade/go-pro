// +build ignore
// STEP 4: Input Handling
//
// Allow the player to control the snake with arrow keys.
// This introduces:
// - Keyboard input detection
// - Input validation (prevent 180-degree turns)
// - Input buffering (nextDir pattern)
//
// Run with: go run step4_input.go
//
// Controls:
// - Arrow Keys: Change direction
// - The snake cannot turn 180 degrees into itself

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

type Point struct {
	X, Y int
}

type GameStateStep4 struct {
	snake    []Point
	direction Point // Current direction
	nextDir   Point // Buffered next direction (for smooth input)
	tickCount int
	tickSpeed int
}

type GameStep4 struct {
	state *GameStateStep4
}

func NewGameStep4() *GameStep4 {
	return &GameStep4{
		state: &GameStateStep4{
			snake: []Point{
				{X: 10, Y: 10},
				{X: 9, Y: 10},
				{X: 8, Y: 10},
			},
			direction: Point{X: 1, Y: 0},
			nextDir:   Point{X: 1, Y: 0},
			tickSpeed: 8,
		},
	}
}

func (g *GameStep4) Update() error {
	state := g.state

	// Handle keyboard input
	g.handleInput()

	// Update position
	state.tickCount++
	if state.tickCount >= state.tickSpeed {
		state.tickCount = 0
		g.moveSnake()
	}

	return nil
}

// handleInput processes keyboard input
func (g *GameStep4) handleInput() {
	state := g.state

	// Check each arrow key and update nextDir if it's valid
	// Valid means: the new direction isn't opposite of current direction
	// This prevents the snake from turning 180 degrees

	// Up arrow (0, -1)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && state.direction.Y == 0 {
		// Only allow if we're not moving vertically
		state.nextDir = Point{X: 0, Y: -1}
	}

	// Down arrow (0, 1)
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && state.direction.Y == 0 {
		// Only allow if we're not moving vertically
		state.nextDir = Point{X: 0, Y: 1}
	}

	// Left arrow (-1, 0)
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && state.direction.X == 0 {
		// Only allow if we're not moving horizontally
		state.nextDir = Point{X: -1, Y: 0}
	}

	// Right arrow (1, 0)
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && state.direction.X == 0 {
		// Only allow if we're not moving horizontally
		state.nextDir = Point{X: 1, Y: 0}
	}
}

func (g *GameStep4) moveSnake() {
	state := g.state

	// Apply the buffered direction if it's valid
	// (not exactly opposite to current direction)
	if state.nextDir != state.direction {
		// Check if nextDir is opposite of direction
		isOpposite := (state.direction.X == -state.nextDir.X &&
			state.direction.Y == -state.nextDir.Y)

		if !isOpposite {
			state.direction = state.nextDir
		}
	}

	// Calculate new head position
	head := state.snake[0]
	newHead := Point{
		X: head.X + state.direction.X,
		Y: head.Y + state.direction.Y,
	}

	// Add new head and remove tail
	state.snake = append([]Point{newHead}, state.snake...)
	state.snake = state.snake[:len(state.snake)-1]
}

func (g *GameStep4) Draw(screen *ebiten.Image) {
	screen.Fill(ebiten.ColorBlack)

	// Draw snake
	for i, segment := range g.state.snake {
		pixelX := segment.X * cellSize
		pixelY := segment.Y * cellSize

		var cellColor color.Color
		if i == 0 {
			cellColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
		} else {
			cellColor = color.RGBA{R: 0, G: 200, B: 0, A: 255}
		}

		drawCell(screen, pixelX, pixelY, cellColor)
	}
}

func (g *GameStep4) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func drawCell(screen *ebiten.Image, x, y int, c color.Color) {
	for yy := y; yy < y+cellSize; yy++ {
		for xx := x; xx < x+cellSize; xx++ {
			screen.Set(xx, yy, c)
		}
	}
}

func main() {
	game := NewGameStep4()

	ebiten.SetWindowTitle("Snake Game - Step 4: Input")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

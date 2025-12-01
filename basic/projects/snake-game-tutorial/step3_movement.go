// +build ignore
// STEP 3: Snake Movement
//
// Now the snake will move continuously across the grid.
// This introduces:
// - Game tick/frame counting
// - Direction vectors
// - Moving the snake (prepend head, remove tail)
//
// Run with: go run step3_movement.go

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

type GameStateStep3 struct {
	snake     []Point // Snake segments
	direction Point   // Current movement direction (X, Y velocity)
	tickCount int     // Frame counter
	tickSpeed int     // How many frames between each move
}

type GameStep3 struct {
	state *GameStateStep3
}

func NewGameStep3() *GameStep3 {
	return &GameStep3{
		state: &GameStateStep3{
			snake: []Point{
				{X: 10, Y: 10},
				{X: 9, Y: 10},
				{X: 8, Y: 10},
			},
			// Start moving to the right (1 cell right per move, 0 vertical)
			direction: Point{X: 1, Y: 0},
			tickSpeed: 8, // Move every 8 frames (~8 moves per second at 60 FPS)
		},
	}
}

func (g *GameStep3) Update() error {
	state := g.state

	// Increment frame counter
	state.tickCount++

	// Move snake when enough frames have passed
	if state.tickCount >= state.tickSpeed {
		state.tickCount = 0
		g.moveSnake()
	}

	return nil
}

// moveSnake performs one step of snake movement
func (g *GameStep3) moveSnake() {
	state := g.state

	// Get current head
	head := state.snake[0]

	// Calculate new head position based on direction
	newHead := Point{
		X: head.X + state.direction.X,
		Y: head.Y + state.direction.Y,
	}

	// Add new head to the front of the snake
	// This is a slice operation: prepend newHead to state.snake
	state.snake = append([]Point{newHead}, state.snake...)

	// Remove the tail so the snake maintains its length
	// The snake appears to move forward
	state.snake = state.snake[:len(state.snake)-1]

	// Optional: Add simple boundary wrapping (snake wraps around edges)
	// if state.snake[0].X < 0 {
	//     state.snake[0].X = gridWidth - 1
	// }
	// if state.snake[0].X >= gridWidth {
	//     state.snake[0].X = 0
	// }
	// etc.
}

func (g *GameStep3) Draw(screen *ebiten.Image) {
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

func (g *GameStep3) Layout(outsideWidth, outsideHeight int) (int, int) {
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
	game := NewGameStep3()

	ebiten.SetWindowTitle("Snake Game - Step 3: Movement")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

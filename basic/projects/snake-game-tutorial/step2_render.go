// +build ignore
// STEP 2: Rendering the Snake
//
// We'll now render a snake on the screen using a grid system.
// This introduces:
// - Grid-based coordinate system
// - Data structure for the snake
// - Rendering with pixel manipulation
//
// Run with: go run step2_render.go
//
// The snake is represented as a slice of Points where:
// - Index 0 is the head
// - Index 1+ are the body segments

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

const (
	// Grid dimensions
	gridWidth  = 20
	gridHeight = 20
	cellSize   = 20 // pixels per cell

	// Screen dimensions (automatically calculated from grid)
	screenWidth  = gridWidth * cellSize
	screenHeight = gridHeight * cellSize
)

// Point represents a position in the grid
type Point struct {
	X, Y int
}

type GameStateStep2 struct {
	snake []Point // List of snake segments
}

type GameStep2 struct {
	state *GameStateStep2
}

func NewGameStep2() *GameStep2 {
	return &GameStep2{
		state: &GameStateStep2{
			// Create initial snake in the middle of the grid
			// Head is at index 0, body extends to the left
			snake: []Point{
				{X: 10, Y: 10}, // Head
				{X: 9, Y: 10},  // Body segment
				{X: 8, Y: 10},  // Body segment
			},
		},
	}
}

func (g *GameStep2) Update() error {
	// No game logic yet - just rendering
	return nil
}

func (g *GameStep2) Draw(screen *ebiten.Image) {
	// Fill background with black
	screen.Fill(ebiten.ColorBlack)

	// Draw each snake segment
	for i, segment := range g.state.snake {
		// Convert grid coordinates to pixel coordinates
		pixelX := segment.X * cellSize
		pixelY := segment.Y * cellSize

		// Different color for head vs body
		var cellColor color.Color
		if i == 0 {
			// Head: bright green
			cellColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
		} else {
			// Body: darker green
			cellColor = color.RGBA{R: 0, G: 200, B: 0, A: 255}
		}

		// Draw the cell
		drawCell(screen, pixelX, pixelY, cellColor)
	}

	// TODO: Draw food here in later steps
}

func (g *GameStep2) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// drawCell fills a single grid cell with a color
// This is a simple approach - later we can optimize with images
func drawCell(screen *ebiten.Image, x, y int, c color.Color) {
	// Draw all pixels within this cell
	for yy := y; yy < y+cellSize; yy++ {
		for xx := x; xx < x+cellSize; xx++ {
			screen.Set(xx, yy, c)
		}
	}
}

func main() {
	game := NewGameStep2()

	ebiten.SetWindowTitle("Snake Game - Step 2: Rendering")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

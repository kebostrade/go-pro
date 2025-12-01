// Package main implements a classic Snake game using Ebiten.
// This is a complete, runnable implementation demonstrating:
// - Game loop and state management
// - Input handling (keyboard controls)
// - Collision detection
// - Score tracking and game over logic
package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	// Grid settings
	gridWidth  = 20
	gridHeight = 20
	cellSize   = 20

	// Screen dimensions
	screenWidth  = gridWidth * cellSize
	screenHeight = gridHeight * cellSize
)

// Point represents a position on the game grid
type Point struct {
	X, Y int
}

// GameState represents all game data
type GameState struct {
	snake       []Point // Snake body, head is at index 0
	food        Point
	direction   Point // Current movement direction
	nextDir     Point // Next direction (allows buffering input)
	score       int
	gameOver    bool
	gameStarted bool
	tickCount   int
	tickSpeed   int // How many frames until next move
}

// Game implements ebiten.Game interface
type Game struct {
	state *GameState
}

// NewGame creates a new game instance
func NewGame() *Game {
	return &Game{
		state: &GameState{
			snake: []Point{
				{X: gridWidth / 2, Y: gridHeight / 2},
				{X: gridWidth/2 - 1, Y: gridHeight / 2},
				{X: gridWidth/2 - 2, Y: gridHeight / 2},
			},
			direction:   Point{X: 1, Y: 0}, // Start moving right
			nextDir:     Point{X: 1, Y: 0},
			food:        Point{X: 10, Y: 5},
			score:       0,
			gameOver:    false,
			gameStarted: false,
			tickSpeed:   8, // Move every 8 frames (~8 times per second at 60 FPS)
		},
	}
}

// Update is called every frame to update game state
func (g *Game) Update() error {
	state := g.state

	if !state.gameStarted && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		state.gameStarted = true
		return nil
	}

	if state.gameStarted && !state.gameOver {
		// Handle input
		g.handleInput()

		// Update game state
		state.tickCount++
		if state.tickCount >= state.tickSpeed {
			state.tickCount = 0
			g.update()
		}
	}

	// Reset game on game over
	if state.gameOver && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.state = NewGame().state
	}

	return nil
}

// handleInput processes keyboard input
func (g *Game) handleInput() {
	state := g.state

	// Allow direction changes (with simple validation to prevent 180-degree turns)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && state.direction.Y == 0 {
		state.nextDir = Point{X: 0, Y: -1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && state.direction.Y == 0 {
		state.nextDir = Point{X: 0, Y: 1}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) && state.direction.X == 0 {
		state.nextDir = Point{X: -1, Y: 0}
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) && state.direction.X == 0 {
		state.nextDir = Point{X: 1, Y: 0}
	}
}

// update performs one game step
func (g *Game) update() {
	state := g.state

	// Apply next direction if it's valid (not opposite of current direction)
	if state.nextDir != state.direction {
		if !(state.direction.X == -state.nextDir.X && state.direction.Y == -state.nextDir.Y) {
			state.direction = state.nextDir
		}
	}

	// Calculate new head position
	head := state.snake[0]
	newHead := Point{
		X: head.X + state.direction.X,
		Y: head.Y + state.direction.Y,
	}

	// Check wall collision
	if newHead.X < 0 || newHead.X >= gridWidth || newHead.Y < 0 || newHead.Y >= gridHeight {
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
		// Increase difficulty slightly every 50 points
		if state.score%50 == 0 && state.tickSpeed > 4 {
			state.tickSpeed--
		}
	} else {
		// Remove tail if no food eaten
		state.snake = state.snake[:len(state.snake)-1]
	}
}

// spawnFood generates a new food position
func (g *Game) spawnFood() {
	// Simple random-like spawn (deterministic for testing)
	// In a real game, use math/rand for true randomness
	state := g.state
	for {
		x := (state.score/10 + 7) % gridWidth
		y := (state.score/10 + 3) % gridHeight
		newFood := Point{X: x, Y: y}

		// Make sure food doesn't spawn on snake
		canSpawn := true
		for _, segment := range state.snake {
			if segment == newFood {
				canSpawn = false
				break
			}
		}

		if canSpawn {
			state.food = newFood
			return
		}
	}
}

// Draw is called to render the game
func (g *Game) Draw(screen *ebiten.Image) {
	state := g.state

	// Fill background
	screen.Fill(color.Black)

	// Draw snake
	for i, segment := range state.snake {
		x := segment.X * cellSize
		y := segment.Y * cellSize

		// Head is brighter green
		var cellColor color.Color
		if i == 0 {
			cellColor = color.RGBA{R: 0, G: 255, B: 0, A: 255} // Bright green for head
		} else {
			cellColor = color.RGBA{R: 0, G: 200, B: 0, A: 255} // Darker green for body
		}

		drawRect(screen, x, y, cellSize, cellSize, cellColor)
	}

	// Draw food (red)
	foodX := state.food.X * cellSize
	foodY := state.food.Y * cellSize
	drawRect(screen, foodX, foodY, cellSize, cellSize, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Draw UI text
	if !state.gameStarted {
		ebitenutil.DebugPrintAt(screen, "SNAKE GAME", 10, 10)
		ebitenutil.DebugPrintAt(screen, "Press SPACE to start", 10, 30)
	} else if state.gameOver {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("GAME OVER - Score: %d", state.score), 10, 10)
		ebitenutil.DebugPrintAt(screen, "Press R to restart", 10, 30)
	} else {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", state.score), 10, 10)
	}
}

// Layout defines the screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// drawRect is a helper to draw filled rectangles
func drawRect(screen *ebiten.Image, x, y, width, height int, cellColor color.Color) {
	for yy := y; yy < y+height; yy++ {
		for xx := x; xx < x+width; xx++ {
			screen.Set(xx, yy, cellColor)
		}
	}
}

func main() {
	game := NewGame()

	ebiten.SetWindowTitle("Snake Game")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

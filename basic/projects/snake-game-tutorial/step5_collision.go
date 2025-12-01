// +build ignore
// STEP 5: Collision Detection and Food
//
// Add collision detection for walls and self, plus food mechanics.
// This introduces:
// - Food spawning
// - Collision detection (walls and body)
// - Growing the snake (not removing tail when eating)
// - Game over state
//
// Run with: go run step5_collision.go

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

type GameStateStep5 struct {
	snake     []Point
	food      Point
	direction Point
	nextDir   Point
	tickCount int
	tickSpeed int
	gameOver  bool // New: track if game ended
	score     int  // New: track score
}

type GameStep5 struct {
	state *GameStateStep5
}

func NewGameStep5() *GameStep5 {
	game := &GameStep5{
		state: &GameStateStep5{
			snake: []Point{
				{X: 10, Y: 10},
				{X: 9, Y: 10},
				{X: 8, Y: 10},
			},
			direction: Point{X: 1, Y: 0},
			nextDir:   Point{X: 1, Y: 0},
			tickSpeed: 8,
			gameOver:  false,
			score:     0,
			food:      Point{X: 15, Y: 10}, // Initial food position
		},
	}
	return game
}

func (g *GameStep5) Update() error {
	state := g.state

	// Don't update if game is over
	if state.gameOver {
		return nil
	}

	g.handleInput()

	state.tickCount++
	if state.tickCount >= state.tickSpeed {
		state.tickCount = 0
		g.moveSnake()
	}

	return nil
}

func (g *GameStep5) handleInput() {
	state := g.state

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

func (g *GameStep5) moveSnake() {
	state := g.state

	// Apply buffered direction
	if state.nextDir != state.direction {
		isOpposite := (state.direction.X == -state.nextDir.X &&
			state.direction.Y == -state.nextDir.Y)
		if !isOpposite {
			state.direction = state.nextDir
		}
	}

	// Calculate new head
	head := state.snake[0]
	newHead := Point{
		X: head.X + state.direction.X,
		Y: head.Y + state.direction.Y,
	}

	// CHECK 1: Wall collision
	if newHead.X < 0 || newHead.X >= gridWidth ||
		newHead.Y < 0 || newHead.Y >= gridHeight {
		state.gameOver = true
		return
	}

	// CHECK 2: Self collision - check if new head hits any body segment
	for _, segment := range state.snake {
		if newHead == segment {
			state.gameOver = true
			return
		}
	}

	// Add new head
	state.snake = append([]Point{newHead}, state.snake...)

	// CHECK 3: Food collision
	if newHead == state.food {
		// Snake grows (don't remove tail)
		state.score += 10
		g.spawnFood()
		// Note: we don't remove the tail, so the snake grows by 1
	} else {
		// Normal movement (remove tail)
		state.snake = state.snake[:len(state.snake)-1]
	}
}

// spawnFood generates a new food position not occupied by snake
func (g *GameStep5) spawnFood() {
	state := g.state

	// Simple deterministic spawning (better to use rand.Intn in real game)
	// Try positions based on current score until we find an empty one
	for attempt := 0; attempt < 100; attempt++ {
		x := (state.score/10 + 7 + attempt) % gridWidth
		y := (state.score/10 + 3 + attempt) % gridHeight
		candidate := Point{X: x, Y: y}

		// Check if position is empty (not on snake)
		empty := true
		for _, segment := range state.snake {
			if segment == candidate {
				empty = false
				break
			}
		}

		if empty {
			state.food = candidate
			return
		}
	}

	// Fallback (shouldn't reach here with a 20x20 grid)
	state.food = Point{X: 0, Y: 0}
}

func (g *GameStep5) Draw(screen *ebiten.Image) {
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

	// Draw food (red)
	foodX := g.state.food.X * cellSize
	foodY := g.state.food.Y * cellSize
	drawCell(screen, foodX, foodY, color.RGBA{R: 255, G: 0, B: 0, A: 255})
}

func (g *GameStep5) Layout(outsideWidth, outsideHeight int) (int, int) {
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
	game := NewGameStep5()

	ebiten.SetWindowTitle("Snake Game - Step 5: Collision & Food")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

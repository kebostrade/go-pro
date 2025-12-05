// +build ignore
// STEP 7: Difficulty Progression
//
// Increase game difficulty as the score increases.
// This introduces:
// - Progressive difficulty
// - Dynamic speed scaling
// - Game balance tuning
//
// Run with: go run step7_difficulty.go

package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

type GameStateStep7 struct {
	snake       []Point
	food        Point
	direction   Point
	nextDir     Point
	tickCount   int
	tickSpeed   int // Gets smaller (faster) as difficulty increases
	gameOver    bool
	gameStarted bool
	score       int
}

type GameStep7 struct {
	state *GameStateStep7
}

func NewGameStep7() *GameStep7 {
	return &GameStep7{
		state: &GameStateStep7{
			snake: []Point{
				{X: 10, Y: 10},
				{X: 9, Y: 10},
				{X: 8, Y: 10},
			},
			direction:   Point{X: 1, Y: 0},
			nextDir:     Point{X: 1, Y: 0},
			tickSpeed:   8,    // Start relatively slow
			gameOver:    false,
			gameStarted: false,
			score:       0,
			food:        Point{X: 15, Y: 10},
		},
	}
}

func (g *GameStep7) Update() error {
	state := g.state

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

	if state.gameOver && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		*state = *NewGameStep7().state
	}

	return nil
}

func (g *GameStep7) handleInput() {
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

func (g *GameStep7) moveSnake() {
	state := g.state

	if state.nextDir != state.direction {
		isOpposite := (state.direction.X == -state.nextDir.X &&
			state.direction.Y == -state.nextDir.Y)
		if !isOpposite {
			state.direction = state.nextDir
		}
	}

	head := state.snake[0]
	newHead := Point{
		X: head.X + state.direction.X,
		Y: head.Y + state.direction.Y,
	}

	// Wall collision
	if newHead.X < 0 || newHead.X >= gridWidth ||
		newHead.Y < 0 || newHead.Y >= gridHeight {
		state.gameOver = true
		return
	}

	// Self collision
	for _, segment := range state.snake {
		if newHead == segment {
			state.gameOver = true
			return
		}
	}

	state.snake = append([]Point{newHead}, state.snake...)

	if newHead == state.food {
		state.score += 10
		g.spawnFood()

		// DIFFICULTY SCALING
		// Increase difficulty every 50 points
		// Don't go faster than tickSpeed of 4 (too hard)
		if state.score%50 == 0 && state.tickSpeed > 4 {
			state.tickSpeed--
		}
	} else {
		state.snake = state.snake[:len(state.snake)-1]
	}
}

func (g *GameStep7) spawnFood() {
	state := g.state

	for attempt := 0; attempt < 100; attempt++ {
		x := (state.score/10 + 7 + attempt) % gridWidth
		y := (state.score/10 + 3 + attempt) % gridHeight
		candidate := Point{X: x, Y: y}

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

	state.food = Point{X: 0, Y: 0}
}

func (g *GameStep7) Draw(screen *ebiten.Image) {
	state := g.state

	screen.Fill(ebiten.ColorBlack)

	// Draw snake
	for i, segment := range state.snake {
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

	// Draw food
	foodX := state.food.X * cellSize
	foodY := state.food.Y * cellSize
	drawCell(screen, foodX, foodY, color.RGBA{R: 255, G: 0, B: 0, A: 255})

	// Draw UI
	if !state.gameStarted {
		ebitenutil.DebugPrintAt(screen, "SNAKE GAME", 10, 10)
		ebitenutil.DebugPrintAt(screen, "Press SPACE to start", 10, 30)
	} else if state.gameOver {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("GAME OVER - Score: %d", state.score), 10, 10)
		ebitenutil.DebugPrintAt(screen, "Press R to restart", 10, 30)
	} else {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf("Score: %d | Speed: %d", state.score, 9-state.tickSpeed), 10, 10)
	}
}

func (g *GameStep7) Layout(outsideWidth, outsideHeight int) (int, int) {
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
	game := NewGameStep7()

	ebiten.SetWindowTitle("Snake Game - Step 7: Difficulty")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

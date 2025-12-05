// This file is not compiled when main.go is present.
// To use this step, rename main.go or run: go run step1_window.go
//
// STEP 1: Basic Window and Game Loop
//
// In this step, we create a basic Ebiten game window and implement
// the required Game interface methods.
//
// Key concepts:
// - Implementing the ebiten.Game interface
// - Game loop basics (Update and Draw)
// - Window configuration
//
// Run with: go run step1_window.go
//
// +build ignore

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

// Game struct implements the ebiten.Game interface
type GameStep1 struct{}

// Update is called every frame (60 times per second by default)
// This is where you put game logic
func (g *GameStep1) Update() error {
	// For now, we don't have any game logic
	return nil
}

// Draw is called every frame to render graphics
// The screen parameter is what we draw to
func (g *GameStep1) Draw(screen *ebiten.Image) {
	// Fill the entire screen with black
	screen.Fill(ebiten.ColorBlack)
	// Try other colors: ebiten.ColorWhite, ebiten.ColorRed, etc.
}

// Layout defines the screen dimensions
// It's called when the window is resized
func (g *GameStep1) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &GameStep1{}

	// Set window properties
	ebiten.SetWindowTitle("Snake Game - Step 1: Basic Window")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	// ebiten.SetWindowResizable(true) // Allow window resizing
	// ebiten.SetWindowFloating(true)  // Always on top

	// Run the game (this blocks until window closes)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

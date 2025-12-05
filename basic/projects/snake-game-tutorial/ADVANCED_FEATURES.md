# Advanced Features for Snake Game

This guide provides implementations for extending the basic Snake game with advanced features. Each section is self-contained and can be integrated into the main game.

## Table of Contents

1. [Using Image Assets](#using-image-assets)
2. [Sound Effects and Music](#sound-effects-and-music)
3. [Pause Functionality](#pause-functionality)
4. [High Score Persistence](#high-score-persistence)
5. [Game Modes](#game-modes)
6. [Obstacles](#obstacles)
7. [Power-ups](#power-ups)
8. [Particle Effects](#particle-effects)
9. [Multi-player Local](#multi-player-local)
10. [Performance Optimization](#performance-optimization)

---

## Using Image Assets

Replace pixel manipulation with pre-rendered images for better performance and visuals.

### Option 1: Embedded Images

```go
import (
    _ "image/png"
    "image/color"
    "github.com/hajimehoshi/ebiten/v2"
)

var (
    snakeHeadImage *ebiten.Image
    snakeBodyImage *ebiten.Image
    foodImage      *ebiten.Image
)

func init() {
    // Create simple colored images
    snakeHeadImage, _ = ebiten.NewImage(20, 20, nil)
    snakeHeadImage.Fill(color.RGBA{R: 0, G: 255, B: 0, A: 255})

    snakeBodyImage, _ = ebiten.NewImage(20, 20, nil)
    snakeBodyImage.Fill(color.RGBA{R: 0, G: 200, B: 0, A: 255})

    foodImage, _ = ebiten.NewImage(20, 20, nil)
    foodImage.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
}

func (g *Game) Draw(screen *ebiten.Image) {
    screen.Fill(ebiten.ColorBlack)

    // Draw snake with images
    for i, segment := range g.state.snake {
        opts := &ebiten.DrawImageOptions{}
        opts.GeoM.Translate(
            float64(segment.X*cellSize),
            float64(segment.Y*cellSize),
        )

        img := snakeBodyImage
        if i == 0 {
            img = snakeHeadImage
        }

        screen.DrawImage(img, opts)
    }

    // Draw food
    opts := &ebiten.DrawImageOptions{}
    opts.GeoM.Translate(
        float64(g.state.food.X*cellSize),
        float64(g.state.food.Y*cellSize),
    )
    screen.DrawImage(foodImage, opts)
}
```

### Option 2: Load PNG Files

```go
import (
    "os"
    "image"
    _ "image/png"
    "github.com/hajimehoshi/ebiten/v2"
)

func loadImage(filepath string) (*ebiten.Image, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    img, _, err := image.Decode(file)
    if err != nil {
        return nil, err
    }

    return ebiten.NewImageFromImage(img)
}

func init() {
    var err error
    snakeHeadImage, err = loadImage("assets/snake_head.png")
    if err != nil {
        log.Fatal(err)
    }
    // Load other images...
}
```

### Asset Organization

```
snake-game-tutorial/
├── assets/
│   ├── snake_head.png
│   ├── snake_body.png
│   ├── food.png
│   └── background.png
├── main.go
└── ...
```

---

## Sound Effects and Music

Add audio feedback for player actions.

### Setup

```go
import (
    "github.com/hajimehoshi/ebiten/v2/audio"
    "github.com/hajimehoshi/ebiten/v2/audio/wav"
)

type GameState struct {
    // ... existing fields ...
    audioContext *audio.Context
    eatSound     *audio.Player
    crashSound   *audio.Player
    bgmPlayer    *audio.Player
}

func init() {
    audioContext := audio.NewContext(44100)

    // Load and play background music
    bgm, _ := loadAudioFile("assets/bgm.wav")
    bgmPlayer, _ := audioContext.NewPlayer(bgm)
    bgmPlayer.SetVolume(0.3)
    bgmPlayer.Play()
}

func loadAudioFile(filepath string) (io.ReadSeeker, error) {
    return os.Open(filepath)
}
```

### Playing Sound Effects

```go
func (g *Game) onFoodEaten() {
    state := g.state
    state.score += 10

    // Play eat sound
    g.state.eatSound.Rewind()
    g.state.eatSound.Play()

    g.spawnFood()
}

func (g *Game) onGameOver() {
    state := g.state
    state.gameOver = true

    // Play crash sound
    g.state.crashSound.Rewind()
    g.state.crashSound.Play()

    // Stop background music
    g.state.bgmPlayer.Pause()
}
```

### Asset Files

```
assets/
├── eat.wav          # Food eaten sound
├── crash.wav        # Collision sound
├── gameover.wav     # Game over sound
└── bgm.wav          # Background music
```

---

## Pause Functionality

Allow players to pause the game.

```go
type GameState struct {
    // ... existing fields ...
    paused bool
}

func (g *Game) Update() error {
    state := g.state

    // Toggle pause with P key
    if inpututil.IsKeyJustPressed(ebiten.KeyP) {
        state.paused = !state.paused
    }

    // Don't update game logic when paused
    if state.paused {
        return nil
    }

    // ... rest of update logic ...
}

func (g *Game) Draw(screen *ebiten.Image) {
    // ... normal drawing ...

    if g.state.paused {
        // Draw pause overlay
        screen.Fill(color.RGBA{0, 0, 0, 128})
        ebitenutil.DebugPrintAt(screen, "PAUSED - Press P to resume", 50, 150)
    }
}
```

---

## High Score Persistence

Save and load high scores to disk.

```go
import (
    "encoding/json"
    "os"
)

type ScoreData struct {
    Score int
    Time  string
}

type Leaderboard struct {
    Scores []ScoreData
}

const scoresFile = "highscores.json"

func (g *Game) loadHighScores() error {
    data, err := os.ReadFile(scoresFile)
    if err != nil {
        if os.IsNotExist(err) {
            g.leaderboard = &Leaderboard{}
            return nil
        }
        return err
    }

    return json.Unmarshal(data, g.leaderboard)
}

func (g *Game) saveHighScores() error {
    // Add current score
    g.leaderboard.Scores = append(g.leaderboard.Scores, ScoreData{
        Score: g.state.score,
        Time:  time.Now().Format("2006-01-02 15:04"),
    })

    // Sort scores
    sort.Slice(g.leaderboard.Scores, func(i, j int) bool {
        return g.leaderboard.Scores[i].Score > g.leaderboard.Scores[j].Score
    })

    // Keep only top 10
    if len(g.leaderboard.Scores) > 10 {
        g.leaderboard.Scores = g.leaderboard.Scores[:10]
    }

    // Save to file
    data, _ := json.MarshalIndent(g.leaderboard, "", "  ")
    return os.WriteFile(scoresFile, data, 0644)
}

func (g *Game) onGameOver() {
    g.state.gameOver = true
    g.saveHighScores()
}
```

### Display High Scores

```go
func (g *Game) drawHighScores(screen *ebiten.Image) {
    y := 50
    ebitenutil.DebugPrintAt(screen, "HIGH SCORES", 50, y)
    y += 20

    for i, score := range g.leaderboard.Scores {
        text := fmt.Sprintf("%d. %d points - %s", i+1, score.Score, score.Time)
        ebitenutil.DebugPrintAt(screen, text, 50, y)
        y += 20
    }
}
```

---

## Game Modes

Implement different gameplay variations.

### Timed Mode

```go
type GameMode int

const (
    ClassicMode GameMode = iota
    TimedMode
    SurvivalMode
)

type GameState struct {
    // ... existing fields ...
    mode          GameMode
    timeRemaining int // In frames
    initialTime   int // 30 seconds = 1800 frames at 60 FPS
}

func (g *Game) Update() error {
    state := g.state

    if state.gameStarted && !state.gameOver && state.mode == TimedMode {
        state.timeRemaining--

        if state.timeRemaining <= 0 {
            state.gameOver = true
        }
    }

    // ... rest of update ...
}

func (g *Game) Draw(screen *ebiten.Image) {
    // ... normal drawing ...

    if g.state.gameStarted && g.state.mode == TimedMode {
        seconds := g.state.timeRemaining / 60
        ebitenutil.DebugPrintAt(screen,
            fmt.Sprintf("Time: %d", seconds), 10, 30)
    }
}
```

### Survival Mode

```go
// In survival mode, food doesn't grow the snake
// Instead, the snake loses length over time

func (g *Game) moveSnake() {
    state := g.state

    // ... collision checks ...

    state.snake = append([]Point{newHead}, state.snake...)

    if state.mode == SurvivalMode {
        // Always remove tail (no growth)
        state.snake = state.snake[:len(state.snake)-1]

        // Lose length if don't eat often
        // Use a timer to remove segments
    } else {
        if newHead == state.food {
            state.score += 10
            g.spawnFood()
        } else {
            state.snake = state.snake[:len(state.snake)-1]
        }
    }
}
```

---

## Obstacles

Add static obstacles to the map.

```go
type GameState struct {
    // ... existing fields ...
    obstacles []Point
}

func (g *Game) initializeLevel() {
    // Create a grid pattern of obstacles
    obstacles := []Point{}

    for x := 5; x < gridWidth; x += 5 {
        for y := 5; y < gridHeight; y += 5 {
            obstacles = append(obstacles, Point{X: x, Y: y})
        }
    }

    g.state.obstacles = obstacles
}

func (g *Game) moveSnake() {
    state := g.state

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

    // OBSTACLE COLLISION - NEW
    for _, obstacle := range state.obstacles {
        if newHead == obstacle {
            state.gameOver = true
            return
        }
    }

    // ... rest of movement ...
}

func (g *Game) Draw(screen *ebiten.Image) {
    // ... existing drawing ...

    // Draw obstacles (gray)
    for _, obstacle := range g.state.obstacles {
        x := obstacle.X * cellSize
        y := obstacle.Y * cellSize
        drawCell(screen, x, y, color.RGBA{R: 100, G: 100, B: 100, A: 255})
    }
}
```

---

## Power-ups

Implement temporary beneficial effects.

```go
type PowerUpType int

const (
    SpeedBoost PowerUpType = iota
    Shield
    GrowthSpurt
)

type PowerUp struct {
    Position Point
    Type     PowerUpType
    ExpiresAt int // Frame when it disappears
}

type GameState struct {
    // ... existing fields ...
    powerUps []PowerUp
    shielded bool
    shieldExpiresAt int
}

func (g *Game) spawnPowerUp() {
    // Random chance to spawn
    if rand.Float64() < 0.001 {
        powerUp := PowerUp{
            Position: Point{
                X: rand.Intn(gridWidth),
                Y: rand.Intn(gridHeight),
            },
            Type: PowerUpType(rand.Intn(3)),
            ExpiresAt: g.frameCount + 300, // 5 seconds at 60 FPS
        }
        g.state.powerUps = append(g.state.powerUps, powerUp)
    }
}

func (g *Game) checkPowerUpCollision(head Point) {
    state := g.state

    for i, powerUp := range state.powerUps {
        if head == powerUp.Position {
            g.activatePowerUp(powerUp.Type)

            // Remove power-up
            state.powerUps = append(state.powerUps[:i], state.powerUps[i+1:]...)
            break
        }
    }
}

func (g *Game) activatePowerUp(powerUpType PowerUpType) {
    state := g.state

    switch powerUpType {
    case SpeedBoost:
        // Double speed for 3 seconds
        if state.tickSpeed > 2 {
            state.tickSpeed--
        }

    case Shield:
        // Protect from one collision
        state.shielded = true
        state.shieldExpiresAt = g.frameCount + 300

    case GrowthSpurt:
        // Grow snake by 5 segments
        tail := state.snake[len(state.snake)-1]
        for i := 0; i < 5; i++ {
            state.snake = append(state.snake, tail)
        }
    }
}

func (g *Game) moveSnake() {
    // ... collision checks ...

    // Self collision with shield
    if g.state.shielded {
        for _, segment := range g.state.snake {
            if newHead == segment {
                g.state.shielded = false
                return // Don't die, just lose shield
            }
        }
    }

    // ... rest of movement ...
}

func (g *Game) Draw(screen *ebiten.Image) {
    // ... existing drawing ...

    // Draw power-ups
    for _, powerUp := range g.state.powerUps {
        x := powerUp.Position.X * cellSize
        y := powerUp.Position.Y * cellSize

        var color color.Color
        switch powerUp.Type {
        case SpeedBoost:
            color = color.RGBA{R: 255, G: 255, B: 0, A: 255} // Yellow
        case Shield:
            color = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue
        case GrowthSpurt:
            color = color.RGBA{R: 255, G: 0, B: 255, A: 255} // Magenta
        }

        drawCell(screen, x, y, color)
    }

    // Draw shield indicator
    if g.state.shielded {
        ebitenutil.DebugPrintAt(screen, "SHIELD ACTIVE", 10, 50)
    }
}
```

---

## Particle Effects

Add visual feedback with particles.

```go
type Particle struct {
    Position Point
    Velocity Point
    LifeSpan int
    Color    color.Color
}

type GameState struct {
    // ... existing fields ...
    particles []Particle
}

func (g *Game) createParticles(x, y int, color color.Color) {
    for i := 0; i < 8; i++ {
        angle := float64(i) * 2 * math.Pi / 8

        g.state.particles = append(g.state.particles, Particle{
            Position: Point{X: x, Y: y},
            Velocity: Point{
                X: int(math.Cos(angle) * 2),
                Y: int(math.Sin(angle) * 2),
            },
            LifeSpan: 30,
            Color:    color,
        })
    }
}

func (g *Game) Update() error {
    // ... existing update ...

    // Update particles
    for i := 0; i < len(g.state.particles); i++ {
        g.state.particles[i].Position.X += g.state.particles[i].Velocity.X
        g.state.particles[i].Position.Y += g.state.particles[i].Velocity.Y
        g.state.particles[i].LifeSpan--

        if g.state.particles[i].LifeSpan <= 0 {
            g.state.particles = append(
                g.state.particles[:i],
                g.state.particles[i+1:]...)
            i--
        }
    }

    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // ... existing drawing ...

    // Draw particles
    for _, particle := range g.state.particles {
        alpha := uint8(255 * particle.LifeSpan / 30)
        c := particle.Color
        c.A = alpha
        drawCell(screen, particle.Position.X, particle.Position.Y, c)
    }
}
```

---

## Multi-player Local

Add local 2-player support.

```go
type Player struct {
    snake     []Point
    food      Point
    direction Point
    score     int
    gameOver  bool
}

type GameState struct {
    player1 Player
    player2 Player
    mode    GameMode
}

func (g *Game) Update() error {
    // Player 1 controls: WASD
    if ebiten.IsKeyPressed(ebiten.KeyW) && g.state.player1.direction.Y == 0 {
        g.state.player1.direction = Point{X: 0, Y: -1}
    }
    // ... WASD controls ...

    // Player 2 controls: Arrow keys
    if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.state.player2.direction.Y == 0 {
        g.state.player2.direction = Point{X: 0, Y: -1}
    }
    // ... Arrow key controls ...

    // Update both players
    g.movePlayer(&g.state.player1)
    g.movePlayer(&g.state.player2)

    // Check if snakes collide
    g.checkSnakeCollision()

    return nil
}

func (g *Game) checkSnakeCollision() {
    head1 := g.state.player1.snake[0]
    head2 := g.state.player2.snake[0]

    // Heads collide
    if head1 == head2 {
        g.state.player1.gameOver = true
        g.state.player2.gameOver = true
    }

    // Head1 hits player2's body
    for _, segment := range g.state.player2.snake {
        if head1 == segment {
            g.state.player1.gameOver = true
        }
    }

    // Head2 hits player1's body
    for _, segment := range g.state.player1.snake {
        if head2 == segment {
            g.state.player2.gameOver = true
        }
    }
}
```

---

## Performance Optimization

Optimize for smooth gameplay.

### 1. Object Pooling

```go
type Snake Pool for reusing segments
type SegmentPool struct {
    available []Point
    inUse     []Point
}

func (p *SegmentPool) Get() Point {
    if len(p.available) > 0 {
        point := p.available[len(p.available)-1]
        p.available = p.available[:len(p.available)-1]
        p.inUse = append(p.inUse, point)
        return point
    }
    return Point{}
}
```

### 2. Caching Pre-calculated Values

```go
type GameState struct {
    snakeCellPixels []ebiten.Image // Cache rendered cells

}

func (g *Game) cacheSnakeCells() {
    for i, segment := range g.state.snake {
        x := segment.X * cellSize
        y := segment.Y * cellSize
        // Cache position calculations
    }
}
```

### 3. Limiting Collision Checks

```go
func (g *Game) moveSnake() {
    // Only check nearby segments instead of all
    head := g.state.snake[0]

    // Check only segments within 5 cells
    for _, segment := range g.state.snake[1:10] {
        if head.nearby(segment) {
            g.state.gameOver = true
        }
    }
}
```

### 4. Frame Skip

```go
func (g *Game) Update() error {
    if ebiten.IsDrawingSkipped() {
        return nil // Skip rendering this frame if needed
    }

    // ... normal update ...
    return nil
}
```

---

## Putting It All Together

Integrate multiple features into one enhanced version:

```go
type AdvancedGameState struct {
    // Core mechanics
    player1, player2 Player
    mode             GameMode

    // Visual features
    particles        []Particle
    powerUps         []PowerUp
    obstacles        []Point

    // Audio
    audioContext     *audio.Context

    // UI
    paused           bool
    leaderboard      *Leaderboard

    // Optimization
    frameCount       int
    cachedPositions  map[Point]bool
}
```

---

## Testing Your Features

```go
func TestPowerUpActivation(t *testing.T) {
    game := NewGame()
    game.state.player1.food = Point{X: 10, Y: 10}

    // Simulate eating power-up
    game.activatePowerUp(SpeedBoost)

    if game.state.player1.tickSpeed == game.initialTickSpeed {
        t.Fatal("PowerUp didn't activate")
    }
}
```

---

## Next Steps

1. Pick one feature and implement it fully
2. Test thoroughly before adding more
3. Optimize based on performance measurements
4. Consider creating a settings menu to enable/disable features
5. Package your enhanced game and share it!

Good luck! 🎮

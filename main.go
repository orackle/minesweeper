package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
  "os"
	"math/rand"
	"time"
)

// Define parameters
const (
	screenWidth  = 800
	screenHeight = 600
	tileSize     = 40
	circleRadius = tileSize / 3
	buttonWidth  = 200
	buttonHeight = 50
)

var (
	numRows  = screenHeight / tileSize
	numCols  = screenWidth / tileSize
	numMines = 40
)

type Game struct {
	board        [][]int
	revealed     [][]bool
	flagged      [][]bool
	gameOver     bool
	gameWon      bool
	numberColors map[int]color.Color // Map to store colors for numbers
	restartBtn   Button             // Button for restarting the game
	exitBtn      Button             // Button for exiting the game
	shouldExit   bool               // Flag to indicate if the game should exit
}

type Button struct {
	x, y   float64     // Position
	width  float64     // Width
	height float64     // Height
	label  string      // Text on the button
	action func()      // Action to perform when clicked
	active bool        // Whether the button is active (visible and clickable)
}




// -------------------------------------------------------------------
func (g *Game) initNumberColors() {
	// Define pastel colors for numbers
	pastelBlue := color.RGBA{135, 206, 250, 255}
	pastelGreen := color.RGBA{144, 238, 144, 255}
	pastelYellow := color.RGBA{255, 255, 224, 255}
	pastelPink := color.RGBA{255, 182, 193, 255}
	// Define colors for numbers
	g.numberColors = map[int]color.Color{
		1: pastelBlue,
		2: pastelGreen,
		3: pastelBlue,
		4: pastelYellow,
		5: pastelPink,
		6: colornames.Cyan,
		7: colornames.Orange,
		8: colornames.Gray,
	}
}



// -------------------------------------------------------------------
func (g *Game) restart() {
	// Reset game state
	g.board = make([][]int, numRows)
	g.revealed = make([][]bool, numRows)
	g.flagged = make([][]bool, numRows)

	for i := range g.board {
		g.board[i] = make([]int, numCols)
		g.revealed[i] = make([]bool, numCols)
		g.flagged[i] = make([]bool, numCols)
	}

	// Place mines randomly around the board
	for i := 0; i < numMines; i++ {
		for {
			row := rand.Intn(numRows)
			col := rand.Intn(numCols)
			if g.board[row][col] == 0 {
				g.board[row][col] = -1
				break
			}
		}
	}

	// Calculate number of adjacent mines for each tile
	g.countMines()

	g.gameOver = false
	g.gameWon = false
	g.restartBtn.active = false
	g.exitBtn.active = false
}



// -------------------------------------------------------------------

func (g *Game) Update() error {
	// Game update logic
	if g.gameOver || g.gameWon {
		if !g.restartBtn.active {
			// Show restart and exit buttons
			g.restartBtn = Button{
				x:      screenWidth/2 - buttonWidth/2,
				y:      screenHeight/2 - buttonHeight/2 - 30,
				width:  buttonWidth,
				height: buttonHeight,
				label:  "Restart",
				action: g.restart,
				active: true,
			}
			g.exitBtn = Button{
				x:      screenWidth/2 - buttonWidth/2,
				y:      screenHeight/2 - buttonHeight/2 + 30,
				width:  buttonWidth,
				height: buttonHeight,
				label:  "Exit",
        action: func() { os.Exit(0) }, // Exit the program
				active: true,
			}
		}

		// Handle button clicks
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			if g.restartBtn.active && g.isMouseInButton(x, y, g.restartBtn) {
				g.restartBtn.action()
			} else if g.exitBtn.active && g.isMouseInButton(x, y, g.exitBtn) {
				g.exitBtn.action()
			}
		}

		return nil
	}

	// Handle left mouse button click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		col := x / tileSize
		row := y / tileSize

		// Check if the click is valid within the game board
		if row >= 0 && row < numRows && col >= 0 && col < numCols {
			g.revealTile(row, col)
		}
	}

	// Handle right mouse button click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		col := x / tileSize
		row := y / tileSize

		// Check if the click is valid within the game board
		if row >= 0 && row < numRows && col >= 0 && col < numCols {
			g.flag(row, col)
		}
	}

	// Check if the player has won
	win := true
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			if g.board[row][col] != -1 && !g.revealed[row][col] {
				win = false
				break
			}
		}
	}

	if win {
		g.gameWon = true
	}

	return nil
}




// -------------------------------------------------------------------
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear the screen with white color
	screen.Fill(color.White)

	// Game drawing logic
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			// Calculate circle position
			x := col*tileSize + tileSize/2
			y := row*tileSize + tileSize/2

			if g.revealed[row][col] {
				if g.board[row][col] == -1 {
					drawCircle(screen, float64(x), float64(y), float64(circleRadius), colornames.Red)
				} else {
					// Calculate the color based on the number of adjacent mines
					var tileColor color.Color
					switch g.board[row][col] {
					case 0:
						tileColor = colornames.White
					case 1:
						tileColor = g.numberColors[1]
					case 2:
						tileColor = g.numberColors[2]
					case 3:
						tileColor = g.numberColors[3]
					case 4:
						tileColor = g.numberColors[4]
					case 5:
						tileColor = g.numberColors[5]
					case 6:
						tileColor = g.numberColors[6]
					case 7:
						tileColor = g.numberColors[7]
					case 8:
						tileColor = g.numberColors[8]
					default:
						tileColor = colornames.White
					}
					drawCircle(screen, float64(x), float64(y), float64(circleRadius), tileColor)
				}
			} else {
				drawCircle(screen, float64(x), float64(y), float64(circleRadius), colornames.Lightgray)
			}

			// Draw flag symbol if flagged
			if g.flagged[row][col] {
				drawCircle(screen, float64(x), float64(y), float64(circleRadius), colornames.Orange)
				drawFlagSymbol(screen, float64(x), float64(y), float64(circleRadius), colornames.Black)
			}
		}
	}

	// Draw restart and exit buttons if active
	if g.restartBtn.active {
		drawButton(screen, g.restartBtn)
	}
	if g.exitBtn.active {
		drawButton(screen, g.exitBtn)
	}
}




// -------------------------------------------------------------------
// Function to draw a filled circle using vector.DrawFilledCircle
func drawCircle(screen *ebiten.Image, x, y, radius float64, clr color.Color) {
	vector.DrawFilledCircle(screen, float32(x), float32(y), float32(radius), clr, true)
}



// -------------------------------------------------------------------
// Function to draw a flag symbol
func drawFlagSymbol(screen *ebiten.Image, x, y, size float64, clr color.Color) {
	// Draw lines to create an "X" symbol
	ebitenutil.DrawLine(screen, x-size/2, y-size/2, x+size/2, y+size/2, clr)
	ebitenutil.DrawLine(screen, x-size/2, y+size/2, x+size/2, y-size/2, clr)
}


// -------------------------------------------------------------------
// Function to draw a button
func drawButton(screen *ebiten.Image, btn Button) {
	// Draw button background
	ebitenutil.DrawRect(screen, btn.x, btn.y, btn.width, btn.height, colornames.Lightblue)

	// Draw button label
	textX := btn.x + btn.width/2 - float64(len(btn.label))*5
	textY := btn.y + btn.height/2 + 5
	ebitenutil.DebugPrintAt(screen, btn.label, int(textX), int(textY))
}


// -------------------------------------------------------------------
func (g *Game) isMouseInButton(mx, my int, btn Button) bool {
	return float64(mx) >= btn.x && float64(mx) <= btn.x+btn.width &&
		float64(my) >= btn.y && float64(my) <= btn.y+btn.height
}


// -------------------------------------------------------------------
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}



// -------------------------------------------------------------------
func (g *Game) revealTile(row, col int) {
	// Reveal tile logic
	if row < 0 || row >= numRows || col < 0 || col >= numCols || g.revealed[row][col] {
		return // Invalid click -- outside game board or on an already revealed tile
	}

	g.revealed[row][col] = true

	if g.board[row][col] == -1 {
		g.gameOver = true
		return
	}

	if g.board[row][col] == 0 {
		// Recursive reveal adjacent tiles
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				g.revealTile(row+dr, col+dc)
			}
		}
	}
}


// -------------------------------------------------------------------

func (g *Game) flag(row, col int) {
	// Toggle flagging of tile
	if g.revealed[row][col] {
		return // Cannot flag revealed tile
	}

	g.flagged[row][col] = !g.flagged[row][col]
}


// ------------------------------------------------------------------
func (g *Game) countMines() {
	// Count adjacent mines for each tile
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			if g.board[row][col] == -1 {
				continue
			}
			count := 0

			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					r, c := row+dr, col+dc
					if r >= 0 && r < numRows && c >= 0 && c < numCols && g.board[r][c] == -1 {
						count++
					}
				}
			}
			g.board[row][col] = count
		}
	}
}


// ---------------------------------------------------------------------------------
func main() {
	rand.Seed(time.Now().UnixNano())

	// Initializing the game struct
	game := &Game{
		board:    make([][]int, numRows),
		revealed: make([][]bool, numRows),
		flagged:  make([][]bool, numRows),
	}

	game.initNumberColors() // Initialize number colors map

	for i := range game.board {
		game.board[i] = make([]int, numCols)
		game.revealed[i] = make([]bool, numCols)
		game.flagged[i] = make([]bool, numCols)
	}

	// Place mines randomly around the board
	for i := 0; i < numMines; i++ {
		for {
			row := rand.Intn(numRows)
			col := rand.Intn(numCols)
			if game.board[row][col] == 0 {
				game.board[row][col] = -1
				break
			}
		}
	}

	game.countMines()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Minesweeper")

	// Start the game loop
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

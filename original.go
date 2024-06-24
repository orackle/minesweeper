package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

// Define parameters
const (
	screenWidth  = 800
	screenHeight = 600
	tileSize     = 40
	numRows      = screenHeight / tileSize
	numCols      = screenWidth / tileSize
	numMines     = 40
)

type Game struct {
	board        [][]int
	revealed     [][]bool
	flagged      [][]bool
	gameOver     bool
	gameWon      bool
	startTime    time.Time
	font         font.Face     // Font face for rendering text
	emojiMap     map[int]rune // Map to store emojis for numbers
}

func (g *Game) initEmojiMap() {
	// Define emojis for numbers (adjust as needed)
	g.emojiMap = map[int]rune{
		1: '🍄', // Placeholder - replace with actual emoji
		2: '👀',
		3: '🚨',
		4: '😬',
		5: '👀',
		6: '👀',
		7: '👀',
		8: '👀',
	}
}

func (g *Game) restart() {
	g.board = make([][]int, numRows)
	g.revealed = make([][]bool, numRows)
	g.flagged = make([][]bool, numRows)

	for i := range g.board {
		g.board[i] = make([]int, numCols)
		g.revealed[i] = make([]bool, numCols)
		g.flagged[i] = make([]bool, numCols)
	}

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

	g.countMines()

	g.gameOver = false
	g.gameWon = false
}

func (g *Game) Update() error {
	if g.gameOver || g.gameWon {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.restart()
		}
		return nil
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		row := y / tileSize
		col := x / tileSize

		if row >= 0 && row < numRows && col >= 0 && col < numCols {
			g.revealTile(row, col)
		}
	}

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

func (g *Game) Draw(screen *ebiten.Image) {
	for row := 0; row < numRows; row++ {
		for col := 0; col < numCols; col++ {
			x := col * tileSize
			y := row * tileSize

			if g.revealed[row][col] {
				if g.board[row][col] == -1 {
					ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, colornames.Red)
				} else {
					ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, colornames.Lightgray)
					if g.board[row][col] > 0 {
						// Calculate text position
						textWidth := tileSize / 2 // Estimate text width
						textX := x + (tileSize/2 - textWidth/2)
						textY := y + tileSize/2

						// Get emoji rune based on the number of adjacent mines
						emoji := g.emojiMap[g.board[row][col]]

						// Draw emoji using Noto Emoji font
						text.Draw(screen, string(emoji), g.font, textX, textY, colornames.White)
					}
				}
			} else {
				ebitenutil.DrawRect(screen, float64(x), float64(y), tileSize, tileSize, colornames.Darkgray)
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) revealTile(row, col int) {
	if row < 0 || row >= numRows || col < 0 || col >= numCols || g.revealed[row][col] {
		return // Invalid click -- outside game board or on an already revealed tile
	}

	g.revealed[row][col] = true

	if g.board[row][col] == -1 {
		g.gameOver = true
		return
	}

	if g.board[row][col] == 0 {
		for dr := -1; dr <= 1; dr++ {
			for dc := -1; dc <= 1; dc++ {
				g.revealTile(row+dr, col+dc)
			}
		}
	}
}

func (g *Game) countMines() {
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

func main() {
	rand.Seed(time.Now().UnixNano())

	// Load the Noto Emoji font
	fontData, err := ioutil.ReadFile("path/to/NotoColorEmoji.ttf") // Replace with your path
	if err != nil {
		log.Fatal(err)
	}

	tt, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	face := truetype.NewFace(tt, &truetype.Options{
		Size:    20,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	// Initializing the game struct
	game := &Game{
		board:    make([][]int, numRows),
		revealed: make([][]bool, numRows),
		flagged:  make([][]bool, numRows),
		font:     face, // Assigning the font face
	}

	game.initEmojiMap() // Initialize emoji map

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

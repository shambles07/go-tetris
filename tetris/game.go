// This is a simple implementation of a console-based Tetris clone. It only works where termbox-go works
// (linux/mac should be fine).
package tetris

import (
	"github.com/nsf/termbox-go"
	"math"
	"math/rand"
	"time"
)

type Direction int

const (
	Up Direction = iota + 1
	Down
	Left
	Right
)

// A Game tracks the entire game state of tetris, including the Board, the upcoming piece, the game speed
// (dropDelayMillis), the score, and various other internal data.
type Game struct {
	board           *Board
	nextPiece       *Piece
	pieces          []Piece
	paused          bool
	over            bool
	dropDelayMillis int
	ticker          *time.Ticker
	score           int
}

// Initialize a new game, ready to be started with Start().
func NewGame() *Game {
	game := new(Game)
	game.pieces = tetrisPieces()
	game.board = newBoard()
	game.board.currentPiece = game.GeneratePiece()
	game.board.currentPosition = Vector{initialX, 0}
	game.nextPiece = game.GeneratePiece()
	game.paused = false
	game.over = false
	game.score = 0
	game.startTicker()
	return game
}

// Start up the ticker with the appropriate interval based on the current score.
func (game *Game) startTicker() {
	// Set the speed as a function of score. Starts at 800ms, decreases to 200ms by 100ms each 500 points.
	game.dropDelayMillis = 800 - game.score/5
	if game.dropDelayMillis < 200 {
		game.dropDelayMillis = 200
	}
	game.ticker = time.NewTicker(time.Duration(game.dropDelayMillis) * time.Millisecond)
}

// Stop the game ticker. This stops automatic advancement of the piece.
func (game *Game) stopTicker() {
	game.ticker.Stop()
}

// A game event, generated by user input or by the game ticker.
type GameEvent int

const (
	MoveLeft GameEvent = iota
	MoveRight
	MoveDown
	Rotate
	QuickDrop
	Pause
	Quit
	// An event that doesn't cause a change to game state but causes a full redraw; e.g., a window resize.
	Redraw
)

// Start running the game. It will continue indefinitely until the user exits.
func (game *Game) Start() {

	drawStaticBoardParts()
	game.DrawDynamic(false)

	eventQueue := make(chan GameEvent, 100)
	go func() {
		for {
			eventQueue <- waitForUserEvent()
		}
	}()
gameOver:
	for {
		var event GameEvent
		select {
		case event = <-eventQueue:
		case <-game.ticker.C:
			event = MoveDown
		}
		// If the game is paused, ignore all commands except for Pause, Quit, and Redraw. On Redraw, redraw
		// the pause screen.
		if game.paused {
			switch event {
			case Pause:
				game.PauseToggle()
			case Quit:
				return
			case Redraw:
				drawStaticBoardParts()
				game.DrawPauseScreen()
			}
		} else {
			switch event {
			case MoveLeft:
				game.Move(Left)
			case MoveRight:
				game.Move(Right)
			case MoveDown:
				game.Move(Down)
			case QuickDrop:
				game.QuickDrop()
			case Rotate:
				game.Rotate()
			case Pause:
				game.PauseToggle()
			case Quit:
				return
			case Redraw:
				drawStaticBoardParts()
			}
			// Update screen only if game is not now paused.
			if !game.paused {
				game.DrawDynamic(false)
			}
		}
		if game.over {
			break gameOver
		}
	}
	game.DrawGameOver()
	for waitForUserEvent() != Quit {
	}
}

// A blocking function that waits on a ticker and then emits a MoveDown event.
func waitForTick(ticker *time.Ticker) GameEvent {
	<-ticker.C
	return MoveDown
}

// A blocking function that waits for user input and then emits the appropriate GameEvent.
func waitForUserEvent() GameEvent {
	switch event := termbox.PollEvent(); event.Type {
	// Movement: arrow keys or vim controls (h, j, k, l)
	// Pause: 'p'
	// Exit: 'q' or ctrl-c.
	case termbox.EventKey:
		if event.Ch == 0 { // A special key combo was pressed
			switch event.Key {
			case termbox.KeyCtrlC:
				return Quit
			case termbox.KeyArrowLeft:
				return MoveLeft
			case termbox.KeyArrowUp:
				return Rotate
			case termbox.KeyArrowRight:
				return MoveRight
			case termbox.KeyArrowDown:
				return MoveDown
			case termbox.KeySpace:
				return QuickDrop
			}
		} else {
			switch event.Ch {
			case 'p':
				return Pause
			case 'q':
				return Quit
			case 'h':
				return MoveLeft
			case 'k':
				return Rotate
			case 'l':
				return MoveRight
			case 'j':
				return MoveDown
			}
		}
	case termbox.EventResize:
		return Redraw
	case termbox.EventError:
		panic(event.Err)
	}
	return Redraw // Should never be reached
}

// Randomly choose a new game piece from among the the available pieces.
func (game *Game) GeneratePiece() *Piece {
	return &game.pieces[rand.Intn(len(game.pieces))]
}

// Anchor the current piece to the board, clears lines, increments the score, and generate a new piece. Sets
// the 'game over' state if the new piece overlaps existing pieces.
func (game *Game) anchor() {
	game.board.mergeCurrentPiece()

	// Clear any completed rows (with animation) and increment the score if necessary.
	rowsCleared := game.board.clearedRows()

	if len(rowsCleared) > 0 {
		// Animate the cleared rows disappearing
		game.stopTicker()
		flickerCells := make(map[Vector]termbox.Attribute)
		for _, y := range rowsCleared {
			for x := 0; x < width; x++ {
				point := Vector{x, y}
				flickerCells[point] = game.board.cells[point]
			}
		}
		for i := 0; i < 5; i++ {
			for point, color := range flickerCells {
				if i%2 == 0 {
					color = backgroundColor
				}
				setBoardCell((point.x*2)+2, headerHeight+point.y+2, color)
			}
			termbox.Flush()
			time.Sleep(80 * time.Millisecond)
		}

		// Get rid of the rows
		game.board.clearRows()

		// Scoring -- 1 row -> 100, 2 rows -> 200, ... 4 rows -> 800
		points := 100 * math.Pow(2, float64(len(rowsCleared)-1))
		game.score += int(points)

		game.startTicker()
	}

	// Bring in the next piece.
	game.board.currentPiece = game.nextPiece
	game.board.currentPosition = Vector{initialX, 0}
	game.nextPiece = game.GeneratePiece()
	game.nextPiece.currentRotation = 0

	if game.board.currentPieceInCollision() {
		game.over = true
	}
}

// Attempt to move.
func (game *Game) Move(where Direction) {
	translation := Vector{0, 0}
	switch where {
	case Down:
		translation = Vector{0, 1}
	case Left:
		translation = Vector{-1, 0}
	case Right:
		translation = Vector{1, 0}
	}
	// Attempt to make the move.
	moved := game.board.moveIfPossible(translation)

	// Perform anchoring if we tried to move down but we were unsuccessful.
	if where == Down && !moved {
		game.anchor()
	}
}

// Drop the piece all the way and anchor it.
func (game *Game) QuickDrop() {
	// Move down as far as possible
	for game.board.moveIfPossible(Vector{0, 1}) {
	}
	game.DrawDynamic(false)
	game.anchor()
}

// Rotates the current game piece, if possible.
func (game *Game) Rotate() {
	game.board.currentPiece.rotate()
	if game.board.currentPieceInCollision() {
		game.board.currentPiece.unrotate()
	}
}

// Draw the dynamic parts of the game interface (the board, the next piece preview pane, and the score).  The
// static parts should be drawn with the drawStaticBoardParts() function, if needed.  If clearOnly is true, 
// the board and preview pane will be cleared rather than redrawn. 
func (game *Game) DrawDynamic(clearOnly bool) {

	// Print the board contents. Each block will correspond to a side-by-side pair of cells in the termbox, so
	// that the visible blocks will be roughly square.  If clearOnly is true, draw background color.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if clearOnly {
				setBoardCell((x*2)+2, headerHeight+y+2, backgroundColor)
			} else {
				color := game.board.CellColor(Vector{x, y})
				setBoardCell((x*2)+2, headerHeight+y+2, color)
			}
		}
	}

	// Print the preview piece. Need to clear the box first.  Draw next piece only if clearOnly is false
	previewPieceOffset := Vector{(width * 2) + 8, headerHeight + 3}
	for x := 0; x < 6; x++ {
		for y := 0; y < 4; y++ {
			cursor := previewPieceOffset.plus(Vector{x, y})
			setCell(cursor.x, cursor.y, ' ', termbox.ColorDefault)
		}
	}
	if !clearOnly {
		for _, point := range game.nextPiece.rotations[0] {
			cursor := previewPieceOffset.plus(Vector{point.x * 2, point.y})
			setBoardCell(cursor.x, cursor.y, game.nextPiece.color)
		}
	}

	// Draw the current score.  If clearOnly, do the same.
	score := game.score
	cursor := Vector{(width * 2) + 18, headerHeight + previewHeight + 7}
	for {
		digit := score % 10
		score /= 10
		drawDigitAsAscii(cursor.x, cursor.y, digit)
		cursor = cursor.plus(Vector{-4, 0})
		if score == 0 {
			break
		}
	}

	// Flush termbox's internal state to the screen.
	termbox.Flush()
}

// Pause or unpause the game, depending on game.paused.
func (game *Game) PauseToggle() {
	if game.paused {
		drawStaticBoardParts()
		game.DrawDynamic(false)
		game.startTicker()
	} else {
		game.stopTicker()
		game.DrawPauseScreen()
	}
	game.paused = !game.paused
}

// Draw the pause screen, hiding the game board and next piece.
func (game *Game) DrawPauseScreen() {
	// Clear the board and preview screen
	game.DrawDynamic(true)

	// Draw PAUSED overlay
	for y := (totalHeight/2 - 1); y <= (totalHeight/2)+1; y++ {
		for x := 1; x < totalWidth+3; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorBlue)
		}
	}
	for i, ch := range "PAUSED" {
		termbox.SetCell(totalWidth/2-2+i, totalHeight/2, ch, termbox.ColorWhite, termbox.ColorBlue)
	}

	// Flush termbox to screen
	termbox.Flush()
}

// Draw the "GAME OVER" overlay on top of the game interface.
func (game *Game) DrawGameOver() {
	for y := (totalHeight/2 - 1); y <= (totalHeight/2)+1; y++ {
		for x := 1; x < totalWidth+3; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault, termbox.ColorBlue)
		}
	}
	for i, ch := range "GAME OVER" {
		termbox.SetCell(totalWidth/2-4+i, totalHeight/2, ch, termbox.ColorWhite, termbox.ColorBlue)
	}
	termbox.Flush()
}

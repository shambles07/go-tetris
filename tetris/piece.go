package tetris

import (
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

// Bag stores the piece number and pieces
type Bag struct {
	pieceNumber int
	pieces      []int
	bagNumber   int
}

// PieceInstance stores a piece's particular rotational instance
type PieceInstance []Vector

// Piece holds a Tetris piece, including the all the possible rotations, a color, and an index indicating the current
// rotation of the piece.
type Piece struct {
	name            string
	rotations       []PieceInstance
	currentRotation int
	initialLocation Vector
	color           termbox.Attribute
}

// initialize is called after each bag has been
func (b *Bag) initialize() {
	b.pieceNumber = 6 // reset the piece number
	b.shuffle()       // randomize this bag
}

func (b *Bag) shuffle() {
	b.pieces = []int{0, 1, 2, 3, 4, 5, 6} // piece set is constant
	rand.Seed(time.Now().UnixNano())      // seed rng and shuffle the bag!
	rand.Shuffle(len(b.pieces), func(i, j int) { b.pieces[i], b.pieces[j] = b.pieces[j], b.pieces[i] })
}

func (b *Bag) getPiece() (retPiece int) {
	// call to reset pieceNumber and reshuffle after every bag
	if b.pieceNumber < 0 {
		b.initialize()
		// increment our bag number
		b.bagNumber++
	}
	retPiece = b.pieces[b.pieceNumber]
	// decrement piece number in bag and return piece
	b.pieceNumber--
	return retPiece
}

/* func (p *Piece) hold() PieceInstance {
	return p.name
} */

// Find the current PieceInstance of this piece.
func (p *Piece) instance() PieceInstance {
	return p.rotations[p.currentRotation]
}

// Go to the next rotation.
func (p *Piece) rotate() {
	p.currentRotation = (p.currentRotation + 1) % len(p.rotations)
}

// Go to the previous rotation.
func (p *Piece) unrotate() {
	p.currentRotation = (p.currentRotation - 1) % len(p.rotations)
	if p.currentRotation < 0 {
		p.currentRotation += len(p.rotations)
	}
}

// Rotate the piece 180 degrees
func (p *Piece) dblrotate() {
	p.currentRotation = (p.currentRotation + 2) % len(p.rotations)
	if p.currentRotation < 0 {
		p.currentRotation += len(p.rotations)
	}
}

// This has all the hard-coded tetris pieces.
// TODO: It might be nice to have a way to parse these from a configuration file. Maybe the format could look
// something like this:
//
// ##   #
//	## ##
//	   #
//
// ##
// ##
//
// etc.

func tetrisPieces() []Piece {

	return []Piece{
		Piece{"O",
			[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{0, 1}, Vector{1, 1}}},
			0, Vector{4, 0}, termbox.ColorYellow},
		Piece{"Z",
			[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{1, 1}, Vector{2, 1}}, []Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{0, 2}}},
			0, Vector{3, 0}, termbox.ColorRed},
		Piece{"S",
			[]PieceInstance{[]Vector{Vector{1, 0}, Vector{2, 0}, Vector{0, 1}, Vector{1, 1}}, []Vector{Vector{0, 0}, Vector{0, 1}, Vector{1, 1}, Vector{1, 2}}},
			0, Vector{3, 0}, termbox.ColorGreen},
		Piece{"T",
			[]PieceInstance{[]Vector{Vector{0, 0}, Vector{1, 0}, Vector{2, 0}, Vector{1, 1}}, []Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{1, 2}}, []Vector{Vector{1, 0}, Vector{0, 1}, Vector{1, 1}, Vector{2, 1}}, []Vector{Vector{0, 0}, Vector{0, 1}, Vector{1, 1}, Vector{0, 2}}},
			0, Vector{3, 0}, termbox.ColorMagenta},
		Piece{"L",
			[]PieceInstance{[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{0, 2}}, []Vector{Vector{0, 0}, Vector{1, 0}, Vector{1, 1}, Vector{1, 2}}, []Vector{Vector{2, 0}, Vector{0, 1}, Vector{1, 1}, Vector{2, 1}}, []Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{2, 2}}},
			0, Vector{3, -1}, termbox.ColorWhite},
		Piece{"J",
			[]PieceInstance{[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{2, 2}}, []Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{0, 2}}, []Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{0, 0}}, []Vector{Vector{1, 0}, Vector{2, 0}, Vector{1, 1}, Vector{1, 2}}},
			0, Vector{3, -1}, termbox.ColorBlue},
		Piece{"I",
			[]PieceInstance{[]Vector{Vector{0, 1}, Vector{1, 1}, Vector{2, 1}, Vector{3, 1}}, []Vector{Vector{1, 0}, Vector{1, 1}, Vector{1, 2}, Vector{1, 3}}},
			0, Vector{3, -1}, termbox.ColorCyan},
	}
}

package sudoku

// Puzzle represents sudoku's board
type Puzzle struct {
	board []uint16
}

func newPuzzle() *Puzzle {
	p := &Puzzle{
		board: make([]uint16, 81),
	}
	return p
}

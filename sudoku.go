package sudoku

// Puzzle represents sudoku's board
type Puzzle struct {
	board []uint16
}

// New parses string and returns a puzzle, nil if board is incorrect
func New(board string) *Puzzle {
	if len(board) != 81 {
		return nil
	}
	p := parseGrid(board)
	return p
}

func newPuzzle() *Puzzle {
	p := &Puzzle{
		board: make([]uint16, 81),
	}
	return p
}

func parseGrid(grid string) *Puzzle {
	if len(grid) != 81 {
		return nil
	}
	puzzle := newPuzzle()
	for i, c := range grid {
		if '1' <= c && c <= '9' {
			puzzle.board[i] = getMask(int(c) - int('0'))
		}
	}
	return puzzle
}

func getMask(x int) uint16 {
	return 1 << uint16(x)
}

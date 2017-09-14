package sudoku

import "math/bits"

var neighbors [][]int

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

func hasBit(x, mask uint16) bool {
	return (x & mask) == mask
}

func offBit(x, mask uint16) uint16 {
	return x &^ mask
}

func isPower2(x uint16) bool {
	return (x & (x - 1)) == 0
}

func prepare(puzzle *Puzzle) *Puzzle {
	solution := newPuzzle()
	for i := 0; i < 81; i++ {
		solution.board[i] = (getMask(10) - 1) - 1
	}
	for i := 0; i < 81; i++ {
		v := puzzle.board[i]
		if v > 0 && !assign(solution, i, v) {
			return nil
		}
	}
	return solution
}

func assign(puzzle *Puzzle, idx int, mask uint16) bool {
	for i := 1; i <= 9; i++ {
		tmp := getMask(i)
		if tmp != mask && !eliminate(puzzle, idx, tmp) {
			return false
		}
	}
	return true
}

func eliminate(puzzle *Puzzle, idx int, mask uint16) bool {
	value := puzzle.board[idx]
	if !hasBit(value, mask) {
		return true
	}
	value = offBit(value, mask)
	puzzle.board[idx] = value
	if value == 0 {
		return false
	}
	if bits.OnesCount16(value) == 1 {
		for _, i := range neighbors[idx] {
			if !eliminate(puzzle, i, value) {
				return false
			}
		}
	}
	return true
}

func search(puzzle *Puzzle) *Puzzle {
	for i := 0; i < 81; i++ {
		value := puzzle.board[i]
		if isPower2(value) {
			continue
		}

		minSquare := getMin(puzzle)
		value = puzzle.board[minSquare]

		for i := 1; i <= 9; i++ {
			mask := getMask(i)
			if !hasBit(value, mask) {
				continue
			}
			tmp := newPuzzle()
			copy(tmp.board, puzzle.board)

			if assign(tmp, minSquare, mask) {
				res := search(tmp)
				if res != nil {
					return res
				}
			}
		}
		return nil
	}
	return puzzle
}

func getMin(puzzle *Puzzle) int {
	minSquare, minSize := 0, 9
	for j := 0; j < 81; j++ {
		size := bits.OnesCount16(puzzle.board[j])
		if size > 1 && size <= minSize {
			minSquare = j
			minSize = size
		}
	}
	return minSquare
}

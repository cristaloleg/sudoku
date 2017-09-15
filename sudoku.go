package sudoku

import (
	"math/bits"
)

var neighbors [][]int

// Puzzle represents Sudoku's board
type Puzzle struct {
	board []uint16
}

// New parses a string and returns a puzzle, nil if board is incorrect
func New(board string) *Puzzle {
	if len(board) != 81 {
		return nil
	}
	p := parseGrid(board)
	return p
}

// newPuzzle ...
func newPuzzle() *Puzzle {
	p := &Puzzle{
		board: make([]uint16, 81),
	}
	return p
}

// parseGrid ...
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

// prepare ...
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

// assign tries to assign mask for a given tile
// and than eliminates it from neighbors of the tile
func assign(puzzle *Puzzle, idx int, mask uint16) bool {
	for i := 1; i <= 9; i++ {
		tmp := getMask(i)
		if tmp != mask && !eliminate(puzzle, idx, tmp) {
			return false
		}
	}
	return true
}

// eliminate ...
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

// search starts recursive search for a solution
func search(puzzle *Puzzle) *Puzzle {
	for i := 0; i < 81; i++ {
		value := puzzle.board[i]
		if hasOneBit(value) {
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

// getMin returns an element with minimum possible options
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

func init() {
	makeNeighbors()
}

func makeNeighbors() {
	neighbors = make([][]int, 81)
	for i := 0; i < 81; i++ {
		neighbors[i] = make([]int, 20)
		copy(neighbors[i][:4], closest(i))

		k := 4
		for j := 0; j < 10; j++ {
			if i != (i/9)*9+j {
				neighbors[i][k] = (i/9)*9 + j
				k++
			}
		}

		k = 12
		for j := 0; j < 9; j++ {
			if i != i%9+(j*9) {
				neighbors[i][k] = i%9 + j*9
				k++
			}
		}
	}
}

// closest returns 4 elements neighbors for a given one
func closest(i int) []int {
	t, p := i%3, (i/9)%3
	a, b, c, d := 0, 0, 0, 0
	switch t {
	case 0:
		a, b = 1, 2
	case 1:
		a, b = -1, 1
	case 2:
		a, b = -1, -2
	}
	switch p {
	case 0:
		c, d = 1, 2
	case 1:
		c, d = -1, 1
	case 2:
		c, d = -1, -2
	}
	return []int{
		i + c*9 + a,
		i + c*9 + b,
		i + d*9 + a,
		i + d*9 + b,
	}
}

func getMask(x int) uint16 {
	return 1 << uint16(x)
}

func hasOneBit(x uint16) bool {
	return (x & (x - 1)) == 0
}

func hasBit(x, mask uint16) bool {
	return (x & mask) == mask
}

func offBit(x, mask uint16) uint16 {
	return x &^ mask
}

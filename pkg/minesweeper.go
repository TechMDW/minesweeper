package minesweeper

import (
	"fmt"
	"math/rand"
)

// Cell represents a cell in a minesweeper board.
type Cell struct {
	IsMine      bool
	IsRevealed  bool
	IsFlagged   bool
	MinesAround int
}

// Board represents a minesweeper board.
type Board struct {
	Rows     int
	Cols     int
	NumMines int
	Cells    [][]Cell
	Rand     *rand.Rand
}

// privateRand is a private random number generator
// var privateRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// Symbols used to display the board
var (
	sMine      = "X"
	sFlag      = "F"
	sHidden    = "•"
	sSeperator = "  "
)

const (
	dStartIndex = 1
)

type BoardOptions struct {
	Seed int64
}

// NewBoard creates a new board with the given number of rows, columns, and mines.
//
// The board is initialized with all cells hidden and no mines placed.
func NewBoard(rows, cols, numMines int, options *BoardOptions) *Board {
	board := &Board{
		Rows:     rows,
		Cols:     cols,
		NumMines: numMines,
		Cells:    make([][]Cell, rows),

		Rand: rand.New(rand.NewSource(options.Seed)),
	}

	for i := 0; i < rows; i++ {
		board.Cells[i] = make([]Cell, cols)
	}

	board.placeMines()

	return board
}

// placeMines places mines randomly on the board.
func (b *Board) placeMines() {
	for i := 0; i < b.NumMines; i++ {
		for {
			row := b.Rand.Intn(b.Rows)
			col := b.Rand.Intn(b.Cols)

			if !b.Cells[row][col].IsMine {
				b.Cells[row][col].IsMine = true

				b.incrementMinesAround(row, col)

				break
			}
		}
	}
}

// incrementMinesAround increments the MinesAround field of all cells in the
// board that are adjacent to the cell at row, col.
//
// The MinesAround field of the cell at row, col is not incremented.
func (b *Board) incrementMinesAround(row, col int) {
	for r := row - 1; r <= row+1; r++ {
		for c := col - 1; c <= col+1; c++ {
			if r >= 0 && r < b.Rows && c >= 0 && c < b.Cols && !b.Cells[r][c].IsMine {
				b.Cells[r][c].MinesAround++
			}
		}
	}
}

// Reveal recursively reveals all the cells around a cell. If a mine is revealed, it returns true.
// Otherwise, it returns false.
func (b *Board) Reveal(row, col int) bool {
	if row < 0 || row >= b.Rows || col < 0 || col >= b.Cols || b.Cells[row][col].IsRevealed {
		return false
	}

	b.Cells[row][col].IsRevealed = true

	if b.Cells[row][col].IsMine {
		return true
	}

	if b.Cells[row][col].MinesAround == 0 {
		for r := row - 1; r <= row+1; r++ {
			for c := col - 1; c <= col+1; c++ {
				b.Reveal(r, c)
			}
		}
	}

	return false
}

func (b *Board) CellsNonRevealed() int {
	count := 0

	for _, row := range b.Cells {
		for _, cell := range row {
			if !cell.IsRevealed {
				count++
			}
		}
	}

	return count
}

func (b *Board) CellsRevealed() int {
	count := 0

	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.IsRevealed {
				count++
			}
		}
	}

	return count
}

func (b *Board) FlagsCount() int {
	count := 0

	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.IsFlagged {
				count++
			}
		}
	}

	return count
}

func (b *Board) RevealAll() {
	for i := 0; i < b.Rows; i++ {
		for j := 0; j < b.Cols; j++ {
			cell := &b.Cells[i][j]

			if cell.IsMine {
				cell.IsFlagged = true
				continue
			}

			cell.IsRevealed = true
		}
	}
}

func (b *Board) RevealedPercentage() float64 {
	nonMineCells := 0
	revealedNonMineCells := 0

	for _, row := range b.Cells {
		for _, cell := range row {
			if !cell.IsMine {
				nonMineCells++

				if cell.IsRevealed {
					revealedNonMineCells++
				}
			}
		}
	}

	percentage := float64(revealedNonMineCells) / float64(nonMineCells)

	return percentage
}

func (b *Board) ToggleFlag(row, col int) {
	if row >= 0 && row < b.Rows && col >= 0 && col < b.Cols {
		b.Cells[row][col].IsFlagged = !b.Cells[row][col].IsFlagged
	}
}

type DisplayOptions struct {
	StartIndex *int
}

func (b *Board) Display(showMines bool, options *DisplayOptions) {
	// Add padding to the left for the column numbers
	fmt.Print("   ")

	startIndex := dStartIndex
	if options.StartIndex != nil {
		startIndex = *options.StartIndex
	}

	for c := 0; c < b.Cols; c++ {
		// TODO: Replace with user start position
		fmt.Printf("\x1b[34m%2d\x1b[0m ", c+startIndex) // Blue color for column numbers
	}

	// Print new line after the top row
	fmt.Println()

	for r := 0; r < b.Rows; r++ {
		// TODO: Replace with user start position
		fmt.Printf("\x1b[34m%2d\x1b[0m| ", r+startIndex) // Blue color for row numbers

		for c := 0; c < b.Cols; c++ {
			cell := b.Cells[r][c]

			if cell.IsRevealed {
				if cell.IsMine {
					fmt.Printf("\x1b[41m%s\x1b[0m%s", sMine, sSeperator) // Bg Red
					// fmt.Print("X  ")
				} else {
					switch cell.MinesAround {
					case 1:
						fmt.Printf("\x1b[94m%d\x1b[0m  ", cell.MinesAround) // Light red
					case 2:
						fmt.Printf("\x1b[32m%d\x1b[0m  ", cell.MinesAround) // Yellow
					case 3:
						fmt.Printf("\x1b[31m%d\x1b[0m  ", cell.MinesAround) // Bright yellow
					case 4:
						fmt.Printf("\x1b[34m%d\x1b[0m  ", cell.MinesAround) // Light green
					case 5:
						fmt.Printf("\x1b[33m%d\x1b[0m  ", cell.MinesAround) // Green
					case 6:
						fmt.Printf("\x1b[36m%d\x1b[0m  ", cell.MinesAround) // Light cyan
					case 7:
						fmt.Printf("\x1b[30m%d\x1b[0m  ", cell.MinesAround) // Cyan
					default:
						fmt.Printf("\x1b[90m%d\x1b[0m  ", cell.MinesAround) // Grey
					}
				}
			} else {
				if showMines && cell.IsMine {
					fmt.Printf("\x1b[41m%s\x1b[0m%s", sMine, sSeperator) // Bg Red
					// fmt.Print("X  ")
				} else if cell.IsFlagged {
					fmt.Printf("\x1b[91m%s\x1b[0m%s", sFlag, sSeperator)
					// fmt.Print("F  ")
				} else {
					fmt.Printf("\x1b[37m%s\x1b[0m%s", sHidden, sSeperator)
					// fmt.Print("•  ")
				}
			}
		}
		fmt.Println()
	}
}

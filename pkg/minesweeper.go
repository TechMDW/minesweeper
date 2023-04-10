package minesweeper

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/TechMDW/minesweeper/internal/util"
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

	BoardOptions   *BoardOptions
	DisplayOptions *DisplayOptions
}

type BoardOptions struct {
	Seed int64
}

type DisplayOptions struct {
	StartIndex  *int
	ANSI        *bool
	TopIndex    *bool
	BottomIndex *bool
	RightIndex  *bool
	LeftIndex   *bool

	// Symbols used to display the board
	SymbolMine      *string
	SymbolFlag      *string
	SymbolHidden    *string
	SymbolSeperator *string
}

// Symbols used to display the board
const (
	SymbolMine      = "X"
	SymbolFlag      = "F"
	SymbolHidden    = "•"
	SymbolSeperator = "  "
)

const (
	dStartIndex = 1
)

// NewBoard creates a new board with the given number of rows, columns, and mines.
//
// The board is initialized with all cells hidden and no mines placed.
func NewBoard(rows, cols, numMines int, boardOptions *BoardOptions, displayOptions *DisplayOptions) *Board {
	board := &Board{
		Rows:     rows,
		Cols:     cols,
		NumMines: numMines,
		Cells:    make([][]Cell, rows),

		Rand:           rand.New(rand.NewSource(boardOptions.Seed)),
		BoardOptions:   boardOptions,
		DisplayOptions: displayOptions,
	}

	// Check if the display options are nil and set them to their default values
	if board.DisplayOptions == nil {
		board.DisplayOptions = &DisplayOptions{}
	}

	if board.DisplayOptions.StartIndex == nil {
		board.DisplayOptions.StartIndex = util.IntPtr(1)
	}

	if board.DisplayOptions.ANSI == nil {
		board.DisplayOptions.ANSI = util.BoolPtr(true)
	}

	if board.DisplayOptions.LeftIndex == nil {
		board.DisplayOptions.LeftIndex = util.BoolPtr(true)
	}

	if board.DisplayOptions.RightIndex == nil {
		board.DisplayOptions.RightIndex = util.BoolPtr(false)
	}

	if board.DisplayOptions.SymbolMine == nil {
		board.DisplayOptions.SymbolMine = util.StringPtr(SymbolMine)
	}

	if board.DisplayOptions.SymbolFlag == nil {
		board.DisplayOptions.SymbolFlag = util.StringPtr(SymbolFlag)
	}

	if board.DisplayOptions.SymbolHidden == nil {
		board.DisplayOptions.SymbolHidden = util.StringPtr(SymbolHidden)
	}

	if board.DisplayOptions.SymbolSeperator == nil {
		board.DisplayOptions.SymbolSeperator = util.StringPtr(SymbolSeperator)
	}

	if board.DisplayOptions.TopIndex == nil {
		board.DisplayOptions.TopIndex = util.BoolPtr(true)
	}

	if board.DisplayOptions.BottomIndex == nil {
		board.DisplayOptions.BottomIndex = util.BoolPtr(false)
	}

	// Check board options are nil and set them to their default values
	if board.BoardOptions == nil {
		board.BoardOptions = &BoardOptions{}
	}

	if board.BoardOptions.Seed == 0 {
		board.BoardOptions.Seed = time.Now().UnixNano()
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

func (b *Board) Display(showMines bool) {
	// Add padding to the left for the column numbers
	fmt.Print("   ")

	startIndex := dStartIndex
	if b.DisplayOptions.StartIndex != nil {
		startIndex = *b.DisplayOptions.StartIndex
	}

	symbolMine := *b.DisplayOptions.SymbolMine
	symbolFlag := *b.DisplayOptions.SymbolFlag
	symbolHidden := *b.DisplayOptions.SymbolHidden
	symbolSeperator := *b.DisplayOptions.SymbolSeperator

	if *b.DisplayOptions.TopIndex {
		for c := 0; c < b.Cols; c++ {
			b.Printf("\x1b[34m%2d\x1b[0m ", c+startIndex)
		}
	}

	// Print new line after the top row
	fmt.Println()

	for r := 0; r < b.Rows; r++ {
		if *b.DisplayOptions.LeftIndex {
			b.Printf("\x1b[34m%2d\x1b[0m| ", r+startIndex)
		}

		for c := 0; c < b.Cols; c++ {
			cell := b.Cells[r][c]

			seperator := symbolSeperator

			if c == b.Cols-1 {
				seperator = ""
			}

			if cell.IsRevealed {
				if cell.IsMine {
					b.Printf("\x1b[41m%s\x1b[0m%s", symbolMine, seperator)
				} else {
					switch cell.MinesAround {
					case 1:
						b.Printf("\x1b[94m%d\x1b[0m%s", cell.MinesAround, seperator)
					case 2:
						b.Printf("\x1b[32m%d\x1b[0m%s", cell.MinesAround, seperator)
					case 3:
						b.Printf("\x1b[31m%d\x1b[0m%s", cell.MinesAround, seperator)
					case 4:
						b.Printf("\x1b[34m%d\x1b[0m%s", cell.MinesAround, seperator)
					case 5:
						b.Printf("\x1b[33m%d\x1b[0m%s", cell.MinesAround, seperator)
					case 6:
						b.Printf("\x1b[36m%d\x1b[0m%s", cell.MinesAround, seperator)
					case 7:
						b.Printf("\x1b[30m%d\x1b[0m%s", cell.MinesAround, seperator)
					default:
						b.Printf("\x1b[90m%d\x1b[0m%s", cell.MinesAround, seperator)
					}
				}
			} else {
				if showMines && cell.IsMine {
					b.Printf("\x1b[41m%s\x1b[0m%s", symbolMine, seperator)
				} else if cell.IsFlagged {
					b.Printf("\x1b[91m%s\x1b[0m%s", symbolFlag, seperator)
				} else {
					b.Printf("\x1b[37m%s\x1b[0m%s", symbolHidden, seperator)
				}
			}
		}

		if *b.DisplayOptions.RightIndex {
			b.Printf(" |\x1b[34m%d\x1b[0m", r+startIndex)
		}

		fmt.Println()
	}

	if *b.DisplayOptions.BottomIndex {
		fmt.Print("   ")

		for c := 0; c < b.Cols; c++ {
			b.Printf("\x1b[34m%2d\x1b[0m ", c+startIndex)
		}

		fmt.Println()
	}
}

func (b *Board) Printf(format string, a ...any) {
	if !*b.DisplayOptions.ANSI {
		format = RemoveAnsiEscapeCodes(format)

		for i, v := range a {
			t := reflect.TypeOf(v)

			if t.Kind() == reflect.String {
				a[i] = RemoveAnsiEscapeCodes(v.(string))
			} else {
				a[i] = v
			}
		}
	}

	fmt.Printf(format, a...)
}

func (b *Board) Print(a ...any) {
	if !*b.DisplayOptions.ANSI {
		for i, v := range a {
			t := reflect.TypeOf(v)

			if t.Kind() == reflect.String {
				a[i] = RemoveAnsiEscapeCodes(v.(string))
			} else {
				a[i] = v
			}
		}
	}

	fmt.Print(a...)
}

func (b *Board) Println(a ...any) {
	if !*b.DisplayOptions.ANSI {
		for i, v := range a {
			t := reflect.TypeOf(v)

			if t.Kind() == reflect.String {
				a[i] = RemoveAnsiEscapeCodes(v.(string))
			} else {
				a[i] = v
			}
		}
	}

	fmt.Println(a...)
}

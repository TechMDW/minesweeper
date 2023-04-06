package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
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
}

// privateRand is a private random number generator
var privateRand = rand.New(rand.NewSource(time.Now().UnixNano()))

var (
	// Clear the display
	dClear = "\033[H\033[2J"
)

// Symbols used to display the board
var (
	sMine      = "X"
	sFlag      = "F"
	sHidden    = "•"
	sSeperator = "  "
)

// NewBoard creates a new board with the given number of rows, columns, and mines.
//
// The board is initialized with all cells hidden and no mines placed.
func NewBoard(rows, cols, numMines int) *Board {
	board := &Board{
		Rows:     rows,
		Cols:     cols,
		NumMines: numMines,
		Cells:    make([][]Cell, rows),
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
			row := privateRand.Intn(b.Rows)
			col := privateRand.Intn(b.Cols)

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

	for c := 0; c < b.Cols; c++ {
		// TODO: Replace with user start position
		fmt.Printf("\x1b[34m%2d\x1b[0m ", c+1) // Blue color for column numbers
	}

	// Print new line after the top row
	fmt.Println()

	for r := 0; r < b.Rows; r++ {
		// TODO: Replace with user start position
		fmt.Printf("\x1b[34m%2d\x1b[0m| ", r+1) // Blue color for row numbers

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

// printHelp prints the help message.
func printHelp() {
	fmt.Println("r <row> <col> = reveal cell at position (row, col)")

	fmt.Println()

	fmt.Println("f <row> <col> = flag cell at position (row, col)")

	fmt.Println()

	fmt.Println("header = hide header (show only board + footer)")

	fmt.Println()

	fmt.Println("footer = hide footer (show only board + header)")

	fmt.Println()

	fmt.Println("h/help/imlost = show this help")

	fmt.Println()

	fmt.Println("q/quit/exit = quit game")

	fmt.Println()

	fmt.Println("Just press ENTER(↵) to continue ...")
}

func formatDuration(d time.Duration) string {
	nanoseconds := int(d.Nanoseconds()) % 1000
	microseconds := int(d.Microseconds()) % 1000
	seconds := int(d.Seconds()) % 60
	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours())

	format := ""
	if hours > 0 {
		format = fmt.Sprintf("%dh %dm %ds %dms %dns", hours, minutes, seconds, microseconds, nanoseconds)
	} else if minutes > 0 {
		format = fmt.Sprintf("%dm %ds %dms %dns", minutes, seconds, microseconds, nanoseconds)
	} else if seconds > 0 {
		format = fmt.Sprintf("%ds %dms %dns", seconds, microseconds, nanoseconds)
	} else if microseconds > 0 {
		format = fmt.Sprintf("%dms %dns", microseconds, nanoseconds)
	} else {
		format = fmt.Sprintf("%dns", nanoseconds)
	}

	return format
}

func formatPercentageBar(percentage float64, width int) string {
	filledWidth := int(percentage * float64(width))
	filled := strings.Repeat("=", filledWidth)

	unfilledWidth := width - filledWidth
	unfilled := strings.Repeat(" ", unfilledWidth)

	return fmt.Sprintf("[%s%s] %.1f%%", filled, unfilled, percentage*100)
}

func main() {
	flags := flag.NewFlagSet("minesweeper", flag.ExitOnError)

	rows := *flags.Int("rows", 10, "Number of rows")
	cols := *flags.Int("cols", 10, "Number of columns")
	numMines := *flags.Int("mines", 10, "Number of mines")
	showHelp := flags.Bool("help", false, "Show help")

	flags.Parse(os.Args[1:])

	if *showHelp {
		fmt.Println("Usage: minesweeper [OPTIONS]")
		fmt.Println("Options:")
		flags.PrintDefaults()
		return
	}

	board := NewBoard(rows, cols, numMines)

	gameOver := false
	footer := true
	header := true
	inHelp := false

	startTime := time.Now()

	for !gameOver {
		if inHelp {
			fmt.Scanln()

			inHelp = false
		}

		fmt.Println(dClear)

		if header {
			fmt.Println("Cells left: ", board.CellsNonRevealed())
			fmt.Println("Flags: ", board.FlagsCount())
			fmt.Println("Mines: ", board.NumMines)
			fmt.Println(" ", formatPercentageBar(board.RevealedPercentage(), cols*3-2))
		}

		board.Display(false)

		if footer {
			fmt.Println("Enter command: (r <row> <col> = reveal, f <row> <col> = flag, h/help/imlost = help)")
		} else {
			fmt.Println()
		}

		// Quick and dirty way to read input.
		var command string
		var row, col int
		fmt.Scanf("%s %d %d", &command, &row, &col)

		// Convert to lowercase
		command = strings.ToLower(command)

		// TODO: Replace with user start position
		if !reflect.ValueOf(row).IsZero() && !reflect.ValueOf(col).IsZero() {
			row--
			col--
		}

		switch command {
		case "r":
			if board.Reveal(row, col) {
				gameOver = true
			}
		case "f":
			board.ToggleFlag(row, col)
		case "h", "help", "imlost":
			fmt.Println(dClear)
			printHelp()
			inHelp = true
		case "footer":
			footer = !footer
		case "header":
			header = !header
		case "cheat":
			board.RevealAll()
			gameOver = true
		case "q", "quit", "exit":
			gameOver = true
		default:
			if reflect.ValueOf(command).IsZero() && reflect.ValueOf(row).IsZero() && reflect.ValueOf(col).IsZero() {
				continue
			}
			fmt.Println("Invalid command.")
		}
	}

	cellNonRevealed := board.CellsNonRevealed()
	cellsRevealed := board.CellsRevealed()
	flagCount := board.FlagsCount()
	percentage := board.RevealedPercentage()

	fmt.Println(dClear)

	board.Display(true)

	fmt.Println()

	if percentage == 1 {
		fmt.Printf("\x1b[32m%s\x1b[0m\n", "You won!") // Green color
	} else {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "You lost!") // Red color
	}

	fmt.Printf("You completed %d/%d cells in %s (%.2f%%)\n\n", cellsRevealed, cellsRevealed+cellNonRevealed, formatDuration(time.Since(startTime)), percentage*100)

	fmt.Printf("Size: %dX%d\n", rows, cols)
	fmt.Println("Amount of cells:", rows*cols)
	fmt.Println("Mines:", numMines)
	fmt.Println("Cells revealed:", cellsRevealed)
	fmt.Println("Cells left:", cellNonRevealed)
	fmt.Println("Flags:", flagCount)
}

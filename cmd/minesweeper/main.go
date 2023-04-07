package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/TechMDW/minesweeper/internal/util"
	minesweeper "github.com/TechMDW/minesweeper/pkg"
)

var (
	// Clear the display
	dClear = "\033[H\033[2J"
)

// printHelp prints the help message.
func printHelp() {
	fmt.Println("r <row> <col> = reveal cell at position (row, col)")

	fmt.Println()

	fmt.Println("c <col> <row> = reveal cell at position (col, row)")

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

func main() {
	flags := flag.NewFlagSet("minesweeper", flag.ExitOnError)

	rows := flags.Int("rows", 10, "Number of rows")
	cols := flags.Int("cols", 10, "Number of columns")
	mines := flags.Int("mines", 10, "Number of mines")
	startIndex := flags.Int("start", 1, "Start index (row and column start at this index)")
	clear := flags.Bool("clear", true, "Automatically clear the screen")

	showHelp := flags.Bool("help", false, "Show help")

	flags.Parse(os.Args[1:])

	if *showHelp {
		fmt.Println("Usage: minesweeper [OPTIONS]")
		fmt.Println("Options:")
		flags.PrintDefaults()
		return
	}

	if *startIndex < 0 {
		startIndex = util.IntPtr(0)
	}

	displayOptions := &minesweeper.DisplayOptions{
		StartIndex: startIndex,
	}

	board := minesweeper.NewBoard(*rows, *cols, *mines)

	gameOver := false
	inHelp := false
	footer := true
	header := true

	scanner := bufio.NewScanner(os.Stdin)

	startTime := time.Now()

	for !gameOver {
		if inHelp {
			fmt.Scanln()

			inHelp = false
		}

		if *clear {
			fmt.Println(dClear)
		}

		if header {
			fmt.Println("Cells left: ", board.CellsNonRevealed())
			fmt.Println("Flags: ", board.FlagsCount())
			fmt.Println("Mines: ", board.NumMines)
			fmt.Println(" ", util.FormatPercentageBar(board.RevealedPercentage(), *cols*3-2))
		}

		board.Display(false, displayOptions)

		if footer {
			fmt.Println("Enter command: (r <row> <col> = reveal, f <row> <col> = flag, h = help)")
		} else {
			fmt.Println()
		}

		// Read user input
		if !scanner.Scan() {
			fmt.Println("Error reading input.")

			continue
		}

		input := scanner.Text()

		var command string
		var row, col int

		// Parse user input
		_, err := fmt.Sscanf(input, "%s %d %d", &command, &row, &col)
		if err != nil && err != io.EOF {
			fmt.Println("Invalid input format.")
			continue
		}

		// Convert to lowercase
		command = strings.ToLower(command)

		if startIndex != nil {
			// Convert user input to 0-based index
			row -= *startIndex
			col -= *startIndex
		}

		switch command {
		case "r":
			if board.Reveal(row, col) {
				gameOver = true
			}
		case "c":
			if board.Reveal(col, row) {
				gameOver = true
			}
		case "f":
			board.ToggleFlag(row, col)
		case "h", "help", "imlost":
			if *clear {
				fmt.Println(dClear)
			}

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
			// Invalid command with red bg and white text
			fmt.Printf("\x1b[41;37m%s\x1b[0m\n", "Invalid command!")
		}
	}

	cellNonRevealed := board.CellsNonRevealed()
	cellsRevealed := board.CellsRevealed()
	flagCount := board.FlagsCount()
	percentage := board.RevealedPercentage()

	if *clear {
		fmt.Println(dClear)
	}

	board.Display(true, displayOptions)

	fmt.Println()

	if percentage == 1 {
		fmt.Printf("\x1b[32m%s\x1b[0m\n", "You won!") // Green color
	} else {
		fmt.Printf("\x1b[31m%s\x1b[0m\n", "You lost!") // Red color
	}

	fmt.Printf("You completed %d/%d cells in %s (%.2f%%)\n\n", cellsRevealed, cellsRevealed+cellNonRevealed, util.FormatDuration(time.Since(startTime)), percentage*100)

	fmt.Printf("Size: %d X %d\n", *rows, *cols)
	fmt.Println("Amount of cells:", *rows**cols)
	fmt.Println("Mines:", *mines)
	fmt.Println("Cells revealed:", cellsRevealed)
	fmt.Println("Cells left:", cellNonRevealed)
	fmt.Println("Flags:", flagCount)
}

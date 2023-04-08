package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
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

	// Game/Board options
	rows := flags.Int("rows", 10, "Number of rows")
	cols := flags.Int("cols", 10, "Number of columns")
	mines := flags.Int("mines", 10, "Number of mines")
	seed := flags.Int64("seed", time.Now().UnixNano(), "Seed for random number generator")

	// Display options
	startIndex := flags.Int("start", 1, "Start index (row and column start at this index)")
	ansi := flags.Bool("ansi", true, "Use ANSI escape codes to color the board")

	// Default/Debug options
	showHelp := flags.Bool("help", false, "Show help")
	clear := flags.Bool("clear", true, "Automatically clear the screen")

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
		ANSI:       ansi,
	}

	boardOptions := &minesweeper.BoardOptions{
		Seed: *seed,
	}

	board := minesweeper.NewBoard(*rows, *cols, *mines, boardOptions, displayOptions)

	gameOver := false
	inHelp := false
	manualQuit := false
	footer := true
	header := true

	scanner := bufio.NewScanner(os.Stdin)

	startTime := time.Now()

	for !gameOver {
		if inHelp {
			fmt.Scanln()

			inHelp = false
		}

		log.Println("clear:", *clear)
		if *clear {
			fmt.Println(dClear)
		}

		percentageDone := board.RevealedPercentage()

		if percentageDone == 1 {
			break
		}

		if header {
			fmt.Println("Cells left: ", board.CellsNonRevealed())
			fmt.Println("Flags: ", board.FlagsCount())
			fmt.Println("Mines: ", board.NumMines)
			fmt.Println(" ", util.FormatPercentageBar(percentageDone, *cols*3-2))
		}

		board.Display(false)

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
		var sIndex int

		// Parse user input
		_, err := fmt.Sscanf(input, "%s %d %d", &command, &row, &col)
		if err != nil && err != io.EOF {
			fmt.Println("Invalid input format.")
			continue
		}

		command = strings.ToLower(command)

		// Quick workaround for now
		if command == "start" {
			sIndex = row
		}

		if startIndex != nil {
			// Convert to 0-based index
			row -= *board.DisplayOptions.StartIndex
			col -= *board.DisplayOptions.StartIndex
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
		case "f", "fr", "rf":
			board.ToggleFlag(row, col)
		case "fc", "cf":
			board.ToggleFlag(col, row)
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
		case "ansi":
			board.DisplayOptions.ANSI = util.BoolPtr(!*board.DisplayOptions.ANSI)
		case "start":
			board.DisplayOptions.StartIndex = util.IntPtr(sIndex)
		case "cheat":
			board.RevealAll()
			gameOver = true
		case "restart":
			board = minesweeper.NewBoard(*rows, *cols, *mines, boardOptions, displayOptions)

			startTime = time.Now()
		case "q", "quit", "exit":
			gameOver = true
			manualQuit = true
		default:
			board.Printf("\x1b[41;37m%s\x1b[0m\n", "Invalid command!")
		}
	}

	gameDuration := time.Since(startTime)
	cellNonRevealed := board.CellsNonRevealed()
	cellsRevealed := board.CellsRevealed()
	flagCount := board.FlagsCount()
	percentage := board.RevealedPercentage()

	if *clear {
		fmt.Println(dClear)
	}

	board.Display(true)

	fmt.Println()

	if percentage == 1 {
		board.Printf("\x1b[32m%s\x1b[0m\n", "You won!")
	} else {
		board.Printf("\x1b[31m%s\x1b[0m\n", "You lost!")
	}

	fmt.Printf("You completed %d/%d cells in %s (%.2f%%)\n", cellsRevealed, cellsRevealed+cellNonRevealed, util.FormatDuration(gameDuration), percentage*100)
	fmt.Println("Seed:", *seed)
	fmt.Println("")
	fmt.Printf("Size: %d X %d\n", *rows, *cols)
	fmt.Println("Amount of cells:", *rows**cols)
	fmt.Println("Mines:", *mines)
	fmt.Println("Cells revealed:", cellsRevealed)
	fmt.Println("Cells left:", cellNonRevealed)
	fmt.Println("Flags:", flagCount)

	if manualQuit {
		return
	}

	// ALlow user to restart or quit
	fmt.Println("Enter command: (r = retry same seed, q = quit)")

	if !scanner.Scan() {
		fmt.Println("Error reading input.")

		return
	}

	input := scanner.Text()

	command := strings.ToLower(input)

	switch command {
	case "r", "restart":
		main()
	case "q", "quit", "exit":
		return
	default:
		fmt.Println("BYE!")
	}
}

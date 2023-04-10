package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TechMDW/minesweeper/internal/util"
	minesweeper "github.com/TechMDW/minesweeper/pkg"
)

const (
	// Clear the display
	dClear = "\033[H\033[2J"
)

type Config struct {
	rows            int
	cols            int
	mines           int
	footer          bool
	header          bool
	seed            int64
	startIndex      int
	ansi            bool
	showHelp        bool
	clear           bool
	topIndex        bool
	bottomIndex     bool
	rightIndex      bool
	leftIndex       bool
	symbolMine      string
	symbolFlag      string
	symbolHidden    string
	symbolSeperator string
}

type Command struct {
	Action string
	Args   []string
}

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

	fmt.Scanln()
}

func printFooter(config *Config) {
	if !config.footer {
		fmt.Println("")
		return
	}

	fmt.Println("Enter command: (r <row> <col> = reveal, f <row> <col> = flag, h = help)")
}

func printHeader(board *minesweeper.Board, config *Config) {
	if !config.header {
		return
	}

	percentageDone := board.RevealedPercentage()

	if config.header {
		fmt.Println("Cells left: ", board.CellsNonRevealed())
		fmt.Println("Flags: ", board.FlagsCount())
		fmt.Println("Mines: ", board.NumMines)
		fmt.Println(" ", util.FormatPercentageBar(percentageDone, config.cols*3-2))
	}
}

func parseFlags() *Config {
	flags := flag.NewFlagSet("minesweeper", flag.ExitOnError)

	// Game/Board options
	rows := flags.Int("rows", 10, "Number of rows")
	cols := flags.Int("cols", 10, "Number of columns")
	mines := flags.Int("mines", 10, "Number of mines")
	seed := flags.Int64("seed", time.Now().UnixNano(), "Seed for random number generator")
	header := flags.Bool("header", true, "Show header")
	footer := flags.Bool("footer", true, "Show footer")

	// Display options
	startIndex := flags.Int("start", 1, "Start index (row and column start at this index)")
	ansi := flags.Bool("ansi", true, "Use ANSI escape codes to color the board")
	symbolMine := flags.String("symbolMine", minesweeper.SymbolMine, "Symbol to use for mines")
	symbolFlag := flags.String("symbolFlag", minesweeper.SymbolFlag, "Symbol to use for flags")
	symbolHidden := flags.String("symbolHidden", minesweeper.SymbolHidden, "Symbol to use for hidden cells")
	symbolSeperator := flags.String("symbolSeperator", minesweeper.SymbolSeperator, "Symbol to use for seperating cells")
	topIndex := flags.Bool("topIndex", true, "Show top index")
	bottomIndex := flags.Bool("bottomIndex", false, "Show bottom index")
	rightIndex := flags.Bool("rightIndex", false, "Show right index")
	leftIndex := flags.Bool("leftIndex", true, "Show left index")

	// Default/Debug options
	showHelp := flags.Bool("help", false, "Show help")
	clear := flags.Bool("clear", true, "Automatically clear the screen")

	flags.Parse(os.Args[1:])

	if *startIndex < 0 {
		startIndex = util.IntPtr(0)
	}

	return &Config{
		rows:            *rows,
		cols:            *cols,
		mines:           *mines,
		seed:            *seed,
		startIndex:      *startIndex,
		ansi:            *ansi,
		showHelp:        *showHelp,
		clear:           *clear,
		header:          *header,
		footer:          *footer,
		topIndex:        *topIndex,
		bottomIndex:     *bottomIndex,
		rightIndex:      *rightIndex,
		leftIndex:       *leftIndex,
		symbolMine:      *symbolMine,
		symbolFlag:      *symbolFlag,
		symbolHidden:    *symbolHidden,
		symbolSeperator: *symbolSeperator,
	}
}

func printStatistics(board *minesweeper.Board, startTime time.Time, config *Config, manualQuit bool) {
	gameDuration := time.Since(startTime)
	cellNonRevealed := board.CellsNonRevealed()
	cellsRevealed := board.CellsRevealed()
	flagCount := board.FlagsCount()
	percentage := board.RevealedPercentage()

	if config.clear {
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
	fmt.Println("Seed:", config.seed)
	fmt.Println("")
	fmt.Printf("Size: %d X %d\n", config.rows, config.cols)
	fmt.Println("Amount of cells:", config.rows*config.cols)
	fmt.Println("Mines:", config.mines)
	fmt.Println("Cells revealed:", cellsRevealed)
	fmt.Println("Cells left:", cellNonRevealed)
	fmt.Println("Flags:", flagCount)

	if manualQuit {
		return
	}

	// Allow user to restart or quit
	fmt.Println("Enter command: (r = retry same seed, rn = retry new seed, q = quit)")

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		fmt.Println("Error reading input.")
		return
	}

	input := scanner.Text()
	command := strings.ToLower(input)

	switch command {
	case "r", "restart":
		playGame(config)
	case "rn", "restartnew":
		config.seed = time.Now().UnixNano()
		playGame(config)
	case "q", "quit", "exit":
		return
	default:
		fmt.Println("BYE!")
	}
}

func parseInput(input string) (*Command, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	return &Command{
		Action: strings.ToLower(parts[0]),
		Args:   parts[1:],
	}, nil
}

func handleInput(input string, board *minesweeper.Board, config *Config) (gameOver, manualQuit bool) {
	command, err := parseInput(input)
	if err != nil {
		fmt.Println("Invalid input format.")
		return
	}

	switch command.Action {
	case "r":
		if handleReveal(command.Args, false, board, config) {
			gameOver = true
		}
	case "c":
		if handleReveal(command.Args, true, board, config) {
			gameOver = true
		}
	case "f", "fr", "rf":
		handleFlag(command.Args, false, board, config)
	case "fc", "cf":
		handleFlag(command.Args, true, board, config)
	case "h", "help", "imlost":
		if config.clear {
			fmt.Println(dClear)
		}

		printHelp()
	case "footer":
		config.footer = !config.footer
	case "header":
		config.header = !config.header
	case "ansi":
		board.DisplayOptions.ANSI = util.BoolPtr(!*board.DisplayOptions.ANSI)
	case "start":
		var sIndex int
		if _, err := fmt.Sscanf(strings.Join(command.Args, " "), "start %d", &sIndex); err != nil && err != io.EOF {
			fmt.Println("Invalid input format.")
			return
		}

		board.DisplayOptions.StartIndex = util.IntPtr(sIndex)
	case "cheat":
		board.RevealAll()
		gameOver = true
	case "q", "quit", "exit":
		gameOver = true
		manualQuit = true
	default:
		board.Printf("\x1b[41;37m%s\x1b[0m\n", "Invalid command!")
	}

	return
}

func handleReveal(args []string, inverted bool, board *minesweeper.Board, config *Config) (gameOver bool) {
	if len(args) < 2 {
		fmt.Println("Invalid input format")
		return
	}

	x, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid input format")
		return
	}

	y := make([]int, 0, len(args)-1)
	for _, arg := range args[1:] {
		num, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Invalid input format")
			return
		}
		y = append(y, num)
	}

	// Convert to 0-based index
	x -= *board.DisplayOptions.StartIndex

	for _, yi := range y {
		yi -= *board.DisplayOptions.StartIndex

		// Determine the correct order of arguments for the board.Reveal function
		row, col := x, yi
		if inverted {
			row, col = col, row
		}

		if board.Reveal(row, col) {
			gameOver = true
			break
		}
	}

	return
}

func handleFlag(args []string, inverted bool, board *minesweeper.Board, config *Config) {
	if len(args) < 2 {
		fmt.Println("Invalid input format")
		return
	}

	x, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Invalid input format")
		return
	}

	y := make([]int, 0, len(args)-1)
	for _, arg := range args[1:] {
		num, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Println("Invalid input format")
			return
		}
		y = append(y, num)
	}

	// Convert to 0-based index
	x -= *board.DisplayOptions.StartIndex

	for _, yi := range y {
		yi -= *board.DisplayOptions.StartIndex

		// Determine the correct order of arguments for the board.Reveal function
		row, col := x, yi
		if inverted {
			row, col = col, row
		}

		board.ToggleFlag(row, col)
	}
}

func playGame(config *Config) {
	boardOptions := &minesweeper.BoardOptions{
		Seed: config.seed,
	}

	displayOptions := &minesweeper.DisplayOptions{
		StartIndex: &config.startIndex,
		ANSI:       &config.ansi,

		SymbolMine:      &config.symbolMine,
		SymbolFlag:      &config.symbolFlag,
		SymbolHidden:    &config.symbolHidden,
		SymbolSeperator: &config.symbolSeperator,

		TopIndex:    &config.topIndex,
		LeftIndex:   &config.leftIndex,
		RightIndex:  &config.rightIndex,
		BottomIndex: &config.bottomIndex,
	}

	board := minesweeper.NewBoard(config.rows, config.cols, config.mines, boardOptions, displayOptions)

	gameOver := false
	manualQuit := false

	scanner := bufio.NewScanner(os.Stdin)

	startTime := time.Now()

	for !gameOver {
		if config.clear {
			fmt.Println(dClear)
		}

		percentageDone := board.RevealedPercentage()

		if percentageDone == 1 {
			break
		}

		printHeader(board, config)

		board.Display(false)

		printFooter(config)

		// Read user input
		if !scanner.Scan() {
			fmt.Println("Error reading input.")
			continue
		}

		input := scanner.Text()
		gameOver, manualQuit = handleInput(input, board, config)
	}

	printStatistics(board, startTime, config, manualQuit)
}

func main() {
	config := parseFlags()

	if config.showHelp {
		fmt.Println("Usage: minesweeper [OPTIONS]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		return
	}

	playGame(config)
}

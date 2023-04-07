package minesweeper_test

import (
	"math"
	"testing"

	minesweeper "github.com/TechMDW/minesweeper/pkg"
)

func TestNewBoard(t *testing.T) {
	rows, cols, numMines := 10, 10, 10
	boardOptions := &minesweeper.BoardOptions{Seed: 5}
	displayOptions := &minesweeper.DisplayOptions{StartIndex: nil, ANSI: nil}

	board := minesweeper.NewBoard(rows, cols, numMines, boardOptions, displayOptions)

	if board.Rows != rows {
		t.Errorf("Expected Rows to be %d, but got %d", rows, board.Rows)
	}

	if board.Cols != cols {
		t.Errorf("Expected Cols to be %d, but got %d", cols, board.Cols)
	}

	if board.NumMines != numMines {
		t.Errorf("Expected NumMines to be %d, but got %d", numMines, board.NumMines)
	}
}

func TestPlaceMines(t *testing.T) {
	rows, cols, numMines := 10, 10, 10
	boardOptions := &minesweeper.BoardOptions{Seed: 5}
	displayOptions := &minesweeper.DisplayOptions{StartIndex: nil, ANSI: nil}

	board := minesweeper.NewBoard(rows, cols, numMines, boardOptions, displayOptions)

	mineCount := 0
	for _, row := range board.Cells {
		for _, cell := range row {
			if cell.IsMine {
				mineCount++
			}
		}
	}

	if mineCount != numMines {
		t.Errorf("Expected mine count to be %d, but got %d", numMines, mineCount)
	}
}

func TestReveal(t *testing.T) {
	rows, cols, numMines := 10, 10, 10
	boardOptions := &minesweeper.BoardOptions{Seed: 5}
	displayOptions := &minesweeper.DisplayOptions{StartIndex: nil, ANSI: nil}

	board := minesweeper.NewBoard(rows, cols, numMines, boardOptions, displayOptions)

	// Reveal cell without a mine
	if board.Reveal(0, 0) {
		t.Error("Expected Reveal to return false for a cell without a mine")
	}

	if board.Reveal(8, 5) {
		t.Error("Expected Reveal to return false for a cell without a mine")
	}

	// Reveal cell with a mine
	if !board.Reveal(9, 0) {
		t.Error("Expected Reveal to return true for a cell with a mine")
	}

	if !board.Reveal(4, 8) {
		t.Error("Expected Reveal to return true for a cell with a mine")
	}
}

func TestToggleFlag(t *testing.T) {
	rows, cols, numMines := 10, 10, 10
	boardOptions := &minesweeper.BoardOptions{Seed: 5}
	displayOptions := &minesweeper.DisplayOptions{StartIndex: nil, ANSI: nil}

	board := minesweeper.NewBoard(rows, cols, numMines, boardOptions, displayOptions)

	// Toggle flag on
	board.ToggleFlag(0, 0)
	if !board.Cells[0][0].IsFlagged {
		t.Error("Expected IsFlagged to be true after toggling flag on")
	}

	// Toggle flag off
	board.ToggleFlag(0, 0)
	if board.Cells[0][0].IsFlagged {
		t.Error("Expected IsFlagged to be false after toggling flag off")
	}
}

func TestFlagsCount(t *testing.T) {
	rows, cols, numMines := 10, 10, 10
	boardOptions := &minesweeper.BoardOptions{Seed: 5}
	displayOptions := &minesweeper.DisplayOptions{StartIndex: nil, ANSI: nil}

	board := minesweeper.NewBoard(rows, cols, numMines, boardOptions, displayOptions)

	// Place flags
	board.ToggleFlag(0, 0)
	board.ToggleFlag(1, 1)
	board.ToggleFlag(2, 2)

	flagsCount := board.FlagsCount()
	expectedFlagsCount := 3

	if flagsCount != expectedFlagsCount {
		t.Errorf("Expected flags count to be %d, but got %d", expectedFlagsCount, flagsCount)
	}
}

func TestRevealedPercentage(t *testing.T) {
	rows, cols, numMines := 10, 10, 10
	boardOptions := &minesweeper.BoardOptions{Seed: 5}
	displayOptions := &minesweeper.DisplayOptions{StartIndex: nil, ANSI: nil}

	board := minesweeper.NewBoard(rows, cols, numMines, boardOptions, displayOptions)

	board.Reveal(0, 0)
	board.Reveal(0, 1)
	board.Reveal(1, 0)

	revealedPercentage := board.RevealedPercentage()
	expectedRevealedPercentage := 0.777778

	if math.Round(revealedPercentage*10)/10 != math.Round(expectedRevealedPercentage*10)/10 {
		t.Errorf("Expected revealed percentage to be %f, but got %f", expectedRevealedPercentage, revealedPercentage)
	}
}

func TestRevealAll(t *testing.T) {
	rows, cols, numMines := 10, 10, 10
	boardOptions := &minesweeper.BoardOptions{Seed: 5}
	displayOptions := &minesweeper.DisplayOptions{StartIndex: nil, ANSI: nil}

	board := minesweeper.NewBoard(rows, cols, numMines, boardOptions, displayOptions)

	board.RevealAll()

	revealedCells := board.CellsRevealed()
	expectedRevealedCells := rows*cols - numMines

	if revealedCells != expectedRevealedCells {
		t.Errorf("Expected revealed cells to be %d, but got %d", expectedRevealedCells, revealedCells)
	}
}

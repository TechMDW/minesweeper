# Minesweeper CLI

This is a command-line version of the classic Minesweeper game written in Go. Please note that this project is a work in progress (WIP) and may have some bugs or incomplete features.

## How to Play

Before getting started I recommend checking out the start flags by running `<binary> -h`. This will allow you to change the "rules" of the game.

The game provides a simple command-line interface for playing Minesweeper. The board is displayed as a grid of cells, with row and column numbers as labels. The game supports the following commands:

- `r <row> <col>`: Reveal the cell at the specified row and column.
- `c <col> <row>`: Reveal the cell at the specified column and row
- `f <row> <col>`: Toggle a flag on or off at the specified row and column.
- `header`: Hide or show the header information.
- `footer`: Hide or show the footer information.
- `q`, `quit`, `exit`: Quit the game.
- `h`, `help`, `imlost`: Display the help message with a list of commands.

The game continues until all non-mine cells are revealed or a mine is revealed.

## Start flags

### Game options

- `-rows <int>`: Number of rows (default: 10)
- `-cols <int>`: Number of columns (default: 10)
- `-mines <int>`: Number of mines (default: 10)
- `-seed <int64>`: Seed for random number generator (default: current Unix timestamp in nanoseconds)

### Display options

- `-start <int>`: Start index (row and column start at this index, default: 1)
- `-ansi=<true|false>`: Use ANSI escape codes to color the board (default: true)

### Help

- `-h / -help`: Show help (default: false)

### Debug

- `-clear=<true|false>`: Automatically clear the screen (default: true)

### Example usage

1. `minesweeper -rows 10 -cols 20 -mines 30`
2. `minesweeper -h`
3. `minesweeper -ansi=false -clear=false -seed 50`
4. `minesweeper -rows 30 -ansi=false`

## Download prebuild package

1. Download the latest version of Minesweeper from the [GitHub releases page](https://github.com/TechMDW/minesweeper/releases/latest).
2. Open your terminal.
3. Navigate to the directory where they downloaded the file `cd <path>` (Linux/macOS/Windows).
4. Run the game binary with `./minesweeper` (Linux/macOS) or `minesweeper.exe` (Windows).

## Download via `go install`

[Go](https://golang.org/) is required to be installed when following this!

To install and run the game, you need to have the [Go](https://go.dev/dl/) installed on your system. Once Go is installed, follow these

1. Open a terminal on your computer.

2. Run the following command to download and install the Minesweeper CLI:
   `go install github.com/TechMDW/minesweeper/cmd/minesweeper@latest`
   This will download and install the latest version of the Minesweeper CLI on your system.

   OR

   `go install github.com/TechMDW/minesweeper/cmd/minesweeper@<version>`
   This will download a specific version if you replace the `<version>` with an available [releases](https://github.com/TechMDW/minesweeper/releases)

3. Once installation is complete, you can run the Minesweeper CLI from anywhere in your terminal by typing:
   `minesweeper`
   This should start the game and display the game board in your terminal.

## Building and Running

To build and run the game, you need to have the [Go](https://go.dev/dl/) installed on your system. Once Go is installed, follow these steps:

1. Clone the repository or download the source code.
2. Open a terminal and navigate to the project directory.
3. Run `go build` to build the game binary.
4. Navigate to the directory where the binary was built.
5. Run the game binary with `./minesweeper` (Linux/macOS) or `minesweeper.exe` (Windows).

You can also run the game directly without building the binary by running `go run main.go` from the project directory.

## Contributing

As this project is a work in progress, contributions are welcome! Feel free to report bugs, suggest improvements, or submit pull requests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

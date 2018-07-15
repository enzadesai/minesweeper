package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/enzadesai/minesweeper"
)

var (
	rows  = flag.Int("rows", 9, "number of rows")
	cols  = flag.Int("cols", 9, "number of columns")
	mines = flag.Int("mines", 10, "number of mines")
)

func main() {
	flag.Parse()
	b := minesweeper.New(*rows, *cols, *mines)
	for {
		redraw(b)
		row, col := askMove()
		if !b.Play(row, col) {
			break
		}
	}

	redraw(b)

	// TODO: Make messages random.
	if b.Lost() {
		fmt.Println("Woops! There's always next time.")
	} else {
		fmt.Println("Winner winner chicken dinner!")
	}
}

func draw(b *minesweeper.Board) {
	// Header
	fmt.Print("    ")
	for i := 0; i < *cols; i++ {
		fmt.Printf("%2d ", i)
	}
	fmt.Println()
	for i := -2; i < *cols; i++ {
		fmt.Print("---")
	}
	fmt.Println()

	// Rows
	for r := 0; r < *rows; r++ {
		row := b.Row(r)
		fmt.Printf("%2d | ", r)
		for c := 0; c < *cols; c++ {
			cell := row[c]
			if !cell.Open {
				fmt.Print(" * ")
				continue
			}

			switch cell.Value {
			case minesweeper.CellEmpty:
				fmt.Print("   ")
			case minesweeper.CellMine:
				fmt.Print(" 💣 ")
			default:
				fmt.Printf(" %v ", cell.Value)

			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func redraw(b *minesweeper.Board) {
	// TODO: Make it cross-platform
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("Minesweeper\n")
	draw(b)
}

func askMove() (r, c int) {
	fmt.Println()
	for {
		r, c, err := ask()
		if err == nil {
			return r, c
		}
		fmt.Printf("Bad input: %v\n", err)
	}
}

func ask() (r, c int, err error) {
	fmt.Print("Enter move (row col): ")

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		panic("Failed to get move")
	}

	text = strings.TrimRight(text, "\n")
	splits := strings.Split(text, " ")
	if len(splits) != 2 {
		return -1, -1, fmt.Errorf("need exactly two numbers")
	}
	row, err := strconv.Atoi(splits[0])
	if err != nil {
		return -1, -1, err
	}
	col, err := strconv.Atoi(splits[1])
	if err != nil {
		return -1, -1, err
	}
	return row, col, nil
}

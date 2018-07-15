package minesweeper

import (
	"math/rand"
	"time"
)

const (
	CellEmpty CellValue = 0
	CellMine  CellValue = -1
)

type CellValue int

// Cell represents individual cells on a Board.
type Cell struct {
	Value CellValue
	Open  bool
}

// Board is an instance of the Minesweeper game.
type Board struct {
	rows  int
	cols  int
	mines int

	state [][]Cell
	openc int

	lost bool
}

func (b *Board) isMine(r, c int) bool {
	return b.state[r][c].Value == CellMine
}

// Row returns the r-th row on the board.
func (b *Board) Row(r int) []Cell {
	return b.state[r]
}

func (b *Board) Lost() bool {
	return b.lost
}

func (b *Board) IsValidPos(r, c int) bool {
	if r <= -1 || r >= b.rows {
		return false
	}
	if c <= -1 || c >= b.cols {
		return false
	}
	return true
}

func (b *Board) Won() bool {
	if b.openc == b.rows*b.cols-b.mines {
		return true
	}
	return false
}

// Play plays the r, c cell. Returns if play can continue.
func (b *Board) Play(r, c int) bool {
	if b.lost {
		return false
	}

	if b.Won() {
		return false
	}

	if b.openc == b.rows*b.cols-b.mines {
		return false
	}

	if !b.IsValidPos(r, c) {
		return true
	}

	if b.isMine(r, c) {
		b.lost = true
		b.state[r][c].Open = true
		return false
	}

	if b.state[r][c].Open {
		return true
	}

	b.expand(r, c)
	return !b.Won()
}

func (b *Board) expand(r, c int) {
	if !b.IsValidPos(r, c) {
		return
	}

	cell := b.state[r][c]
	if cell.Value == CellMine {
		return
	}
	if cell.Open {
		return
	}

	b.state[r][c].Open = true
	b.openc++

	if cell.Value != CellEmpty {
		return
	}

	for row := r - 1; row <= r+1; row++ {
		for col := c - 1; col <= c+1; col++ {
			if row == r && col == c {
				continue
			}
			b.expand(row, col)
		}
	}
}

func (b *Board) fillMines() {
	var (
		seed = int64(time.Now().UnixNano())
		rng  = rand.New(rand.NewSource(seed))

		minec int
		total = b.rows * b.cols
	)
	for minec < b.mines {
		var (
			rc  = rng.Intn(total)
			row = rc / b.cols
			col = rc % b.cols
		)
		if b.isMine(row, col) {
			// TODO: Potentially infinite loop.
			continue
		}
		b.state[row][col] = Cell{Value: CellMine}
		minec++
	}
}

func (b *Board) initCells() {
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			if !b.isMine(row, col) {
				mines := CellValue(b.countMines(row, col))
				b.state[row][col].Value = mines
			}
		}
	}
}

func (b *Board) countMines(r, c int) int {
	var mines int
	for row := r - 1; row <= r+1; row++ {
		for col := c - 1; col <= c+1; col++ {
			if !b.IsValidPos(row, col) {
				continue
			}
			if b.isMine(row, col) {
				mines++
			}
		}
	}
	return mines
}

// New
func New(r, c, mines int) *Board {
	// Init state.
	state := make([][]Cell, r)
	for i := 0; i < r; i++ {
		state[i] = make([]Cell, c)
	}

	b := &Board{
		state: state,
		rows:  r,
		cols:  c,
		mines: mines,
	}

	b.fillMines()
	b.initCells()
	return b
}

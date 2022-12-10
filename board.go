package main

import "github.com/gdamore/tcell/v2"

type BoardProp struct {
	value int
	color int
}

type TetrisBoard struct {
	con [][]BoardProp
}

func (board *TetrisBoard) fillBlank() {
	board.con = make([][]BoardProp, 20)
	for ty := 0; ty < 20; ty++ {
		board.con[ty] = make([]BoardProp, 10)
		for tx := 0; tx < 10; tx++ {
			board.con[ty][tx] = BoardProp{0, 7}
		}
	}
}

func (board *TetrisBoard) bindTetromino(tetro Tetromino, x, y int) {
	for ty := 0; ty < tetro.height; ty++ {
		for tx := 0; tx < tetro.width; tx++ {
			if ty+y >= 0 && tx+x >= 0 {
				if tetro.con[ty][tx] > 0 {
					board.con[ty+y][tx+x].color = tetro.color
					board.con[ty+y][tx+x].value += 1
				}
			}
		}
	}
}

func (board *TetrisBoard) canBindTetromino(tetro Tetromino, x, y int) bool {
	for ty := 0; ty < tetro.height; ty++ {
		for tx := 0; tx < tetro.width; tx++ {
			if ty+y >= 0 && tx+x >= 0 {
				teCon := tetro.con[ty][tx]
				boCon := board.con[ty+y][tx+x]
				if teCon+boCon.value > 1 {
					return false
				}
			}
		}
	}
	return true
}

func (board TetrisBoard) print(scr tcell.Screen, x, y int) {
	for ty, con := range board.con {
		for tx, c := range con {
			if c.value > 0 {
				printText(scr, x+tx*2, y+ty, 2, 1, getStyleByInt(c.color), "██")
			}
		}
	}
}

func (board *TetrisBoard) updateBoard() int {
	var count int
	nextCon := make([][]BoardProp, 0)
	for _, con := range board.con {
		var sum int
		for _, c := range con {
			sum += c.value
		}
		if sum < 10 {
			nextCon = append(nextCon, con)
		} else {
			blankCon := make([][]BoardProp, 1)
			blankCon[0] = make([]BoardProp, 10)
			nextCon = append(blankCon, nextCon...)
			count++
		}
	}
	board.con = nextCon
	return count
}

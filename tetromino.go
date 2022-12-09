package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

type Tetromino struct {
	con    [][]int
	color  int
	width  int
	height int
	ori    int
}

func (tetro *Tetromino) setRandom() {
	switch rand.Intn(7) {
	case 0:
		tetro.con = [][]int{
			{1, 1, 1, 1},
		}
		tetro.color = 1
		tetro.width, tetro.height = 4, 1
	case 1:
		tetro.con = [][]int{
			{1, 0, 0},
			{1, 1, 1},
		}
		tetro.color = 1
		tetro.width, tetro.height = 3, 2
	case 2:
		tetro.con = [][]int{
			{0, 0, 1},
			{1, 1, 1},
		}
		tetro.color = 0
		tetro.width, tetro.height = 3, 2
	case 3:
		tetro.con = [][]int{
			{1, 1},
			{1, 1},
		}
		tetro.color = 3
		tetro.width, tetro.height = 2, 2
	case 4:
		tetro.con = [][]int{
			{0, 1, 1},
			{1, 1, 0},
		}
		tetro.color = 4
		tetro.width, tetro.height = 3, 2
	case 5:
		tetro.con = [][]int{
			{0, 1, 0},
			{1, 1, 1},
		}
		tetro.color = 5
		tetro.width, tetro.height = 3, 2
	case 6:
		tetro.con = [][]int{
			{1, 1, 0},
			{0, 1, 1},
		}
		tetro.color = 6
		tetro.width, tetro.height = 3, 2
	}
}

func (tetro *Tetromino) getOriXY() (int, int) {
	if tetro.color == 4 {
		return 0, 0
	} else if tetro.color == 2 {
		switch tetro.ori {
		case 1:
			return 2, 0
		case 3:
			return 2, 0
		}
	} else {
		switch tetro.ori {
		case 1:
			return 1, 0
		case 2:
			return 0, 1
		}
	}
	return 0, 0
}

func (tetro *Tetromino) print(scr tcell.Screen, x, y, minX, minY, maxX, maxY int, withOri bool) {
	oriX, oriY := 0, 0
	if withOri {
		oriX, oriY = tetro.getOriXY()
	}
	style := getStyleByInt(tetro.color)
	for ty, con := range tetro.con {
		posY := y + ty + oriY
		for tx, c := range con {
			posX := x + (tx+oriX)*2
			if c > 0 && posX >= minX && posX < maxX && posY >= minY && posY < maxY {
				printText(scr, posX, posY, 2, 1, style, "██")
			}
		}
	}
}

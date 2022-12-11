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
	tetro.ori = 0
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
		tetro.color = 2
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
	if tetro.color == 3 {
		return 0, 0
	} else if tetro.color == 1 {
		switch tetro.ori {
		case 0:
			return 0, 1
		case 1:
			return 2, 0
		case 2:
			return 0, 2
		case 3:
			return 1, 0
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

func (tetro *Tetromino) print(scr tcell.Screen, x, y, minY int, withOri bool, runet string) {
	oriX, oriY := 0, 0
	if withOri {
		oriX, oriY = tetro.getOriXY()
	}
	style := getStyleByInt(tetro.color)
	for ty, con := range tetro.con {
		posY := y + ty + oriY
		for tx, c := range con {
			posX := x + (tx+oriX)*2
			if c > 0 && posY >= minY {
				printText(scr, posX, posY, 2, 1, style, runet)
			}
		}
	}
}

func (tetro *Tetromino) clone() Tetromino {
	var clonedTetro Tetromino
	clonedTetro.con = make([][]int, tetro.height)
	for ty := 0; ty < tetro.height; ty++ {
		clonedTetro.con[ty] = make([]int, term_width)
		for tx := 0; tx < tetro.width; tx++ {
			clonedTetro.con[ty][tx] = tetro.con[ty][tx]
		}
	}
	clonedTetro.color = tetro.color
	clonedTetro.width = tetro.width
	clonedTetro.height = tetro.height
	clonedTetro.ori = tetro.ori
	return clonedTetro
}

func (tetro *Tetromino) rotateClockWise() {
	nextWidth, nextHeight := tetro.height, tetro.width
	nextCon := make([][]int, nextHeight)
	nextOri := tetro.ori + 1
	if nextOri > 3 {
		nextOri = 0
	}
	for ty := 0; ty < nextHeight; ty++ {
		nextCon[ty] = make([]int, nextWidth)
		for tx := 0; tx < nextWidth; tx++ {
			nextCon[ty][tx] = tetro.con[tx][ty]
		}
	}
	for _, con := range nextCon {
		for i, j := 0, nextWidth-1; i < j; i, j = i+1, j-1 {
			con[i], con[j] = con[j], con[i]
		}
	}
	tetro.width = nextWidth
	tetro.height = nextHeight
	tetro.con = nextCon
	tetro.ori = nextOri

}

func (tetro *Tetromino) rotateCounterClockWise() {
	nextWidth, nextHeight := tetro.height, tetro.width
	nextCon := make([][]int, nextHeight)
	nextOri := tetro.ori - 1
	if nextOri < 0 {
		nextOri = 3
	}
	for ty := 0; ty < nextHeight; ty++ {
		nextCon[ty] = make([]int, nextWidth)
		for tx := 0; tx < nextWidth; tx++ {
			nextCon[ty][tx] = tetro.con[tx][ty]
		}
	}
	for i, j := 0, nextHeight-1; i < j; i, j = i+1, j-1 {
		nextCon[i], nextCon[j] = nextCon[j], nextCon[i]
	}
	tetro.width = nextWidth
	tetro.height = nextHeight
	tetro.con = nextCon
	tetro.ori = nextOri

}

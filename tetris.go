package main

import (
	"time"

	"github.com/gdamore/tcell/v2"
)

// 45x28
func drawTetrisScreen(scr tcell.Screen, force bool) {
	scr.Fill(0, DEF_SF)
	if term_width < 45 || term_height < 28 {
		is_term_too_small = true
		printYourTermIsTooSmall(scr, 45, 28)
	} else {
		is_term_too_small = false
		centerX, centerY := (term_width-45)/2, (term_height-28)/2
		printHeader(scr, centerX, centerY)
		printTetrisBox(scr)
	}
	if force {
		scr.Sync()
	} else {
		scr.Show()
	}
}

func printTetrisBox(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-28)/2
	printBox(scr, x, y+6, 22, 22, DEF_SF, true)
	cur_board.print(scr, x+1, y+6+1)
	cur_tetro.print(scr, x+1+cur_X, y+6+1+predict_Y, y+6+1, true, "◤◢")
	cur_tetro.print(scr, x+1+cur_X, y+6+1+cur_Y, y+6+1, true, "██")
}

func dropTetromino(x int) {
	cur_tetro = Tetromino{}
	cur_tetro.setRandom()
	if x == -1 {
		x = 4
	}
	if statMove := isMoveAvailable(cur_tetro, cur_X, x, cur_Y); statMove != 0 {
		switch statMove {
		case 1:
			x = 0
		case 2:
			x = 20 - cur_tetro.width*2
		}
	}
	cur_X, cur_Y = x, -2
	createPrediction()
}

func createPrediction() {
	predict_Y = cur_Y
	for isMoveAvailable(cur_tetro, cur_X, cur_X, predict_Y+1) == 0 {
		predict_Y++
	}
}

var sess int

func updateTetris(scr tcell.Screen) {
	if !is_elemination_session {
		printTetrisBox(scr)
		scr.Show()
		isCantGoDown := isMoveAvailable(cur_tetro, cur_X, cur_X, cur_Y+1) != 0
		if isCantGoDown {
			is_elemination_session = true
			select {
			case <-time.After(time.Second / 10):
				sess++
				break
			case <-updated:
				break
			}
			if isMoveAvailable(cur_tetro, cur_X, cur_X, cur_Y+1) != 0 {
				oriX, oriY := cur_tetro.getOriXY()
				cur_board.bindTetromino(cur_tetro, oriX+cur_X/2, oriY+cur_Y)
				dropTetromino(cur_X)
				is_end = predict_Y == cur_Y
			}
			printTetrisBox(scr)
			scr.Show()
			is_elemination_session = false
		}
	} else {
		updated <- true
	}
}

// right = true to move right, right = false to move left
func moveLeftRight(right bool) bool {
	next_X := cur_X - 2
	if right {
		next_X = cur_X + 2
	}

	if isMoveAvailable(cur_tetro, cur_X, next_X, cur_Y) == 0 {
		cur_X = next_X
		createPrediction()
		return true
	}
	return false
}

func moveDown() bool {
	next_Y := cur_Y + 1
	if isMoveAvailable(cur_tetro, cur_X, cur_X, next_Y) == 0 {
		cur_Y = next_Y
		return true
	}
	return false
}

func rotateClockWise() bool {
	var success bool
	nextTetro := cur_tetro.clone()
	nextTetro.rotateClockWise()
	nextX, nextY := cur_X, cur_Y
	statMove := isMoveAvailable(nextTetro, cur_X, cur_X, cur_Y)
	if statMove == 3 {
		searchAvailableYMove(nextTetro, cur_X, &nextX, &nextY)
		success = !(nextX == cur_X && nextY == cur_Y)
	} else {
		switch statMove {
		case 0:
			success = true
		case 1:
			searchAvailableRightMove(nextTetro, cur_X, &nextX, &nextY)
			success = !(nextX == cur_X)
		case 2:
			searchAvailableLeftMove(nextTetro, cur_X, &nextX, &nextY)
			success = !(nextX == cur_X)
		}
	}
	if success {
		cur_tetro = nextTetro
		cur_X, cur_Y = nextX, nextY
		createPrediction()
	}
	return success
}

func rotateCounterClockWise() bool {
	var success bool
	nextTetro := cur_tetro.clone()
	nextTetro.rotateCounterClockWise()
	nextX, nextY := cur_X, cur_Y
	statMove := isMoveAvailable(nextTetro, cur_X, cur_X, cur_Y)
	if statMove == 3 {
		searchAvailableYMove(nextTetro, cur_X, &nextX, &nextY)
		success = !(nextX == cur_X && nextY == cur_Y)
	} else {
		switch statMove {
		case 0:
			success = true
		case 1:
			searchAvailableRightMove(nextTetro, cur_X, &nextX, &nextY)
			success = !(nextX == cur_X)
		case 2:
			searchAvailableLeftMove(nextTetro, cur_X, &nextX, &nextY)
			success = !(nextX == cur_X)
		}
	}
	if success {
		cur_tetro = nextTetro
		cur_X, cur_Y = nextX, nextY
		createPrediction()
	}
	return success
}

/*UTILITY*/

func searchAvailableYMove(tetro Tetromino, x int, nextX, nextY *int) {
	if isMoveAvailable(tetro, x, *nextX+2, *nextY+1) == 0 {
		*nextX += 2
		*nextY++
	} else if isMoveAvailable(tetro, x, *nextX-2, *nextY+2) == 0 {
		*nextX -= 2
		*nextY++
	}
}

func searchAvailableRightMove(tetro Tetromino, x int, nextX, nextY *int) {
	nX := *nextX
	for {
		statMove := isMoveAvailable(tetro, x, nX, *nextY)
		if statMove != 1 {
			if statMove == 0 {
				*nextX = nX
			}
			break
		}
		nX += 2
	}
}

func searchAvailableLeftMove(tetro Tetromino, x int, nextX, nextY *int) {
	nX := *nextX
	for {
		statMove := isMoveAvailable(tetro, x, nX, *nextY)
		if statMove != 2 {
			if statMove == 0 {
				*nextX = nX
			}
			break
		}
		nX -= 2
	}
}

func isMoveAvailable(tetro Tetromino, x, nextX, nextY int) int {
	oriX, oriY := tetro.getOriXY()

	if nextX < oriX*-2 {
		return 1
	} else if nextX+(oriX+tetro.width)*2 > 20 {
		return 2
	} else if nextY+oriY+tetro.height > 20 {
		return 3
	}

	if !cur_board.canBindTetromino(tetro, oriX+nextX/2, oriY+nextY) {
		if nextX < x {
			return 1
		} else if nextX > x {
			return 2
		} else {
			return 3
		}
	}
	return 0
}

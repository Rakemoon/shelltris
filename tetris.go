package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

// 45x28
func drawTetrisScreen(scr tcell.Screen, force bool) {
	scr.Fill(0, DEF_SF)
	if term_width < 45 || term_height < 28 {
		is_term_too_small = true
		is_paused = true
		printYourTermIsTooSmall(scr, 45, 28)
	} else {
		is_term_too_small = false
		centerX, centerY := (term_width-45)/2, (term_height-28)/2
		printHeader(scr, centerX, centerY)
		printTetrisBox(scr)
		printScore(scr)
		printNextTetromino(scr)
		printHoldTetromino(scr)
		printKeyBinds(scr)
	}
	if force {
		scr.Sync()
	} else {
		scr.Show()
	}
}

func printTetrisBox(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-28)/2
	printCurBoard := func() {
		cur_board.print(scr, x+1, y+6+1)
		cur_tetro.print(scr, x+1+cur_X, y+6+1+predict_Y, y+6+1, true, "◤◢")
		cur_tetro.print(scr, x+1+cur_X, y+6+1+cur_Y, y+6+1, true, "██")
	}
	if is_game_over {
		printBox(scr, x, y+6, 22, 22, RED_SF, true)
		printCurBoard()
		printGameOverText(scr)
	} else if is_paused {
		printBox(scr, x, y+6, 22, 22, YELLOW_SF, true)
		printCurBoard()
		printPauseText(scr)
	} else {
		printBox(scr, x, y+6, 22, 22, CYAN_SF, true)
		printCurBoard()
	}
}

func printScore(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-28)/2
	printBox(scr, x+23, y+6, 22, 4, DEF_SF, false)
	printText(scr, x+23+2, y+6+1, 22, 1, DEF_SF.Bold(true), "SCORE")
	printText(scr, x+23+2, y+6+2, 22, 1, DEF_SF.Bold(true), "LEVEL")
	printText(scr, x+23+2+6, y+6+1, 22, 1, GREEN_SF.Bold(true), fmt.Sprintf("%12d", score))
	printText(scr, x+23+2+6, y+6+2, 22, 1, getStyleByInt(cur_level).Bold(true), fmt.Sprintf("%12d", cur_level))
}

func printNextTetromino(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-28)/2
	tx, ty, withOri := 0, 0, true
	if next_tetro.color != 1 {
		if next_tetro.color == 3 {
			tx, ty, withOri = 2, 1, false
		} else {
			tx, ty, withOri = 1, 1, false
		}

	}
	var plus int
	if hold_tetro.color != -1 {
		plus += 12
	}
	printBox(scr, x+23+plus, y+6+4, 10, 6, DEF_SF, false)
	printText(scr, x+23+3+plus, y+6+4, 4, 1, DEF_SF.Bold(true), "NEXT")
	next_tetro.print(scr, x+23+1+tx+plus, y+6+4+1+ty, y+6+4, withOri, "██")
}

func printHoldTetromino(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-28)/2
	tx, ty, withOri := 0, 0, true
	if hold_tetro.color != 1 {
		if hold_tetro.color == 3 {
			tx, ty, withOri = 2, 1, false
		} else {
			tx, ty, withOri = 1, 1, false
		}

	}
	if hold_tetro.color != -1 {
		printBox(scr, x+23, y+6+4, 10, 6, DEF_SF, false)
		printText(scr, x+23+3, y+6+4, 4, 1, DEF_SF.Bold(true), "HOLD")
		hold_tetro.print(scr, x+23+1+tx, y+6+4+1+ty, y+6+4, withOri, "██")
	} else {
		printBox(scr, x+23+12, y+6+4, 10, 6, DEF_SF, false)
		printText(scr, x+23+12+3, y+6+4, 4, 1, DEF_SF.Bold(true), "HOLD")
	}
}

func printKeyBinds(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-28)/2
	Keybinds := [][]string{
		{"Left", "←"},
		{"Right", "→"},
		{"Down", "↓"},
		{"Rotate", "↑"},
		{"Rotate*", "z"},
		{"Drop", "space"},
		{"Hold", "x"},
		{"Pause", "p"},
		{"Reset", "r"},
		{"Quit", "q"},
	}
	printBox(scr, x+23, y+6+10, 22, 12, DEF_SF, false)
	printText(scr, x+23+8, y+6+10, 7, 1, DEF_SF.Bold(true), "KEYBIND")
	for ty, binds := range Keybinds {
		name := binds[0]
		key := fmt.Sprintf("%5s", binds[1])
		printText(scr, x+23+3, ty+y+6+10+1, len(name), 1, DEF_SF.Bold(true), name)
		printText(scr, x+23+3+11, ty+y+6+10+1, len(key), 1, DEF_SF.Bold(true), key)
	}
}

func printGameOverText(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-28)/2
	printBox(scr, x+6, y+5+10, 11, 3, RED_SF, false)
	printText(scr, x+6+1, y+5+10+1, 9, 1, RED_SF.Bold(true), "GAME OVER")
}

func printPauseText(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-28)/2
	printBox(scr, x+8, y+5+10, 7, 3, YELLOW_SF, false)
	printText(scr, x+8+1, y+5+10+1, 5, 1, YELLOW_SF.Bold(true), "PAUSE")
}

func dropTetromino(x int) bool {
	if x == -1 {
		cur_tetro = Tetromino{}
		cur_tetro.setRandom()
		next_tetro = Tetromino{}
		next_tetro.setRandom()
		x = 4
	} else {
		if !is_can_hold && hold_tetro.color == -1 {
			hold_tetro = cur_tetro
			cur_tetro = next_tetro
			next_tetro = Tetromino{}
			next_tetro.setRandom()
		} else if !is_can_hold && hold_tetro.color != -1 {
			cur_tetro = hold_tetro
			hold_tetro.color = -1
			is_can_hold = true
		} else {
			cur_tetro = next_tetro
			next_tetro = Tetromino{}
			next_tetro.setRandom()
		}
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
	return true
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
				eliminated := cur_board.updateBoard()
				calculateScore(eliminated)
				printScore(scr)
				if !is_game_over {
					dropTetromino(cur_X)
					printNextTetromino(scr)
					printHoldTetromino(scr)
				}
				is_game_over = cur_Y == predict_Y
			}
			printTetrisBox(scr)
			scr.Show()
			is_elemination_session = false
		}
	} else {
		updated <- true
	}
}

func calculateScore(eliminated int) {
	sc := 10
	for i := 1; i <= eliminated; i++ {
		sc += i
	}
	sc = sc + 2*cur_level
	sc = sc + 2*cur_height
	score += sc
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

package main

import "github.com/gdamore/tcell/v2"

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
	cur_tetro.print(scr, x+1+cur_X, y+6+1+cur_Y, y+6+1, true)
}

// right = true to move right, right = false to move left
func canMoveLeftRight(right bool) bool {
	next_X := cur_X - 2
	if right {
		next_X = cur_X + 2
	}

	if next_X < 0 || next_X+cur_tetro.width*2 > 20 {
		return false
	}
	return true
}

// right = true to move right, right = false to move left
func moveLeftRight(right bool) {
	next_X := cur_X - 2
	if right {
		next_X = cur_X + 2
	}
	cur_X = next_X
}

func canMoveDown() bool {
	next_Y := cur_Y + 1
	if next_Y+cur_tetro.height > 20 {
		return false
	}
	return true
}

func moveDown() {
	cur_Y++
}

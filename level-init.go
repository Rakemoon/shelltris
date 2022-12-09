package main

import (
	"github.com/gdamore/tcell/v2"
)

// 45x23
func drawInitScreen(scr tcell.Screen, force bool) {
	scr.Fill(0, DEF_SF)
	if term_width < 45 || term_height < 23 {
		is_term_too_small = true
		printYourTermIsTooSmall(scr, 45, 23)
	} else {
		is_term_too_small = false
		centerX, centerY := (term_width-45)/2, (term_height-23)/2
		printHeader(scr, centerX, centerY)
		printBoxInitChoose(scr)
		printInitInstruction(scr)
	}
	if force {
		scr.Sync()
	} else {
		scr.Show()
	}
}

func printBoxInitChoose(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-23)/2
	printBox(scr, x+3, y+8, 9*2, 9, DEF_SF, is_select_level)
	printText(scr, x+3+7, y+8, 5, 1, DEF_SF.Bold(is_select_level), "LEVEL")
	printNumber(scr, x+3+6, y+8+2, cur_level)
	printBox(scr, x+3+9*2+3, y+8, 9*2, 9, DEF_SF, !is_select_level)
	printText(scr, x+3+9*2+3+6, y+8, 6, 1, DEF_SF.Bold(!is_select_level), "HEIGHT")
	printNumber(scr, x+3+9*2+3+6, y+8+2, cur_height)
}

func printInitInstruction(scr tcell.Screen) {
	x, y := (term_width-45)/2, (term_height-23)/2
	printBox(scr, x+8, y+19, 28, 4, DEF_SF, false)
	printText(scr, x+8+2, y+19+1, 24, 1, DEF_SF.Bold(true), "USE ARROW KEYS TO SELECT")
	printText(scr, x+8+2+3, y+19+2, 18, 1, DEF_SF.Bold(true), "AND ENTER TO START")
}

func switchInitChoose() {
	is_select_level = !is_select_level
}

func incInitChoose() {
	if is_select_level {
		cur_level++
		if cur_level > MAX_TETRIS_LEVEL {
			cur_level = MIN_TETRIS_LEVEL
		}
	} else {
		cur_height++
		if cur_height > MAX_TETRIS_HEIGHT {
			cur_height = MIN_TETRIS_HEIGHT
		}
	}
}

func decInitChoose() {
	if is_select_level {
		cur_level--
		if cur_level < MIN_TETRIS_LEVEL {
			cur_level = MAX_TETRIS_LEVEL
		}
	} else {
		cur_height--
		if cur_height < MIN_TETRIS_HEIGHT {
			cur_height = MAX_TETRIS_HEIGHT
		}
	}
}

package main

import "github.com/gdamore/tcell/v2"

const (
	MAX_TETRIS_LEVEL  = 6
	MAX_TETRIS_HEIGHT = 6
	MIN_TETRIS_LEVEL  = 0
	MIN_TETRIS_HEIGHT = 0
)

var (
	DEF_SF, BLUE_SF, CYAN_SF, YELLOW_SF, GREEN_SF, MAGENTA_SF, RED_SF, WHITE_SF tcell.Style

	is_end bool

	term_width, term_height int
	is_term_too_small       bool

	is_initialization = true
	is_select_level   = true

	cur_level  = MIN_TETRIS_LEVEL
	cur_height = MIN_TETRIS_HEIGHT

	cur_X, cur_Y int
	cur_tetro    Tetromino

	pressMoveDown = make(chan bool)
)

// Initialize primary style
func initStyle() {
	DEF_SF = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	BLUE_SF = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorBlue)
	CYAN_SF = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorSkyblue)
	YELLOW_SF = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorYellow)
	GREEN_SF = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorGreen)
	MAGENTA_SF = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorPurple)
	RED_SF = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorRed)
	WHITE_SF = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
}

func bindTermSize(width, height int) {
	term_width, term_height = width, height
}

/*
DEF_SF, BLUE_SF, CYAN_SF, YELLOW_SF, GREEN_SF, MAGENTA_SF, RED_SF, WHITE_SF
Return these style from 0 - 7
*/
func getStyleByInt(index int) tcell.Style {
	styleArray := []tcell.Style{WHITE_SF, BLUE_SF, CYAN_SF, YELLOW_SF, GREEN_SF, MAGENTA_SF, RED_SF, DEF_SF}
	return styleArray[index]
}

package main

import "github.com/gdamore/tcell/v2"

var (
	DEF_SF, BLUE_SF, CYAN_SF, YELLOW_SF, GREEN_SF, MAGENTA_SF, RED_SF, WHITE_SF tcell.Style
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

/*
DEF_SF, BLUE_SF, CYAN_SF, YELLOW_SF, GREEN_SF, MAGENTA_SF, RED_SF, WHITE_SF
Return these style from 0 - 7
*/
func getStyleByInt(index int) tcell.Style {
	styleArray := []tcell.Style{DEF_SF, BLUE_SF, CYAN_SF, YELLOW_SF, GREEN_SF, MAGENTA_SF, RED_SF, WHITE_SF}
	if index > 7 {
		index = 0
	}
	if index < 0 {
		index = 7
	}
	return styleArray[index]
}

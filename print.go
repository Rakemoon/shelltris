package main

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

/*TETRIS PRINTING*/

// Print TETRIS
func printHeader(scr tcell.Screen, x, y int) {
	ascii := [][]string{
		{"████████╗", "███████╗", "████████╗", "██████╗ ", "██╗", "███████╗"},
		{"╚══██╔══╝", "██╔════╝", "╚══██╔══╝", "██╔══██╗", "██║", "██╔════╝"},
		{"   ██║   ", "█████╗  ", "   ██║   ", "██████╔╝", "██║", "███████╗"},
		{"   ██║   ", "██╔══╝  ", "   ██║   ", "██╔══██╗", "██║", "╚════██║"},
		{"   ██║   ", "███████╗", "   ██║   ", "██║  ██║", "██║", "███████║"},
		{"   ╚═╝   ", "╚══════╝", "   ╚═╝   ", "╚═╝  ╚═╝", "╚═╝", "╚══════╝"},
	}
	for ty, r := range ascii {
		ctx := 0
		for tx, c := range r {
			cLen := utf8.RuneCountInString(c)
			printText(scr, x+ctx, y+ty, cLen, 1, getStyleByInt(tx+1), c)
			ctx += cLen
		}
	}
}

// Print the Tetris Session
func printTetris(scr tcell.Screen) {
	width, height := scr.Size()
	scr.Fill(0, DEF_SF)

	printHeader(scr, (width-45)/2, (height-6)/2)
}

/*UTILITY*/

// Print the text using string.
func printText(scr tcell.Screen, x, y, width, height int, style tcell.Style, text string) {
	cols, rows := x, y
	for _, c := range []rune(text) {
		scr.SetContent(cols, rows, c, nil, style)
		cols++
		if (cols - x) >= width {
			rows++
			cols = x
		}
		if (rows - y) >= height {
			break
		}
	}
}

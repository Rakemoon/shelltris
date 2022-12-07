package main

import "github.com/gdamore/tcell/v2"

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

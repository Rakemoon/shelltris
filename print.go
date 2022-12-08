package main

import (
	"fmt"
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

/*TETRIS PRINTING*/

// Print TETRIS 45x6
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

func printNumber(scr tcell.Screen, x, y, index int) {
	artNumber := [][][]uint8{
		{
			{1, 1, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 0, 1},
			{1, 1, 1},
		},
		{
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
			{0, 1, 0},
		},
		{
			{1, 1, 1},
			{0, 0, 1},
			{1, 1, 1},
			{1, 0, 0},
			{1, 1, 1},
		},
		{
			{1, 1, 1},
			{0, 0, 1},
			{1, 1, 1},
			{0, 0, 1},
			{1, 1, 1},
		},
		{
			{1, 0, 1},
			{1, 0, 1},
			{1, 1, 1},
			{0, 0, 1},
			{0, 0, 1},
		},
		{
			{1, 1, 1},
			{1, 0, 0},
			{1, 1, 1},
			{0, 0, 1},
			{1, 1, 1},
		},
		{
			{1, 1, 1},
			{1, 0, 0},
			{1, 1, 1},
			{1, 0, 1},
			{1, 1, 1},
		},
	}

	style := getStyleByInt(index)

	for ty, rows := range artNumber[index] {
		for tx, cols := range rows {
			if cols > 0 {
				printText(scr, x+tx*2, y+ty, 2, 1, style, "██")
			}
		}
	}
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

func printBox(scr tcell.Screen, x, y, width, height int, style tcell.Style, heavy bool) {
	boxRunes := [6]rune{'─', '│', '┌', '┐', '└', '┘'}
	if heavy {
		boxRunes = [6]rune{'━', '┃', '┏', '┓', '┗', '┛'}
	}

	for tx := 0; tx < width; tx++ {
		for ty := 0; ty < height; ty++ {
			scr.SetContent(x+tx, y+ty, 0, nil, DEF_SF)
		}
	}

	for tx := 0; tx < width; tx++ {
		scr.SetContent(x+tx, y, boxRunes[0], nil, style)
		scr.SetContent(x+tx, y+height-1, boxRunes[0], nil, style)
	}

	for ty := 0; ty < height; ty++ {
		scr.SetContent(x, y+ty, boxRunes[1], nil, style)
		scr.SetContent(x+width-1, y+ty, boxRunes[1], nil, style)
	}

	if width > 1 && height > 1 {
		scr.SetContent(x, y, boxRunes[2], nil, style)
		scr.SetContent(x+width-1, y, boxRunes[3], nil, style)
		scr.SetContent(x, y+height-1, boxRunes[4], nil, style)
		scr.SetContent(x+width-1, y+height-1, boxRunes[5], nil, style)
	}
}

func printYourTermIsTooSmall(scr tcell.Screen, exWidth, exHeight int) {
	cX, cY := term_width/2, term_height/2
	ex := fmt.Sprintf("EXPECTED: %dx%d", exWidth, exHeight)
	got := fmt.Sprintf("GOT     : %dx%d", term_width, term_height)

	sep := len(ex) / 2
	printBox(scr, 0, 0, term_width, term_height, RED_SF, false)
	printText(scr, cX-13, cY-1, term_width, 1, DEF_SF.Bold(true), "TERMINAL SIZE IS TOO SMALL!")
	printText(scr, cX-sep, cY, term_width, 1, DEF_SF.Bold(true), ex)
	printText(scr, cX-sep, cY+1, term_width, 1, DEF_SF.Bold(true), got)
}

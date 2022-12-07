package main

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v2"
)

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

func printTetris(scr tcell.Screen) {
	width, height := scr.Size()
	scr.Fill(0, DEF_SF)

	printHeader(scr, (width-45)/2, (height-6)/2)
}

func printTetrisPlease(scr tcell.Screen) {
	printTetris(scr)
	scr.Show()
}

func printTetrisNow(scr tcell.Screen) {
	printTetris(scr)
	scr.Sync()
}

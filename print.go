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

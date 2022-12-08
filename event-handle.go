package main

import "github.com/gdamore/tcell/v2"

func handleTermEvent(scr tcell.Screen) {
	for !is_end {
		event := scr.PollEvent()
		switch event := event.(type) {
		case *tcell.EventResize:
			onResize(scr, event)
		case *tcell.EventKey:
			key := event.Key()
			if key == tcell.KeyRune {
				onPressRune(scr, event, event.Rune())
			} else {
				onPressKey(scr, event, key)
			}
		}
	}
}

func onResize(scr tcell.Screen, event *tcell.EventResize) {
	bindTermSize(event.Size())
	if is_initialization {
		drawInitScreen(scr, true)
	}
}

func onPressRune(scr tcell.Screen, event *tcell.EventKey, c rune) {
	if c == 'q' || c == 'Q' {
		is_end = true

	}
}

func onPressKey(scr tcell.Screen, event *tcell.EventKey, key tcell.Key) {
	if is_initialization {
		if key == tcell.KeyTAB || key == tcell.KeyLeft || key == tcell.KeyRight {
			switchInitChoose()
		} else if key == tcell.KeyUp {
			incInitChoose()
		} else if key == tcell.KeyDown {
			decInitChoose()
		}
		printBoxInitChoose(scr)
		scr.Show()
	}
}

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
	} else {
		drawTetrisScreen(scr, true)
	}
}

func onPressRune(scr tcell.Screen, event *tcell.EventKey, c rune) {
	if c == 'q' || c == 'Q' {
		is_end = true
		return
	}
	if !is_initialization && !is_term_too_small {
		if c == 'c' {
			cur_tetro.setRandom()
			printTetrisBox(scr)
			scr.Show()
		} else if c == 'z' {
			if rotateCounterClockWise() {
				printTetrisBox(scr)
				scr.Show()
			}
		}
	}
}

func onPressKey(scr tcell.Screen, event *tcell.EventKey, key tcell.Key) {
	if key == tcell.KeyCtrlC {
		is_end = true
		return
	}
	if is_initialization && !is_term_too_small {
		if key == tcell.KeyTAB || key == tcell.KeyLeft || key == tcell.KeyRight {
			switchInitChoose()
			printBoxInitChoose(scr)
			scr.Show()
		} else if key == tcell.KeyUp {
			incInitChoose()
			printBoxInitChoose(scr)
			scr.Show()
		} else if key == tcell.KeyDown {
			decInitChoose()
			printBoxInitChoose(scr)
			scr.Show()
		} else if key == tcell.KeyEnter {
			is_initialization = false
			drawTetrisScreen(scr, false)
		}
	} else if !is_initialization && !is_term_too_small {
		if key == tcell.KeyLeft || key == tcell.KeyRight {
			if moveLeftRight(key == tcell.KeyRight) {
				printTetrisBox(scr)
				scr.Show()
			}
		} else if key == tcell.KeyUp {
			if rotateClockWise() {
				printTetrisBox(scr)
				scr.Show()
			}
		} else if key == tcell.KeyDown {
			if moveDown() {
				printTetrisBox(scr)
				scr.Show()
			}
		}
	}
}

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
}

func onPressKey(scr tcell.Screen, event *tcell.EventKey, key tcell.Key) {
	if key == tcell.KeyCtrlC {
		is_end = true
		return
	}
	if is_initialization {
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
	} else {
		if key == tcell.KeyLeft {
			cur_X -= 2
			printTetrisBox(scr)
			scr.Show()
		} else if key == tcell.KeyRight {
			cur_X += 2
			printTetrisBox(scr)
			scr.Show()
		} else if key == tcell.KeyUp {
			cur_Y--
			printTetrisBox(scr)
			scr.Show()
		} else if key == tcell.KeyDown {
			cur_Y++
			printTetrisBox(scr)
			scr.Show()
		}
	}
}
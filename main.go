package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// you know its like go please! and then magicly program run
func main() {
	scr, error := tcell.NewScreen()
	if error != nil {
		log.Fatal(error)
	}
	if error = scr.Init(); error != nil {
		log.Fatal(error)
	}
	quit := func() {
		isPanic := recover()
		scr.Fini()
		if isPanic != nil {
			panic(isPanic)
		}
	}
	defer quit()

	rand.Seed(time.Now().UnixNano())

	initStyle()
	bindTermSize(scr.Size())

	drawInitScreen(scr, false)

	handleTermEvent(scr)
}

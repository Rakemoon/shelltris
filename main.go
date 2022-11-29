package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	. "github.com/gbin/goncurses"
)

type Tetromino struct {
	con    [][]int
	pair   int16
	height int
	width  int
	rotate bool
}

var (
	curY, curX       int
	oldY, oldX       int
	curOriY, curOriX int
	oldOriY, oldOriX int
	ori              int
	curTetro         Tetromino
	oldTetro         Tetromino
	isEnd            bool

	wg sync.WaitGroup
)

func initConf(win *Window) {
	rand.Seed(time.Now().UnixNano())

	StartColor()
	UseDefaultColors()
	Cursor(0)
	Echo(false)

	InitPair(1, -1, C_CYAN)
	InitPair(2, -1, C_BLUE)
	InitPair(3, -1, C_YELLOW)
	InitPair(4, -1, C_GREEN)
	InitPair(5, -1, C_MAGENTA)
	InitPair(6, -1, C_RED)
	InitPair(7, -1, C_WHITE)

	InitPair(11, C_CYAN, -1)
	InitPair(12, C_BLUE, -1)
	InitPair(13, C_YELLOW, -1)
	InitPair(14, C_GREEN, -1)
	InitPair(15, C_MAGENTA, -1)
	InitPair(16, C_RED, -1)
	InitPair(17, C_WHITE, -1)

	win.Keypad(true)

	curY, curX = 1, 1
	ori = 1
	curTetro = createTetromino()
	BackupTetris()
	SetOriTetromino()
}

func main() {
	stdscr, stderr := Init()
	if stderr != nil {
		log.Fatal(stderr)
	}
	defer End()

	initConf(stdscr)

	headerPrint(stdscr, 0, 0)

	win, _ := NewWindow(21, 22, 6, 0)
	win.Box(0, 0)
	win.Refresh()

	UpdateTetris(win)

	wg.Add(1)
	go handleKeyboard(stdscr, win)

	go GoDownPLEASE(win)

	wg.Wait()
}

func GoDownPLEASE(win *Window) {
	for !isEnd {
		time.Sleep(time.Second / 1)
		MoveTetrominoDown()
		UpdateTetris(win)
	}
}

func handleKeyboard(stdscr *Window, win *Window) {
	defer wg.Done()
	for !isEnd {
		input := stdscr.GetChar()
		switch Key(input) {
		default:
			continue
		case 'q':
			isEnd = true
		case KEY_LEFT:
			MoveTetrominoLeft()
		case KEY_RIGHT:
			MoveTetrominoRight()
		case KEY_DOWN:
			MoveTetrominoDown()
		case KEY_UP:
			TetrominoRotateCLockWise()
		case 'z':
			TetrominoRotateCounterCLockWise()
		case ' ':
			MoveTetrominoDown(true)
		case 'c':
			BackupTetris()
			curTetro = createTetromino()
			ori = 1
			SetOriTetromino()
		case 'r':
			curY = 1
		}
		UpdateTetris(win)
	}
}

func UpdateTetris(win *Window) {
	for ty, con := range oldTetro.con {
		for tx, s := range con {
			if s > 0 {
				win.MovePrint(oldY+oldOriY+ty, oldX+oldOriX+tx*2, "  ")
			}
		}
	}
	win.Refresh()
	printTetromino(win, curTetro, curY+curOriY, curX+curOriX)
}

func BackupTetris() {
	oldX, oldY, oldOriX, oldOriY, oldTetro = curX, curY, curOriX, curOriY, curTetro
}

func MoveTetrominoLeft() {
	BackupTetris()
	next := curX - 2
	if next+curOriX > 0 {
		curX = next
	}
}

func MoveTetrominoRight() {
	BackupTetris()
	tetroW := curTetro.width*2 + curOriX
	next := curX + 2
	if next+tetroW < 22 {
		curX = next
	}
}

func MoveTetrominoDown(forces ...bool) {
	BackupTetris()
	tetroH := curTetro.height + curOriY
	next := curY + 1
	if next+tetroH < 21 {
		if len(forces) > 0 && forces[0] {
			for next+tetroH < 20 {
				next++
			}
		}
		curY = next
	}
}

func SetOriTetromino() {
	if ori < 1 {
		ori = 4
	}
	if ori > 4 {
		ori = 1
	}
	curOriY, curOriX = getPosOri(ori)
}

func getPosOri(val int) (int, int) {
	var oriY, oriX int
	switch val {
	case 1:
		oriY, oriX = 0, 0
	case 2:
		oriY, oriX = 0, 2
	case 3:
		oriY, oriX = 1, 0
	case 4:
		oriY, oriX = 0, 0
	}
	return oriY, oriX
}

func TetrominoRotateCLockWise() {
	BackupTetris()
	nextOriY, nextOriX := getPosOri(ori + 1)
	if oldX+nextOriX < 1 || oldX+nextOriX+oldTetro.height*2 > 21 || oldTetro.width+nextOriY+oldY > 20 {
		return
	}
	curTetro.con = make([][]int, oldTetro.width)
	for w := 0; w < oldTetro.width; w++ {
		curTetro.con[w] = make([]int, oldTetro.height)
		for h := 0; h < oldTetro.height; h++ {
			curTetro.con[w][h] = oldTetro.con[h][w]
		}
	}
	for _, con := range curTetro.con {
		for i, j := 0, oldTetro.height-1; i < j; i, j = i+1, j-1 {
			con[i], con[j] = con[j], con[i]
		}
	}
	if curTetro.rotate {
		ori++
		SetOriTetromino()
	}
	curTetro.height, curTetro.width = oldTetro.width, oldTetro.height
}

func TetrominoRotateCounterCLockWise() {
	BackupTetris()
	nextOriY, nextOriX := getPosOri(ori - 1)
	if oldX+nextOriX < 1 || oldX+nextOriX+oldTetro.height*2 > 21 || oldTetro.width+nextOriY+oldY > 20 {
		return
	}
	curTetro.con = make([][]int, oldTetro.width)
	for w := 0; w < oldTetro.width; w++ {
		curTetro.con[w] = make([]int, oldTetro.height)
		for h := 0; h < oldTetro.height; h++ {
			curTetro.con[w][h] = oldTetro.con[h][w]
		}
	}
	for i, j := 0, oldTetro.width-1; i < j; i, j = i+1, j-1 {
		curTetro.con[i], curTetro.con[j] = curTetro.con[j], curTetro.con[i]
	}
	if curTetro.rotate {
		ori--
		SetOriTetromino()
	}
	curTetro.height, curTetro.width = oldTetro.width, oldTetro.height
}

func headerPrint(win *Window, y int, x int) {
	defer win.Refresh()
	ascii := [][]string{
		{"████████╗", "███████╗", "████████╗", "██████╗ ", "██╗", "███████╗"},
		{"╚══██╔══╝", "██╔════╝", "╚══██╔══╝", "██╔══██╗", "██║", "██╔════╝"},
		{"   ██║   ", "█████╗  ", "   ██║   ", "██████╔╝", "██║", "███████╗"},
		{"   ██║   ", "██╔══╝  ", "   ██║   ", "██╔══██╗", "██║", "╚════██║"},
		{"   ██║   ", "███████╗", "   ██║   ", "██║  ██║", "██║", "███████║"},
		{"   ╚═╝   ", "╚══════╝", "   ╚═╝   ", "╚═╝  ╚═╝", "╚═╝", "╚══════╝"},
	}

	for sy, con := range ascii {
		win.Move(y+sy, x)
		for sx, s := range con {
			pair := int16(sx + 11)
			win.ColorOn(pair)
			win.Print(s)
			win.ColorOff(pair)
		}
	}
}

func printTetromino(win *Window, tetro Tetromino, y int, x int) {
	defer win.Refresh()
	for ty, con := range tetro.con {
		for tx, s := range con {
			if s > 0 {
				win.ColorOn(tetro.pair)
				win.MovePrint(y+ty, x+tx*2, "  ")
				win.ColorOff(tetro.pair)
			}
		}
	}
}

func createTetromino() Tetromino {
	var t Tetromino
	switch rand.Intn(7) {
	case 0:
		t.con, t.pair = [][]int{
			{1, 2, 1, 1},
		}, 1
		t.height, t.width = 1, 4
		t.rotate = false
	case 1:
		t.con, t.pair = [][]int{
			{1, 0, 0},
			{1, 2, 1},
		}, 2
		t.height, t.width = 2, 3
		t.rotate = true
	case 2:
		t.con, t.pair = [][]int{
			{0, 0, 1},
			{1, 2, 1},
		}, 7
		t.height, t.width = 2, 3
		t.rotate = true
	case 3:
		t.con, t.pair = [][]int{
			{1, 1},
			{2, 1},
		}, 3
		t.height, t.width = 2, 2
		t.rotate = false
	case 4:
		t.con, t.pair = [][]int{
			{0, 1, 1},
			{1, 2, 0},
		}, 4
		t.height, t.width = 2, 3
		t.rotate = true
	case 5:
		t.con, t.pair = [][]int{
			{0, 1, 0},
			{1, 2, 1},
		}, 5
		t.height, t.width = 2, 3
		t.rotate = true
	case 6:
		t.con, t.pair = [][]int{
			{1, 1, 0},
			{0, 2, 1},
		}, 6
		t.height, t.width = 2, 3
		t.rotate = true
	}

	return t
}

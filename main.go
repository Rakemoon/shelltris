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
	curTetro         Tetromino
	oldTetro         Tetromino
	curY, curX       int
	oldY, oldX       int
	curOriY, curOriX int
	oldOriY, oldOriX int
	ori              int
	isEnd            bool
	isTetroDown      bool

	isPressedSpace bool
	isPressedDown  bool

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

	isPressedSpace = false

	InitTetromino()
}

func InitTetromino() {
	curTetro = createTetromino()
	oldX, oldY = 7, 1
	curX, curY = 7, 1
	ori = 1
	curOriX, curOriY = 0, 0
	oldOriX, oldOriY = 0, 0
	isTetroDown = false
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
		time.Sleep(time.Second)
		if isPressedDown {
			isPressedDown = false
		} else {
			BackupTetris()
			MoveTetrominoDown()
			UpdateTetris(win)
		}
	}
}

func handleKeyboard(stdscr *Window, win *Window) {
	defer wg.Done()
	for !isEnd {
		input := stdscr.GetChar()
		BackupTetris()
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
			isPressedDown = true
			MoveTetrominoDown()
		case KEY_UP:
			TetrominoRotateCLockWise()
		case 'z':
			TetrominoRotateCounterCLockWise()
		case ' ':
			isPressedSpace = true
			isPressedDown = true
			MoveTetrominoDown()
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
	if isTetroDown {
		InitTetromino()
		printTetromino(win, curTetro, curY+curOriY, curX+curOriX)
	}
}

func BackupTetris() {
	oldX, oldY, oldOriX, oldOriY, oldTetro = curX, curY, curOriX, curOriY, curTetro
}

func MoveTetrominoLeft() {
	next := curX - 2
	if next+curOriX > 0 {
		curX = next
	}
}

func MoveTetrominoRight() {
	tetroW := curTetro.width*2 + curOriX
	next := curX + 2
	if next+tetroW < 22 {
		curX = next
	}
}

func MoveTetrominoDown() {
	tetroH := curTetro.height + curOriY
	next := curY + 1
	if isPressedSpace {
		for next+tetroH < 20 {
			next++
		}
		curY = next
		isPressedSpace = false
	} else if next+tetroH < 21 {
		curY = next
	}

	if next+tetroH == 20 {
		isTetroDown = true
	}
}

func getPosOri(val int) (int, int, int) {
	if val < 1 {
		val = 4
	}
	if val > 4 {
		val = 1
	}
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
	return val, oriY, oriX
}

func TetrominoRotateCLockWise() {
	val, nextOriY, nextOriX := 1, 0, 0
	if curTetro.rotate {
		val, nextOriY, nextOriX = getPosOri(ori + 1)
	}
	if curX+nextOriX < 1 || curX+nextOriX+curTetro.height*2 > 21 || curTetro.width+nextOriY+curY > 20 {
		return
	}
	nextTetro := curTetro
	nextTetro.con = make([][]int, curTetro.width)
	for w := 0; w < oldTetro.width; w++ {
		nextTetro.con[w] = make([]int, curTetro.height)
		for h := 0; h < curTetro.height; h++ {
			nextTetro.con[w][h] = curTetro.con[h][w]
		}
	}
	for _, con := range nextTetro.con {
		for i, j := 0, curTetro.height-1; i < j; i, j = i+1, j-1 {
			con[i], con[j] = con[j], con[i]
		}
	}
	if curTetro.rotate {
		ori = val
		curOriY, curOriX = nextOriY, nextOriX
	}
	nextTetro.height, nextTetro.width = curTetro.width, curTetro.height
	curTetro = nextTetro
}

func TetrominoRotateCounterCLockWise() {
	val, nextOriY, nextOriX := 1, 0, 0
	if curTetro.rotate {
		val, nextOriY, nextOriX = getPosOri(ori + 1)
	}
	if curX+nextOriX < 1 || curX+nextOriX+curTetro.height*2 > 21 || curTetro.width+nextOriY+curY > 20 {
		return
	}
	nextTetro := curTetro
	nextTetro.con = make([][]int, curTetro.width)
	for w := 0; w < oldTetro.width; w++ {
		nextTetro.con[w] = make([]int, curTetro.height)
		for h := 0; h < curTetro.height; h++ {
			nextTetro.con[w][h] = curTetro.con[h][w]
		}
	}
	for i, j := 0, curTetro.width-1; i < j; i, j = i+1, j-1 {
		nextTetro.con[i], nextTetro.con[j] = nextTetro.con[j], nextTetro.con[i]
	}
	if curTetro.rotate {
		ori = val
		curOriY, curOriX = nextOriY, nextOriX
	}
	nextTetro.height, nextTetro.width = curTetro.width, curTetro.height
	curTetro = nextTetro
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

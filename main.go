package main

import (
	"log"
	"math"
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

type BoardProp struct {
	val  int
	pair int16
}

var (
	board            [][]BoardProp
	curTetro         Tetromino
	oldTetro         Tetromino
	nextTetro        Tetromino
	curY, curX       int
	oldY, oldX       int
	curOriY, curOriX int
	oldOriY, oldOriX int
	ori              int
	isEnd            bool
	isTetroDown      bool

	isPressedSpace bool
	isPressedDown  bool

	oldPredY int

	lastMovement uint8
	score        int
	level        = 1

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

	InitBoard()
	InitTetromino(7)
}

func InitBoard() {
	board = make([][]BoardProp, 20)
	for i := 0; i < 20; i++ {
		board[i] = make([]BoardProp, 10)
	}
}

func InitTetromino(posX int) {
	curTetro = nextTetro
	if nextTetro.height < 1 {
		curTetro = createTetromino()
	}
	nextTetro = createTetromino()
	oldX, oldY = posX, 2-curTetro.height
	curX, curY = posX, 2-curTetro.height
	ori = 1
	curOriX, curOriY = 0, 0
	oldOriX, oldOriY = 0, 0
	oldPredY = 0
	isTetroDown = false

	for curX+curOriX+curTetro.width*2 >= 22 {
		curX -= 2
	}

	for curX+curOriX < 1 {
		curX += 2
	}
}

func InitializeLevel(stdscr *Window, win *Window) int {
	artNumber := [][][]uint8{
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
	val := 1
	win.Box(0, 0)
	win.Refresh()
	defer func() {
		win.Erase()
		win.Refresh()
		win.Delete()
	}()
	for {
		artNum := artNumber[val-1]
		win.Erase()
		win.Box(0, 0)
		win.AttrOn(A_BOLD)
		win.MovePrint(0, 7, "Chose Level!")
		win.AttrOff(A_BOLD)
		win.Refresh()

		win.ColorOn(int16(val))
		for y, row := range artNum {
			for x, col := range row {
				if col > 0 {
					win.MovePrint(y+2, 10+x*2, "  ")
				}
			}
		}
		win.ColorOff(int16(val))
		win.Refresh()
		c := stdscr.GetChar()
		switch Key(c) {
		case 'q':
			val = 0
			return val
		case ' ':
			return val
		case KEY_UP:
			val++
		case KEY_RIGHT:
			val++
		case KEY_DOWN:
			val--
		case KEY_LEFT:
			val--
		}
		if val > 6 {
			val = 1
		}
		if val < 1 {
			val = 6
		}
	}
}

func main() {
	stdscr, stderr := Init()
	if stderr != nil {
		log.Fatal(stderr)
	}
	defer End()

	initConf(stdscr)

	mY, mX := stdscr.MaxYX()

	mY, mX = (mY/2)-14, (mX/2)-22

	headerPrint(stdscr, mY, mX)

	winit, _ := NewWindow(9, 26, mY+8, mX+10)
	winit2, _ := NewWindow(4, 26, mY+17, mX+10)
	winit2.Box(0, 0)
	winit2.AttrOn(A_BOLD)
	winit2.MovePrint(1, 1, " USE ARROW KEYS TO MOVE")
	winit2.MovePrint(2, 1, "   AND SPACE TO SELECT")
	winit2.AttrOff(A_BOLD)
	winit2.Refresh()
	if level = InitializeLevel(stdscr, winit); level == 0 {
		return
	}
	winit2.Erase()
	winit2.Refresh()
	winit2.Delete()

	stdscr.Refresh()

	win, _ := NewWindow(22, 22, mY+6, mX)
	win.Box(0, 0)
	win.Refresh()

	win2, _ := NewWindow(4, 22, mY+6, mX+23)
	win2.Box(0, 0)
	win2.Refresh()
	PrintScore(win2)

	win3, _ := NewWindow(11, 22, mY+17, mX+23)
	win3.Box(0, 0)
	win3.AttrOn(A_BOLD)
	win3.MovePrint(0, 8, "CONTROL")
	win3.MovePrintf(2, 4, "Left%10s", "←")
	win3.MovePrintf(3, 4, "Right%9s", "→")
	win3.MovePrintf(4, 4, "Down%10s", "↓")
	win3.MovePrintf(5, 4, "Rotate%8s", "↑")
	win3.MovePrintf(6, 4, "Rotate*%7s", "z")
	win3.MovePrintf(7, 4, "Drop%10s", "space")
	win3.MovePrintf(8, 4, "Quit%10s", "q")
	win3.AttrOff(A_BOLD)
	win3.Refresh()

	win4, _ := NewWindow(7, 22, mY+10, mX+23)
	PrintNextTetromino(win4)

	UpdateTetris(win, win2, win4)

	go GoDownPLEASE(win, win2, win4)

	wg.Add(1)
	go handleKeyboard(stdscr, win, win2, win4)

	wg.Wait()
}

func GoDownPLEASE(win *Window, win2 *Window, win3 *Window) {
	for !isEnd {
		lastMovement = 0
		time.Sleep(time.Second / time.Duration(3+level))
		if isPressedDown {
			isPressedDown = false
		} else {
			BackupTetris()
			MoveTetrominoDown()
			UpdateTetris(win, win2, win3)
		}
	}
}

func handleKeyboard(stdscr *Window, win *Window, win2 *Window, win3 *Window) {
	defer wg.Done()
	for !isEnd {
		stdscr.Refresh()
		input := stdscr.GetChar()
		BackupTetris()
		switch Key(input) {
		default:
			continue
		case 'q':
			isEnd = true
		case KEY_LEFT:
			if moved := MoveTetrominoLeft(); !moved {
				continue
			}
			lastMovement = 1
		case KEY_RIGHT:
			if moved := MoveTetrominoRight(); !moved {
				continue
			}
			lastMovement = 2
		case KEY_DOWN:
			isPressedDown = true
			MoveTetrominoDown()
			lastMovement = 3
		case KEY_UP:
			TetrominoRotateCLockWise()
			lastMovement = 4
		case 'z':
			TetrominoRotateCounterCLockWise()
			lastMovement = 5
		case ' ':
			isPressedSpace = true
			isPressedDown = true
			MoveTetrominoDown()
			lastMovement = 6
		}
		UpdateTetris(win, win2, win3)
		if isEnd {
			stdscr.ColorOn(16)
			stdscr.MovePrint(30, 0, "GAME OVER")
			stdscr.ColorOff(16)
		}
	}
}

func UpdateTetris(win *Window, win2 *Window, win3 *Window) {
	if isEnd {
		return
	}
	for ty, con := range oldTetro.con {
		for tx, s := range con {
			if s > 0 && oldY+oldOriY+ty >= 1 {
				win.MovePrint(oldY+oldOriY+ty, oldX+oldOriX+tx*2, "  ")
			}
		}
	}
	win.Refresh()
	PrintPrediction(win)
	printTetromino(win, curTetro, curY+curOriY, curX+curOriX)
	if isTetroDown {
	storeToBoard:
		for ty, con := range curTetro.con {
			for tx, s := range con {
				if s > 0 {
					y := curY + curOriY - 1
					x := (curX + curOriX) / 2
					if y < 0 || x < 0 {
						isEnd = true
						break storeToBoard
					}
					board[ty+y][tx+x].val += s
					board[ty+y][tx+x].pair = curTetro.pair
				}
			}
		}
		CheckAndEliminateRow(win, win2)
		InitTetromino(curX)
		PrintNextTetromino(win3)
		//PrintPrediction(win)
		printTetromino(win, curTetro, curY+curOriY, curX+curOriX)
	}
	for _, s := range board[0] {
		if s.val > 0 {
			isEnd = true
		}
	}
}

func PrintNextTetromino(win *Window) {
	win.Erase()
	win.Box(0, 0)
	win.AttrOn(A_BOLD)
	win.MovePrint(0, 9, "Next")
	win.AttrOff(A_BOLD)
	win.Refresh()
	printTetromino(win, nextTetro, 4-nextTetro.height, 11-nextTetro.width)
}

func CheckAndEliminateRow(win *Window, win2 *Window) {
	nextBoard := make([][]BoardProp, 0)
	eliminated := 0
	for _, row := range board {
		var count int
		for _, col := range row {
			count += col.val
		}
		if count < 10 {
			nextBoard = append(nextBoard, row)
		} else {
			eliminated++
		}
	}
	if eliminated > 0 {
		nNextBoard := make([][]BoardProp, eliminated)
		for i := 0; i < eliminated; i++ {
			nNextBoard[i] = make([]BoardProp, 10)
		}
		nextBoard = append(nNextBoard, nextBoard...)
		for y, row := range nextBoard {
			for x, col := range row {
				if col.pair > 0 {
					win.ColorOn(col.pair)
				}
				win.MovePrint(1+y, 1+x*2, "  ")
				if col.pair > 0 {
					win.ColorOff(col.pair)
				}
			}
		}
		CalculateScore(eliminated)
		PrintScore(win2)
	}
	board = nextBoard
}

func CalculateScore(eliminated int) {
	incr := level
	if lastMovement == 4 || lastMovement == 5 {
		incr += 2
	}
	incr += int(math.Ceil(float64(eliminated+1) / 2))
	score += 22 * incr
}

func PrintScore(win *Window) {
	win.AttrOn(A_BOLD)
	win.MovePrint(1, 1, "Score ")
	win.MovePrint(2, 1, "Level ")
	win.ColorOn(14)
	win.MovePrintf(1, 7, "%14d", score)
	win.ColorOff(14)
	win.ColorOn(12)
	win.MovePrintf(2, 7, "%14d", level)
	win.ColorOff(12)
	win.AttrOff(A_BOLD)
	win.Refresh()
}

func BackupTetris() {
	oldX, oldY, oldOriX, oldOriY, oldTetro = curX, curY, curOriX, curOriY, curTetro
}

func MoveTetrominoRight() bool {
	if IsCanMoveRight() {
		curX += 2
		return true
	}
	return false
}

func MoveTetrominoLeft() bool {
	if IsCanMoveLeft() {
		curX -= 2
		return true
	}
	return false
}

func IsCanMoveDown(yS int) bool {
	if yS+curTetro.height+curOriY+1 < 22 {
		if yS >= 0 {
			y := yS + curOriY
			x := (curX + curOriX) / 2
			for i := 0; i < curTetro.height; i++ {
				for j := 0; j < curTetro.width; j++ {
					sum := curTetro.con[i][j] + board[i+y][x+j].val
					if sum > 1 {
						return false
					}
				}
			}
		}
		return true
	}
	return false
}

func IsCanMoveRight() bool {
	if curX+curOriX+2+curTetro.width*2 < 22 {
		if curY > 0 {
			y := curY + curOriY - 1
			x := ((curX + curOriX) / 2) + 1
			for i := 0; i < curTetro.height; i++ {
				for j := 0; j < curTetro.width; j++ {
					sum := curTetro.con[i][j] + board[i+y][x+j].val
					if sum > 1 {
						return false
					}
				}
			}
		}
		return true
	}
	return false
}

func IsCanMoveLeft() bool {
	if curX+curOriX-2 > 0 {
		if curY > 0 {
			y := curY + curOriY - 1
			x := ((curX + curOriX) / 2) - 1
			for i := 0; i < curTetro.height; i++ {
				for j := 0; j < curTetro.width; j++ {
					sum := curTetro.con[i][j] + board[i+y][x+j].val
					if sum > 1 {
						return false
					}
				}
			}
		}
		return true
	}
	return false
}

func MoveTetrominoDown() {
	if isPressedSpace {
		for IsCanMoveDown(curY) {
			curY++
		}
		isPressedSpace = false
		isTetroDown = true
	} else if IsCanMoveDown(curY) {
		curY++
	} else {
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

func PrintPrediction(win *Window) {
	if oldPredY > 0 {
		for ty, con := range oldTetro.con {
			for tx, s := range con {
				if s > 0 {
					win.MovePrint(oldPredY+oldOriY+ty, oldX+oldOriX+tx*2, "  ")
				}
			}
		}
	}
	y := curY
	for IsCanMoveDown(y) {
		y++
	}
	oldPredY = y
	printTetromino(win, curTetro, y+curOriY, curOriX+curX, "◤◢", 10)
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

func printTetromino(win *Window, tetro Tetromino, y int, x int, args ...interface{}) {
	defer win.Refresh()
	charSt := "  "
	plusPair := 0
	if len(args) > 0 {
		charSt = args[0].(string)
		plusPair = args[1].(int)
	}
	for ty, con := range tetro.con {
		for tx, s := range con {
			if s > 0 && y+ty >= 1 {
				win.ColorOn(tetro.pair + int16(plusPair))
				win.MovePrint(y+ty, x+tx*2, charSt)
				win.ColorOff(tetro.pair + int16(plusPair))
			}
		}
	}
}

func createTetromino() Tetromino {
	var t Tetromino
	switch rand.Intn(7) {
	case 0:
		t.con, t.pair = [][]int{
			{1, 1, 1, 1},
		}, 1
		t.height, t.width = 1, 4
		t.rotate = false
	case 1:
		t.con, t.pair = [][]int{
			{1, 0, 0},
			{1, 1, 1},
		}, 2
		t.height, t.width = 2, 3
		t.rotate = true
	case 2:
		t.con, t.pair = [][]int{
			{0, 0, 1},
			{1, 1, 1},
		}, 7
		t.height, t.width = 2, 3
		t.rotate = true
	case 3:
		t.con, t.pair = [][]int{
			{1, 1},
			{1, 1},
		}, 3
		t.height, t.width = 2, 2
		t.rotate = false
	case 4:
		t.con, t.pair = [][]int{
			{0, 1, 1},
			{1, 1, 0},
		}, 4
		t.height, t.width = 2, 3
		t.rotate = true
	case 5:
		t.con, t.pair = [][]int{
			{0, 1, 0},
			{1, 1, 1},
		}, 5
		t.height, t.width = 2, 3
		t.rotate = true
	case 6:
		t.con, t.pair = [][]int{
			{1, 1, 0},
			{0, 1, 1},
		}, 6
		t.height, t.width = 2, 3
		t.rotate = true
	}

	return t
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

//Input

const PaddleHeight = 4
const PaddleSymbol = 0x2588

type Paddle struct {
	row, col, width, height int
}

var screen tcell.Screen
var Player1 *Paddle
var Player2 *Paddle
var debugLog string

func PrintString(row, col int, str string) {
	for _, c := range str {
		screen.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(row, col, width, height int, ch rune) {

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			screen.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}

	}

}

func DrawState() {
	screen.Clear()
	PrintString(0, 0, debugLog)
	Print(Player1.row, Player1.col, Player1.width, Player1.height, PaddleSymbol)
	Print(Player2.row, Player2.col, Player2.width, Player2.height, PaddleSymbol)
	screen.Show()
}

func main() {
	Initscreen()
	InitGameState()
	inputchan := InitUserInput()

	for {
		DrawState()
		time.Sleep(50 * time.Millisecond)
		key := readInput(inputchan)
		if key == "Rune[q]" {
			screen.Fini()
			os.Exit(0)
		} else if key == "Rune[w]" {
			Player1.row--
		} else if key == "Rune[s]" {
			Player1.row++
		} else if key == "Up" {
			Player2.row--
		} else if key == "Down" {
			Player2.row++
		}
	}
}

func InitUserInput() chan string {
	inputchan := make(chan string)
	go func() {
		for {
			switch ev := screen.PollEvent().(type) {
			case *tcell.EventKey:
				inputchan <- ev.Name()
			}
		}
	}()
	return inputchan
}

func Initscreen() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)
}

func InitGameState() {
	width, height := screen.Size()
	paddleStart := height/2 - PaddleHeight/2
	Player1 = &Paddle{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: PaddleHeight,
	}
	Player2 = &Paddle{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: PaddleHeight,
	}
}

func readInput(inputchan chan string) string {
	var key string
	select {
	case key = <-inputchan:
	default:
		key = ""
	}
	return key
}

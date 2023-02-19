package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

//Collides with Paddle

const PaddleHeight = 4
const PaddleSymbol = 0x2588
const BallSymbol = 0x25CF
const BallVelRow = 1
const BallVelCol = 2

type GameObject struct {
	row, col, width, height int
	VelRow, VelCol          int
	symbol                  rune
}

var screen tcell.Screen
var Player1Paddle *GameObject
var Player2Paddle *GameObject
var ball *GameObject
var debugLog string

var gameObject []*GameObject

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
	for _, obj := range gameObject {
		Print(obj.row, obj.col, obj.width, obj.height, obj.symbol)
	}
	screen.Show()
}

func CollidesWithWall(obj *GameObject) bool {
	_, screenHeight := screen.Size()
	return obj.row+obj.VelRow < 0 || obj.row+obj.VelRow > screenHeight
}

func main() {
	Initscreen()
	InitGameState()
	inputchan := InitUserInput()

	for {
		HandleUserInput(readInput(inputchan))
		UpdateState()
		DrawState()
		time.Sleep(50 * time.Millisecond)

	}
}

func UpdateState() {
	for i := range gameObject {
		gameObject[i].row += gameObject[i].VelRow
		gameObject[i].col += gameObject[i].VelCol
	}
	if CollidesWithWall(ball) {
		ball.VelRow = -ball.VelRow
	}
	if CollidesWithPaddle(ball, Player1Paddle) || CollidesWithPaddle(ball, Player2Paddle) {
		ball.VelCol = -ball.VelCol
	}
}

func CollidesWithPaddle(obj *GameObject, paddle *GameObject) bool {
	var collidesonCol bool
	if ball.col < paddle.col {
		collidesonCol = ball.col+ball.VelCol >= paddle.col
	} else {
		collidesonCol = ball.col+ball.VelCol <= paddle.col
	}
	return collidesonCol &&
		ball.row >= paddle.row &&
		ball.row < paddle.row+paddle.height
}

func HandleUserInput(key string) {
	_, screenHeight := screen.Size()
	if key == "Rune[q]" {
		screen.Fini()
		os.Exit(0)
	} else if key == "Rune[w]" && Player1Paddle.row > 0 {
		Player1Paddle.row--
	} else if key == "Rune[s]" && Player1Paddle.row+Player1Paddle.height < screenHeight {
		Player1Paddle.row++
	} else if key == "Up" && Player2Paddle.row > 0 {
		Player2Paddle.row--
	} else if key == "Down" && Player2Paddle.row+Player2Paddle.height < screenHeight {
		Player2Paddle.row++
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
	Player1Paddle = &GameObject{
		row:    paddleStart,
		col:    0,
		width:  1,
		height: PaddleHeight,
		symbol: PaddleSymbol,
		VelRow: 0,
		VelCol: 0,
	}
	Player2Paddle = &GameObject{
		row:    paddleStart,
		col:    width - 1,
		width:  1,
		height: PaddleHeight,
		symbol: PaddleSymbol,
		VelRow: 0,
		VelCol: 0,
	}
	ball = &GameObject{
		row:    height / 2,
		col:    width / 2,
		width:  1,
		height: 1,
		symbol: BallSymbol,
		VelRow: BallVelRow,
		VelCol: BallVelCol,
	}
	gameObject = []*GameObject{
		Player1Paddle, Player2Paddle, ball,
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

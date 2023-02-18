package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

// Отрисовать обьекты на экране
// Отрисовать обьект для игрока
// мувмент
// Границы
// Отрисовать шар
// мувмент шара
// Столкновение
// Game over

func PrintString(s tcell.Screen, row, col int, str string) {
	for _, c := range str {
		s.SetContent(col, row, c, nil, tcell.StyleDefault)
		col += 1
	}
}

func Print(s tcell.Screen, row, col, width, height int, ch rune) {

	for r := 0; r < height; r++ {
		for c := 0; c < width; c++ {
			s.SetContent(col+c, row+r, ch, nil, tcell.StyleDefault)
		}

	}

}

func displayHelloWorld(screen tcell.Screen) {
	screen.Clear()
	// PrintString(screen, 2, 15, "Hello, World!")
	Print(screen, 0, 0, 5, 5, '*')
	screen.Show()
}

// This program just prints "Hello, World!".  Press ESC to exit.
func main() {
	screen := Initscreen()
	displayHelloWorld(screen)

	for {
		switch ev := screen.PollEvent().(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				screen.Fini()
				os.Exit(0)
			}
		}
	}
}

func Initscreen() tcell.Screen {
	screen, err := tcell.NewScreen()
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

	return screen
}

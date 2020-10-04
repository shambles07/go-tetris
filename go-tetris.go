package main

import (
	"fmt"
	"github.com/shambles07/go-tetris/tetris"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	game := tetris.NewGame()
	game.Start()

	termbox.Close()
	fmt.Println("Bye!")
}

package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"github.com/shambles07/go-tetris"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	game := tetris.NewGame()
	game.Start()

	termbox.Close()
	fmt.Println("Bye!")
}

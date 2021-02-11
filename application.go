package main

import (
	"errors"
	"fmt"
	"github.com/eiannone/keyboard"
	"math/rand"
	"os/exec"
	"time"
)

type Snake struct {
	X int
	Y int
	Next *Snake
}

type Game struct {
	head *Snake
	tail *Snake
	move int
	area [][]int
	fruit []int
	last *Snake
	H int
	W int
}

var GameOver = errors.New("game over")

func Initialise(H int, W int) Game {
	head := &Snake{
		X:    3,
		Y:    3,
		Next: nil,
	}
	tail := &Snake{
		X:    3,
		Y:    1,
		Next: &Snake{
			X:    3,
			Y:    2,
			Next: head,
		},
	}
	area := make([][]int, H)
	for i := range area {
		area[i] = make([]int, W)
	}
	fruit := []int{H/2,W/2}
	return Game{
		head: head,
		tail: tail,
		area:  area,
		fruit: fruit,
		last: tail,
		move: 0,
		H: H,
		W: W,
	}
}

func (g *Game) loop() {
	i := 0
	for {
		g.reset()
		if err := g.process(); err != nil {
			fmt.Println(err.Error())
			return
		}
		g.render()
		g.update()
		i++
		<- time.After(time.Millisecond*500)
	}
}

func (g *Game) render() {
	// render area with snake and fruit
	for i := range g.area {
		for j := range g.area[i] {
			if g.area[i][j] == 0 {
				fmt.Print(" ")
			} else if g.area[i][j] == 1 {
				fmt.Print("x")
			} else if g.area[i][j] == 2 {
				fmt.Print("o")
			}
		}
		fmt.Println()
	}
}

func (g *Game) reset() {
	// clear screen
	for i := 0; i < len(g.area); i++ {
		for j := 0; j < len(g.area[i]); j++ {
			if g.area[i][j] == 2 {
				continue
			}
			if i == 0 || i == len(g.area)-1 || j == 0 || j == len(g.area[i])-1 {
				g.area[i][j] = 1
			} else {
				g.area[i][j] = 0
			}

		}
	}
}

func (g *Game) process() error {
	// process area 2d array with updated values
	node := g.tail
	for node != nil {
		if g.area[node.X][node.Y] == 1 {
			return GameOver
		}
		if g.area[node.X][node.Y] == 2 {
			g.tail = g.last
			g.fruit = []int{rand.Intn(g.H-3)+1, rand.Intn(g.W-3)+1}
		}
		g.area[node.X][node.Y] = 1
		node = node.Next
	}
	g.area[g.fruit[0]][g.fruit[1]] = 2
	return nil
}

func (g *Game) update() {
	head := &Snake{}
	head.X = g.head.X
	head.Y = g.head.Y
	switch g.move {
	case 0: head.Y += 1
	case 1: head.X -= 1
	case 2: head.Y -= 1
	case 3: head.X += 1
	}
	g.last = g.tail
	g.tail = g.tail.Next
	g.head.Next = head
	g.head = head
}

func inputListener(g *Game) {
	go func(g *Game) {
		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer func() {
			_ = keyboard.Close()
		}()
		for {
			_, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}
			switch key {
			case keyboard.KeyArrowRight: g.move = 0
			case keyboard.KeyArrowUp: g.move = 1
			case keyboard.KeyArrowLeft: g.move = 2
			case keyboard.KeyArrowDown: g.move = 3
			}
		}
	}(g)
}

func main() {
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	g := Initialise(25, 50)
	inputListener(&g)
	g.loop()
}
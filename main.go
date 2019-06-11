package main

import (
	"./lichess"
	"fmt"
	"github.com/notnil/chess"
	"github.com/notnil/chessimg"
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	game, err := lichess.GetGame("https://lichess.org/game/export/bR4b8jno")
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for i, pos := range game.Positions() {
		wg.Add(1)
		go drawPNG(pos, "pos_"+fmt.Sprintf("%03d", i), &wg)
	}

	wg.Wait()

	log.Println(game.Position().Board().Draw())
}

func drawPNG(pos *chess.Position, filebase string, wg *sync.WaitGroup) {
	defer wg.Done()

	// create file
	f, err := os.Create(filebase + ".svg")
	if err != nil {
		log.Fatal(err)
	}

	// write board SVG to file
	if err := chessimg.SVG(f, pos.Board()); err != nil {
		log.Fatal(err)
	}

	f.Close()

	if r := exec.Command("inkscape", "-z", "-e", filebase+".png", filebase+".svg").Run(); r != nil {
		log.Fatal(err)
	}

	os.Remove(filebase + ".svg")
}

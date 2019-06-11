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
	game, gameID, err := lichess.GetGame("https://lichess.org/bR4b8jno") // 5 moves game
	//game, err := lichess.GetGame("https://lichess.org/oyJ7H81yImOI") // 98 moves game
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for i, pos := range game.Positions() {
		wg.Add(1)
		go drawPNG(pos, fileBaseFor(gameID, i), &wg)
	}

	wg.Wait()
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

	// close svg file
	f.Close()

	// Use inkscape to convert svg -> png
	if r := exec.Command("inkscape", "-z", "-e", filebase+".png", filebase+".svg").Run(); r != nil {
		log.Fatal(err)
	}

	// remove temp svg file
	os.Remove(filebase + ".svg")
}

func fileBaseFor(gameID string, i int) string {
	return gameID + fmt.Sprintf("%03d", i)
}

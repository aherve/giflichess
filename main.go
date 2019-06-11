package main

import (
	"./lichess"
	"github.com/notnil/chess"
	"github.com/notnil/chessimg"
	"log"
	"os"
	"strconv"
)

func main() {
	game, err := lichess.GetGame("https://lichess.org/game/export/bR4b8jno")
	if err != nil {
		log.Fatal(err)
	}

	for i, pos := range game.Positions() {
		drawPosition(pos, "pos_"+strconv.Itoa(i)+".svg")
	}

	log.Println(game.Position().Board().Draw())
}

func drawPosition(pos *chess.Position, filename string) {
	// create file
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// write board SVG to file
	if err := chessimg.SVG(f, pos.Board()); err != nil {
		log.Fatal(err)
	}
}

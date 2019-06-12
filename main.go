package main

import (
	"./gifmaker"
	"./lichess"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("Please provide a lichess game url as parameter (example: https://lichess.org/bR4b8jno )")
	}
	game, gameID, err := lichess.GetGame(os.Args[1])
	outputFile := "out/" + gameID + ".gif"
	handle(err)
	gifmaker.GenerateGIF(game, gameID, outputFile)
	log.Println("gif successfully outputed to ", outputFile)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

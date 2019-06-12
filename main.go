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
	handle(err)

	outputFile := "out/" + gameID + ".gif"
	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0755)
	handle(err)
	defer f.Close()

	gifmaker.GenerateGIF(game, gameID, f)
	log.Println("gif successfully outputed to ", outputFile)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

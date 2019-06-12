package main

import (
	"./lichess"
	"fmt"
	"github.com/notnil/chess"
	"github.com/notnil/chessimg"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("Please provide a lichess game url as parameter (example: https://lichess.org/bR4b8jno )")
	}
	game, gameID, err := lichess.GetGame(os.Args[1]) // 5 moves game
	handle(err)

	// Generate PNGs
	var wg sync.WaitGroup
	for i, pos := range game.Positions() {
		wg.Add(1)
		go drawPNG(pos, fileBaseFor(gameID, i), &wg)
		defer cleanup(gameID, i)
	}
	wg.Wait()

	// Generate GIF
	outGIF := &gif.GIF{}
	for i, _ := range game.Positions() {
		f, err := os.Open(fileBaseFor(gameID, i) + ".png")
		handle(err)
		inPNG, err := png.Decode(f)
		handle(err)
		f.Close()

		bounds := inPNG.Bounds()
		palettedImage := image.NewPaletted(bounds, palette.Plan9)
		draw.Draw(palettedImage, palettedImage.Rect, inPNG, bounds.Min, draw.Over)

		// Add new frame to animated GIF
		outGIF.Image = append(outGIF.Image, palettedImage)
		if i == len(game.Positions())-1 {
			outGIF.Delay = append(outGIF.Delay, 450)
		} else if i < 10 {
			outGIF.Delay = append(outGIF.Delay, 50)
		} else {
			outGIF.Delay = append(outGIF.Delay, 120)
		}

	}

	f, err := os.OpenFile(outputFile(gameID), os.O_WRONLY|os.O_CREATE, 0755)
	handle(err)
	defer f.Close()
	gif.EncodeAll(f, outGIF)
	log.Println("gif successfully outputed to ", outputFile(gameID))
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func cleanup(gameID string, i int) {
	fileBase := fileBaseFor(gameID, i)
	os.Remove(fileBase + ".svg")
	os.Remove(fileBase + ".png")
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
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		if r := exec.Command("inkscape", "-z", "-e", filebase+".png", filebase+".svg").Run(); r != nil {
			log.Fatal(err)
		}
		wg.Done()
	}(wg)
}

func outputFile(gameID string) string {
	return "out/" + gameID + ".gif"
}

func fileBaseFor(gameID string, i int) string {
	return "/tmp/" + gameID + fmt.Sprintf("%03d", i)
}

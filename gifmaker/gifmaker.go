package gifmaker

import (
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

func GenerateGIF(game *chess.Game, gameID string, outputFile string) {

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
		defer f.Close()
		inPNG, err := png.Decode(f)
		handle(err)

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

	f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0755)
	handle(err)
	defer f.Close()
	gif.EncodeAll(f, outGIF)
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
}

func cleanup(gameID string, i int) {
	fileBase := fileBaseFor(gameID, i)
	os.Remove(fileBase + ".svg")
	os.Remove(fileBase + ".png")
}

func fileBaseFor(gameID string, i int) string {
	return "/tmp/" + gameID + fmt.Sprintf("%03d", i)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

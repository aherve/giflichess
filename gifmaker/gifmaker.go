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
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
)

type imgOutput struct {
	index int
	img   *image.Paletted
}

// GenerateGIF will use *chess.Game to write a gif into an io.Writer
// This uses inkscape as a dependency
func GenerateGIF(game *chess.Game, gameID string, out io.Writer) {

	// Generate PNGs
	var wg sync.WaitGroup
	for i, pos := range game.Positions() {
		wg.Add(1)
		go drawPNG(pos, fileBaseFor(gameID, i), &wg)
		defer cleanup(gameID, i)
	}
	wg.Wait()

	// Generate atomic GIFs
	images := make([]*image.Paletted, len(game.Positions()), len(game.Positions()))
	imgChan := make(chan imgOutput)
	quit := make(chan bool)
	for i, _ := range game.Positions() {
		wg.Add(1)
		go func(gameID string, i int, outChan chan imgOutput) {
			defer wg.Done()
			outChan <- imgOutput{i, encodeGIFImage(gameID, i)}
		}(gameID, i, imgChan)

	}
	go func() {
		wg.Wait()
		quit <- true
	}()

loop:
	for {
		select {
		case res := <-imgChan:
			images[res.index] = res.img
		case <-quit:
			break loop
		}
	}

	// Generate final GIF
	outGIF := &gif.GIF{}
	for i, img := range images {

		// Add new frame to animated GIF
		outGIF.Image = append(outGIF.Image, img)
		if i == len(game.Positions())-1 {
			outGIF.Delay = append(outGIF.Delay, 450)
		} else if i < 10 {
			outGIF.Delay = append(outGIF.Delay, 50)
		} else {
			outGIF.Delay = append(outGIF.Delay, 120)
		}
	}

	gif.EncodeAll(out, outGIF)
}

// encodeGIFImage reads a png from gameID & index, and returns a palettedImage
func encodeGIFImage(gameID string, i int) *image.Paletted {
	f, err := os.Open(fileBaseFor(gameID, i) + ".png")
	handle(err)
	defer f.Close()
	inPNG, err := png.Decode(f)
	handle(err)

	bounds := inPNG.Bounds()
	palettedImage := image.NewPaletted(bounds, palette.Plan9)
	draw.Draw(palettedImage, palettedImage.Rect, inPNG, bounds.Min, draw.Over)

	return palettedImage
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

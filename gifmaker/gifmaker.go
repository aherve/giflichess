package gifmaker

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/aherve/chessimg"
	"github.com/aherve/gopool"
	"github.com/notnil/chess"
)

type imgOutput struct {
	index int
	img   *image.Paletted
	err   error
}

func whiteName(game *chess.Game) string {
	var elo, name string
	for _, tag := range game.TagPairs() {
		if strings.ToLower(tag.Key) == "white" {
			name = tag.Value
		}
		if strings.ToLower(tag.Key) == "whiteelo" {
			elo = tag.Value
		}
		if len(elo) > 0 && len(name) > 0 {
			break
		}
	}
	if len(name) > 0 && len(elo) > 0 {
		return fmt.Sprintf("%s (%s)", name, elo)
	} else if len(name) > 0 {
		return name
	}
	return "unknown"
}

func blackName(game *chess.Game) string {
	var elo, name string
	for _, tag := range game.TagPairs() {
		if strings.ToLower(tag.Key) == "black" {
			name = tag.Value
		}
		if strings.ToLower(tag.Key) == "blackelo" {
			elo = tag.Value
		}
		if len(elo) > 0 && len(name) > 0 {
			break
		}
	}
	if len(name) > 0 && len(elo) > 0 {
		return fmt.Sprintf("%s (%s)", name, elo)
	} else if len(name) > 0 {
		return name
	}
	return "unknown"
}

// GenerateGIF will use *chess.Game to write a gif into an io.Writer
// This uses inkscape & imagemagick as a dependency
func GenerateGIF(game *chess.Game, gameID string, reversed bool, out io.Writer, maxConcurrency int) error {

	// Generate PNGs
	pool := gopool.NewPool(maxConcurrency)
	for i, pos := range game.Positions() {
		pool.Add(1)
		go drawPNG(pos, whiteName(game), blackName(game), reversed, fileBaseFor(gameID, i), pool)
		defer cleanup(gameID, i)
	}
	pool.Wait()

	// add Result to last png image
	fileName := fileBaseFor(gameID, len(game.Positions())-1) + ".png"
	cmd := exec.Command("gifmaker/addResult.sh", fileName, game.Outcome().String())
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Generate atomic GIFs
	images := make([]*image.Paletted, len(game.Positions()), len(game.Positions()))
	imgChan := make(chan imgOutput)
	quit := make(chan bool)

	go func() {
		for i := range game.Positions() {
			pool.Add(1)
			go func(gameID string, i int, outChan chan imgOutput, pool *gopool.GoPool) {
				defer pool.Done()
				encoded, err := encodeGIFImage(gameID, i)
				outChan <- imgOutput{i, encoded, err}
				return
			}(gameID, i, imgChan, pool)

		}
		pool.Wait()
		quit <- true
	}()

loop:
	for {
		select {
		case res := <-imgChan:
			if res.err != nil {
				return res.err
			}
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
	return nil
}

// encodeGIFImage reads a png from gameID & index, and returns a palettedImage
func encodeGIFImage(gameID string, i int) (*image.Paletted, error) {
	f, err := os.Open(fileBaseFor(gameID, i) + ".png")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	inPNG, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := inPNG.Bounds()
	palettedImage := image.NewPaletted(bounds, palette.Plan9)
	draw.Draw(palettedImage, palettedImage.Rect, inPNG, bounds.Min, draw.Over)

	return palettedImage, nil
}

func drawPNG(pos *chess.Position, whiteName string, blackName string, reversed bool, filebase string, pool *gopool.GoPool) error {
	defer pool.Done()

	// create file
	f, err := os.Create(filebase + ".svg")
	if err != nil {
		return err
	}

	// write board SVG to file
	if reversed {
		err = chessimg.ReversedSVG(f, pos.Board())
	} else {
		err = chessimg.SVG(f, pos.Board())
	}
	if err != nil {
		return err
	}

	// close svg file
	f.Close()

	// Use inkscape to convert svg -> png
	cmd := exec.Command("inkscape", "-z", "-e", filebase+".png", filebase+".svg")
	cmd.Stderr = os.Stderr
	if r := cmd.Run(); r != nil {
		return err
	}

	// add annotation
	cmd = exec.Command("gifmaker/annotate.sh", filebase+".png", whiteName, blackName)
	cmd.Stderr = os.Stderr
	if r := cmd.Run(); r != nil {
		return err
	}
	return nil
}

func cleanup(gameID string, i int) {
	fileBase := fileBaseFor(gameID, i)
	os.Remove(fileBase + ".svg")
	os.Remove(fileBase + ".png")
}

func fileBaseFor(gameID string, i int) string {
	return "/tmp/" + gameID + fmt.Sprintf("%03d", i)
}

package main

import (
	"./gifmaker"
	"./lichess"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	var output string
	var input string
	app := cli.NewApp()
	app.Name = "giflichess"
	app.Usage = "generate fancy gifs from your lichess games"
	app.Description = "giflichess can turn a lichess game id into an animated gif. You can either use this program as a cli tool, or as a web server"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate a gif and output to a file",
			Action: func(c *cli.Context) error {
				if len(input) == 0 {
					return fmt.Errorf("Please pass an input game: example --game https://lichess.org/bR4b8jno")
				}
				return GenerateFile(input, output)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "game, g",
					Value:       "",
					Usage:       "lichess game url, or lichess game ID",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "ouput, o",
					Value:       "out.gif",
					Usage:       "output file path",
					Destination: &output,
				},
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "run as a server",
			Action: func(c *cli.Context) error {
				fmt.Println("server")
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port, p",
					Value: 8080,
					Usage: "server port",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateFile(urlOrID string, outFile string) error {
	fmt.Printf("generating file %s from game %s...\n", outFile, urlOrID)
	game, gameID, err := lichess.GetGame(urlOrID)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	gifmaker.GenerateGIF(game, gameID, f)
	fmt.Printf("gif successfully outputed to %s\n", outFile)
	return nil
}

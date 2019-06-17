package main

import (
	"fmt"
	"github.com/aherve/giflichess/lichess"
	"github.com/aherve/giflichess/server"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	var output string
	var input string
	var port int
	var reversed bool
	app := cli.NewApp()
	app.Name = "giflichess"
	app.Usage = "generate fancy gifs from your lichess games"
	app.Description = "giflichess can turn a lichess game id into an animated gif. You can either use this program as a cli tool, or as a web server"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate a gif and output to a file",
			Action: func(c *cli.Context) error {
				if len(input) == 0 {
					return fmt.Errorf("Please pass an input game: example --game https://lichess.org/bR4b8jno")
				}
				return lichess.GenerateFile(input, reversed, output)
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "game, g",
					Value:       "",
					Usage:       "lichess game url, or lichess game ID",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output, o",
					Value:       "out.gif",
					Usage:       "output file path",
					Destination: &output,
				},
				cli.BoolFlag{
					Name:        "reversed, r",
					Usage:       "Flip board",
					Destination: &reversed,
				},
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "run as a server",
			Action: func(c *cli.Context) error {
				server.Serve(port)
				return nil
			},
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "port, p",
					Value:       8080,
					Usage:       "server port",
					Destination: &port,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

const Version = "1.1.2"

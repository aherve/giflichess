package main

import (
	"errors"
	"fmt"
	"github.com/aherve/giflichess/gifmaker"
	"github.com/aherve/giflichess/lichess"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	var output string
	var input string
	var port int
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
				return generateFile(input, output)
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
				serve(port)
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

func serve(port int) {
	http.HandleFunc("/", gifHandler)
	log.Println("starting server on port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func gifHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var status int
	log := func() {
		log.Println(r.Method, r.URL, status, time.Since(start))
	}
	defer log()

	// get ID
	maybeID, err := getIDFromQuery(r)
	if err != nil {
		status = 400
		http.Error(w, err.Error(), status)
		return
	}

	// get game
	game, gameID, err := lichess.GetGame(maybeID)
	if err != nil {
		status = 500
		http.Error(w, err.Error(), status)
		return
	}

	// write gif
	w.Header().Set("Content-Disposition", "attachment")
	w.Header().Set("filename", gameID+".gif")
	err = gifmaker.GenerateGIF(game, gameID, w)
	if err != nil {
		status = 500
		http.Error(w, err.Error(), status)
		return
	}
	status = 200
}

func generateFile(urlOrID string, outFile string) error {
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

func getIDFromQuery(r *http.Request) (string, error) {
	split := strings.Split(r.URL.Path, "/")
	if len(split) < 2 {
		return "", errors.New("could not find no id")
	}
	return split[1], nil
}

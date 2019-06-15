package server

import (
	"errors"
	"fmt"
	"github.com/aherve/giflichess/gifmaker"
	"github.com/aherve/giflichess/lichess"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Serve(port int) {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/lichess/", lichessGifHandler)
	http.HandleFunc("/", rootHandler)
	log.Println("starting server on port", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte("<html><head></head><body><h1>Hello</h1> <p>visit /lichess/:id to get a lichess game</p></body></html>"))
	log := func() {
		log.Println(r.Method, r.URL, 200, time.Since(start))
	}
	defer log()
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte("{\"ping\": \"pong\"}"))
	log := func() {
		log.Println(r.Method, r.URL, 200, time.Since(start))
	}
	defer log()
}

func lichessGifHandler(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Cache-Control", "no-cache")
		http.Error(w, err.Error(), status)
		return
	}

	// get game
	game, gameID, err := lichess.GetGame(maybeID)
	if err != nil {
		status = 500
		w.Header().Set("Cache-Control", "no-cache")
		http.Error(w, err.Error(), status)
		return
	}

	// write gif
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.gif\"", gameID))
	w.Header().Set("filename", gameID+".gif")
	w.Header().Set("Cache-Control", cacheControl(1296000))
	err = gifmaker.GenerateGIF(game, gameID, w)
	if err != nil {
		status = 500
		w.Header().Set("Cache-Control", "no-cache")
		http.Error(w, err.Error(), status)
		return
	}
	status = 200
}

func getIDFromQuery(r *http.Request) (string, error) {
	split := strings.Split(r.URL.Path, "/")
	if len(split) < 3 || len(split[2]) < 8 {
		return "", errors.New("No id provided. Please visit /some-id. Example: /bR4b8jno")
	}
	return split[2], nil
}

func cacheControl(seconds int) string {
	return fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds)
}

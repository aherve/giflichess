package server

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aherve/giflichess/gifmaker"
	"github.com/aherve/giflichess/lichess"
	"github.com/notnil/chess"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Serve(port int) {
	http.HandleFunc("/api/pgn", pgnGifHandler)
	http.HandleFunc("/api/ping", pingHandler)
	http.HandleFunc("/api/lichess/", lichessGifHandler)
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/api", rootHandler)
	log.Printf("starting %s server on port %v\n", env(), port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte("<html><head></head><body><h1>Hello</h1> <p>visit /api/lichess/:id to get a lichess game</p></body></html>"))
	logReq := func() {
		log.Println(r.Method, r.URL, 200, time.Since(start))
	}
	defer logReq()
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write([]byte("{\"ping\": \"pong\"}"))
	logReq := func() {
		log.Println(r.Method, r.URL, 200, time.Since(start))
	}
	defer logReq()
}

func lichessGifHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var status int
	logReq := func() {
		log.Println(r.Method, r.URL, status, time.Since(start))
	}
	defer logReq()

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
	htmlGifWriter(gameID, game, w, r, &status)
}

func htmlGifWriter(gameID string, game *chess.Game, w http.ResponseWriter, r *http.Request, status *int) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.gif\"", gameID))
	w.Header().Set("filename", gameID+".gif")
	if env() == "production" {
		w.Header().Set("Cache-Control", cacheControl(1296000))
	} else {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}
	err := gifmaker.GenerateGIF(game, gameID, getReversed(r), w)
	if err != nil {
		w.Header().Set("Cache-Control", "no-cache")
		http.Error(w, err.Error(), *status)
		return
	}
	*status = 200
}

func pgnGifHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var status int
	logReq := func() {
		log.Println(r.Method, r.URL, status, time.Since(start))
	}
	defer logReq()
	if r.Method != "POST" {
		status = http.StatusMethodNotAllowed
		http.Error(w, "Invalid request method", status)
	}

	// get body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status = 400
		w.Header().Set("Cache-Control", "no-cache")
		http.Error(w, err.Error(), status)
		return
	}

	// create game
	log.Println(string(body))
	pgn, err := chess.PGN(r.Body)
	if err != nil {
		status = 400
		w.Header().Set("Cache-Control", "no-cache")
		http.Error(w, err.Error(), status)
		return
	}
	game := chess.NewGame(pgn)
	log.Println(game.Positions())
	hash := md5.Sum(body)
	gameID := hex.EncodeToString(hash[:])

	// write gif
	htmlGifWriter(gameID, game, w, r, &status)
}

func getReversed(r *http.Request) bool {
	if s, ok := r.URL.Query()["reversed"]; ok && len(s) == 1 {
		return s[0] == "true"
	}
	return false
}

func getIDFromQuery(r *http.Request) (string, error) {
	split := strings.Split(r.URL.Path, "/")
	if len(split) < 4 || len(split[3]) < 8 {
		return "", errors.New("No id provided. Please visit /some-id. Example: /bR4b8jno")
	}
	return split[3], nil
}

func cacheControl(seconds int) string {
	return fmt.Sprintf("max-age=%d, public, must-revalidate, proxy-revalidate", seconds)
}

func env() string {
	fromEnv := os.Getenv("APP_ENV")
	if len(fromEnv) > 0 {
		return fromEnv
	}
	return "development"
}

package lichess

import (
	"fmt"
	"github.com/notnil/chess"
	"net/http"
	"net/url"
	"strings"
)

// GetPGN extracts the PGN from a lichess game url
func GetGame(path string) (*chess.Game, error) {
	resp, err := http.Get("https://lichess.org/game/export/bR4b8jno")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	pgn, err := chess.PGN(resp.Body)
	if err != nil {
		return nil, err
	}

	return chess.NewGame(pgn), nil
}

// gameID extracts the id of a lichess game from either analyze or game url
func gameID(path string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	id := strings.Split(u.Path, "/")[1]
	if len(id) < 8 {
		return "", fmt.Errorf("could not find id from string %s", path)
	}
	return id[0:8], nil
}

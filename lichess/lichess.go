package lichess

import (
	"fmt"
	"github.com/aherve/giflichess/gifmaker"
	"github.com/notnil/chess"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func GenerateFile(urlOrID string, outFile string) error {
	fmt.Printf("generating file %s from game %s...\n", outFile, urlOrID)
	game, gameID, err := GetGame(urlOrID)
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

// GetPGN extracts the PGN from a lichess game url
func GetGame(pathOrID string) (*chess.Game, string, error) {
	id, err := gameID(pathOrID)
	if err != nil {
		return nil, "", err
	}
	resp, err := http.Get("https://lichess.org/game/export/" + id)
	if err != nil {
		return nil, id, err
	}

	defer resp.Body.Close()

	pgn, err := chess.PGN(resp.Body)
	if err != nil {
		return nil, id, err
	}

	return chess.NewGame(pgn), id, nil
}

// gameID extracts the id of a lichess game from either analyze url, game url, or id
func gameID(pathOrID string) (string, error) {

	matchId, err := regexp.MatchString(`^[a-zA-Z0-9]{8,}$`, pathOrID)
	if err != nil {
		return "", err
	}

	if matchId {
		return pathOrID[0:8], nil
	}

	u, err := url.Parse(pathOrID)
	if err != nil {
		return "", err
	}

	split := strings.Split(u.Path, "/")
	if len(split) < 2 || len(split[1]) < 8 {
		return "", fmt.Errorf("could not find id from string \"%s\"", pathOrID)
	}
	return split[1][0:8], nil
}

package gifmaker

import (
	"github.com/notnil/chess"
	"os"
	"testing"
)

func TestWhiteName(t *testing.T) {
	file, err := os.Open("../datatest/bR4b8jno.pgn")
	if err != nil {
		t.Error(err)
	}
	readPGN, err := chess.PGN(file)
	if err != nil {
		t.Error(err)
	}

	game := chess.NewGame(readPGN)

	if white := whiteName(game); white != "aherve (1527)" {
		t.Error("expected white to == aherve (1527), got ", white)
	}
	if black := blackName(game); black != "minahabibzadeeh (1558)" {
		t.Error("expected black to == minahabibzadeeh (1558), got ", black)
	}
}

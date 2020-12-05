package gifmaker

import (
	"os"
	"testing"

	"github.com/notnil/chess"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "aherve (1527)", whiteName(game))
	assert.Equal(t, "minahabibzadeeh (1558)", blackName(game))

}

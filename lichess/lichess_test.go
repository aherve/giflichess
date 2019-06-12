package lichess

import "testing"

func TestGetGameID(t *testing.T) {
	gameURL := "https://lichess.org/bR4b8jnouzUP"
	analyzeURL := "https://lichess.org/bR4b8jno/white"
	longID := "bR4b8jnouzUP"
	shortID := "bR4b8jno"

	expectedID := "bR4b8jno"

	if r, _ := gameID(gameURL); r != expectedID {
		t.Error("expected", expectedID, ", got ", r)
	}

	if r, _ := gameID(analyzeURL); r != expectedID {
		t.Error("expected", expectedID, ", got ", r)
	}

	if r, _ := gameID(longID); r != expectedID {
		t.Error("expected", expectedID, ", got ", r)
	}

	if r, _ := gameID(shortID); r != expectedID {
		t.Error("expected", expectedID, ", got ", r)
	}
}

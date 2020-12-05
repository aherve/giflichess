package lichess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGameID(t *testing.T) {
	urls := []string{
		"https://lichess.org/bR4b8jnouzUP",
		"https://lichess.org/bR4b8jno/white",
		"bR4b8jnouzUP",
		"bR4b8jnox",
	}

	expectedID := "bR4b8jno"

	for _, url := range urls {
		res, err := gameID(url)
		assert.Nil(t, err, "err should be nil")
		assert.Equal(t, expectedID, res)
	}

}

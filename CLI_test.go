package poker_test

import (
	"strings"
	"testing"

	poker "github.com/fath-nasrudin/goleague"
)

func TestCLI(t *testing.T) {
	t.Run("Record Chris wins from user input", func(t *testing.T) {
		playerStore := &poker.StubPlayerStore{}
		in := strings.NewReader("Chris wins\n")

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()
		
		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("Record Anto wins from user input", func(t *testing.T) {
		playerStore := &poker.StubPlayerStore{}
		in := strings.NewReader("Anto wins\n")
		
		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Anto")
	})
}

package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("Record Chris wins from user input", func(t *testing.T) {
		playerStore := &StubPlayerStore{}
		in := strings.NewReader("Chris wins\n")

		cli := &CLI{playerStore, in}
		cli.PlayPoker()
		
		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("Record Anto wins from user input", func(t *testing.T) {
		playerStore := &StubPlayerStore{}
		in := strings.NewReader("Anto wins\n")
		
		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Anto")
	})
}

func assertPlayerWin(t testing.TB, store *StubPlayerStore, wantWinner string) {
	t.Helper()

	got := store.winCalls[0]

	if len(store.winCalls) != 1 {
		t.Fatalf("want %d calls but got %d", 1, len(store.winCalls))
	}

	if got != wantWinner {
		t.Fatalf("want %q as winner but got %q", wantWinner, got)
	}
}
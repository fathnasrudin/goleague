package poker

import (
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	playerStore := &StubPlayerStore{}
	in := strings.NewReader("Chris wins\n")
	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	got := playerStore.winCalls[0]
	want := "Chris"

	if len(playerStore.winCalls) != 1 {
		t.Fatalf("want %d calls but got %d", 1, len(playerStore.winCalls))
	}

	if got != want {
		t.Fatalf("want %q but got %q", want, got)
	}
}
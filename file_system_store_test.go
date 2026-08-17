package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("Get League with double read", func(t *testing.T) {
		database := strings.NewReader(`[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}
		]`)
		store := FileSystemPlayerStore{database}
		got := store.GetLeague()
		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, want, got)

		got = store.GetLeague()
		assertLeague(t, want, got)
	})

	t.Run("Get player score", func(t *testing.T) {
		database := strings.NewReader(`[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}
		]`)
		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Chris")
		want := 33

		assertScoreEquals(t, want, got)
	})
}
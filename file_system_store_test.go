package main

import (
	"io"
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("Get League from a reader", func(t *testing.T) {
		database, removeFile := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}
		]`)
		defer removeFile()

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
		database, removeFile := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}
		]`)
		defer removeFile()
		
		store := FileSystemPlayerStore{database}

		got := store.GetPlayerScore("Chris")
		want := 33

		assertScoreEquals(t, want, got)
	})

	t.Run("store wins for existing player", func(t *testing.T) {
		database, removeFile := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
			]`)
		defer removeFile()
		
		store := FileSystemPlayerStore{database}

		store.ProcessWin("Chris")
		got := store.GetPlayerScore("Chris")
		want := 34

		assertScoreEquals(t, want, got)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, removeFile := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
			]`)
		defer removeFile()
		
		store := FileSystemPlayerStore{database}

		store.ProcessWin("Anto")
		got := store.GetPlayerScore("Anto")
		want := 1
		assertScoreEquals(t, want, got)

		store.ProcessWin("Anto")
		assertScoreEquals(t, 2, store.GetPlayerScore("Anto"))
	})
}

func createTempFile(t testing.TB, initialData string) (file io.ReadWriteSeeker, removeFile func()) {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "db")

	if err != nil {
		t.Fatalf("Could not create a temporary file, %v", err)
	}

	tmpFile.Write([]byte(initialData))

	removeFile = func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
	file = tmpFile

	return
}
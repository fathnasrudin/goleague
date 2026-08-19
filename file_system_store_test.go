package poker

import (
	"os"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	t.Run("Get League from a reader and sorted by highest wins", func(t *testing.T) {
		database, removeFile := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}
		]`)
		defer removeFile()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := store.GetLeague()
		want := League{{Name: "Chris", Wins: 33}, {Name: "Cleo", Wins: 10}}

		AssertLeague(t, want, got)

		got = store.GetLeague()
		AssertLeague(t, want, got)
	})

	t.Run("Get player score", func(t *testing.T) {
		database, removeFile := createTempFile(t, `[
		{"Name": "Cleo", "Wins": 10},
		{"Name": "Chris", "Wins": 33}
		]`)
		defer removeFile()
		
		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := store.GetPlayerScore("Chris")
		want := 33

		AssertScoreEquals(t, want, got)
	})

	t.Run("store wins for existing player", func(t *testing.T) {
		database, removeFile := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
			]`)
		defer removeFile()
		
		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		store.RecordWin("Chris")
		got := store.GetPlayerScore("Chris")
		want := 34

		AssertScoreEquals(t, want, got)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, removeFile := createTempFile(t, `[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
			]`)
		defer removeFile()
		
		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		store.RecordWin("Anto")
		got := store.GetPlayerScore("Anto")
		want := 1
		AssertScoreEquals(t, want, got)

		store.RecordWin("Anto")
		AssertScoreEquals(t, 2, store.GetPlayerScore("Anto"))
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, removeFile := createTempFile(t, ``)
		defer removeFile()
		
		_, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)
	})
}

func createTempFile(t testing.TB, initialData string) (file *os.File, removeFile func()) {
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

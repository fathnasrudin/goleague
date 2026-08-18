package main

import (
	"io"
	"testing"
)

func TestTape_Write(t *testing.T) {
	t.Run("write smaller than old file", func(t *testing.T) {
		// write data
		file, clean := createTempFile(t, "12345")
		defer clean()

		data := &tape{file}
		data.Write([]byte("abc"))

		data.file.Seek(0, io.SeekStart)
		fileContents, _ := io.ReadAll(data.file)
		
		got := string(fileContents)
		want := "abc"

		if got != want {
			t.Errorf("Failed on write, want %q but got %q", want, got)
		}
	})
}
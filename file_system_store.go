package main

import (
	"encoding/json"
	"io"
	"log"
)

type FileSystemPlayerStore struct {
	database io.Reader
}

func (s *FileSystemPlayerStore) GetLeague() []Player {
	var league []Player
	err := json.NewDecoder(s.database).Decode(&league)
	if err != nil {
		log.Fatalf("Failed to parse %v to JSON. %q", s.database, err)
	}
	return league
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	file.Seek(0, io.SeekStart)
	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)

	}

	return &FileSystemPlayerStore{database: json.NewEncoder(&tape{file}), league: league}, nil
}

func (s *FileSystemPlayerStore) GetLeague() League {
	return s.league
}

func (s *FileSystemPlayerStore) GetPlayerScore(name string) int {
	p := s.league.Find(name)
	return p.Wins
}


func (s *FileSystemPlayerStore) RecordWin(name string) {
	p := s.league.Find(name)
	
	if p != nil {
		p.Wins++
	} else {
		s.league = append(s.league, Player{name, 1})
	}
	
	s.database.Encode(s.league)
}
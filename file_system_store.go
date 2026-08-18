package main

import (
	"encoding/json"
	"io"
	"os"
)

type FileSystemPlayerStore struct {
	database io.Writer
	league League
}

func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
	file.Seek(0, io.SeekStart)
	league, _ := NewLeague(file)

	return &FileSystemPlayerStore{database: &tape{file}, league: league}
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
	
	json.NewEncoder(s.database).Encode(s.league)
}
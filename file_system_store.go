package main

import (
	"encoding/json"
	"io"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league League
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, io.SeekStart)
	league, _ := NewLeague(database)

	return &FileSystemPlayerStore{database, league}
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
	
	s.database.Seek(0, io.SeekStart)
	json.NewEncoder(s.database).Encode(s.league)
}
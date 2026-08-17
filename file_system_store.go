package main

import (
	"encoding/json"
	"io"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
}

func (s *FileSystemPlayerStore) GetLeague() League {
	s.database.Seek(0, io.SeekStart)
	league, _ := NewLeague(s.database)
	return league
}

func (s *FileSystemPlayerStore) GetPlayerScore(name string) int {
	s.database.Seek(0, io.SeekStart)
	league := s.GetLeague()
	p := league.Find(name)
	return p.Wins
}


func (s *FileSystemPlayerStore) ProcessWin(name string) error {
	league := s.GetLeague()
	p := league.Find(name)
	p.Wins++
	
	s.database.Seek(0, io.SeekStart)
	return json.NewEncoder(s.database).Encode(league)
}
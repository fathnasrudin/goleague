package main

import (
	"io"
)

type FileSystemPlayerStore struct {
	database io.ReadSeeker
}

func (s *FileSystemPlayerStore) GetLeague() []Player {
	s.database.Seek(0, io.SeekStart)
	league, _ := NewLeague(s.database)
	return league
}

func (s *FileSystemPlayerStore) GetPlayerScore(name string) int {
	s.database.Seek(0, io.SeekStart)
	league := s.GetLeague()
	
	for _, p := range league {
		if name == p.Name {
			return p.Wins
		}
	}
	return 0
}
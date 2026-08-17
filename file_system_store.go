package main

import (
	"encoding/json"
	"io"
)

type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
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


func (s *FileSystemPlayerStore) ProcessWin(name string) error {
	league := s.GetLeague()

	for i, p := range league {
		if name == p.Name {
			league[i].Wins++
			break
		}
	}

	s.database.Seek(0, io.SeekStart)
	return json.NewEncoder(s.database).Encode(league)
}
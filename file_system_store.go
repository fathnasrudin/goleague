package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initializePlayerDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initializing player db file, %v", err)
	}

	league, err := NewLeague(file)
	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)

	}

	return &FileSystemPlayerStore{database: json.NewEncoder(&tape{file}), league: league}, nil
}

func (s *FileSystemPlayerStore) GetLeague() League {

	// sort with highest wins first
	slices.SortFunc(s.league, func(a, b Player) int {
		return b.Wins - a.Wins
	})

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

func initializePlayerDBFile(file *os.File) error {
	
	file.Seek(0, io.SeekStart)

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	// add empty array if file is empty
	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, io.SeekStart) // set the cursor to the beginning
	}

	return nil
}
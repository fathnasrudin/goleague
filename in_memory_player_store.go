package main

import "sync"

type InMemoryPlayerStore struct {
	mu sync.Mutex
	scores map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{scores: make(map[string]int)}
}

func (s *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return s.scores[name]
}

func (s *InMemoryPlayerStore) RecordWin(name string) {
	s.mu.Lock()

	defer s.mu.Unlock()
	s.scores[name]++
}

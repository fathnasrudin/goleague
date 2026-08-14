package main

import (
	"log"
	"net/http"
)

type InMemoryPlayerStore struct {}

func (s *InMemoryPlayerStore) GetPlayerScore(string) int {
	return 123
}


func (s *InMemoryPlayerStore) RecordWin(string) {}

func main() {
	store := &InMemoryPlayerStore{}
	server := &PlayerServer{store: store }
	log.Fatal(http.ListenAndServe(":5000", server))
}
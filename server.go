package main

import (
	"fmt"
	"net/http"
	"strings"
)

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	// handle post method
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusAccepted)
		return
	} 

	score := p.store.GetPlayerScore(player)

	// handle if player score not exist
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

type PlayerServer struct {
	store PlayerStore
}

type PlayerStore interface {
	GetPlayerScore(name string) int
}

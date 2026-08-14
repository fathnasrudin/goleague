package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{scores: map[string]int{
		"pepper": 20,
		"floyd": 10,
	}}

	server := &PlayerServer{store: store}
	
	t.Run("returns pepper score", func(t *testing.T) {
		request := newGetScoreRequest("pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("floyd")
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "10")
	})
}

func newGetScoreRequest(name string) *http.Request {
	url := fmt.Sprintf("/players/%s", name)
	request, _:= http.NewRequest(http.MethodGet, url, nil)
	return request
}

func assertResponseBody(t testing.TB, got, want string, ) {
	if got != want {
		t.Errorf("want %q but got %q", want, got)
	}
}
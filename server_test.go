package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
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

		assertStatus(t, http.StatusOK, response.Code)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("floyd")
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)

		assertStatus(t, http.StatusOK, response.Code)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("udin")
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)
		want := http.StatusNotFound
		got := response.Code

		assertStatus(t, want, got)
	})
}

func TestStoreWins(t *testing.T) {
	store := &StubPlayerStore{scores: map[string]int{}}
	server := &PlayerServer{store: store}

	t.Run("it records wins when POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, http.StatusAccepted, response.Code)
		
		if len(store.winCalls) != 1 {
			t.Errorf("Record win not working correctly. Want 1 calls but got %d", 
			len(store.winCalls) )
		}

		if store.winCalls[0] != "pepper" {
			t.Errorf("Did not store correct winner. Want %q but got %q", 
			"pepper", store.winCalls[0])
		}
	})

	
}

func newGetScoreRequest(name string) *http.Request {
	url := fmt.Sprintf("/players/%s", name)
	request, _:= http.NewRequest(http.MethodGet, url, nil)
	return request
}

func newPostWinRequest(name string) *http.Request{
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string, ) {
	t.Helper()

	if got != want {
		t.Errorf("want %q but got %q", want, got)
	}
}


func assertStatus(t testing.TB, want, got int, ) {
	t.Helper()
	
	if got != want {
		t.Errorf("want status code %d but got %d", want, got)
	}
}
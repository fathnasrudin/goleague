package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
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
	assertTotalCalls := func(t testing.TB, want, got int) {
		t.Helper()
		if want != got {
			t.Errorf("Record win not working correctly. Want %d calls but got %d", want, got)
		}
	}

	assertCorrectWinner := func(t testing.TB, want, got string) {
		t.Helper()
		if want != got {
			t.Errorf("Did not store correct winner. Want %q but got %q", want, got)
		}
	}

	t.Run("it records wins when POST", func(t *testing.T) {
		store := &StubPlayerStore{scores: map[string]int{}}
		server := &PlayerServer{store: store}
		player := "pepper"

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newPostWinRequest(player))

		assertStatus(t, http.StatusAccepted, response.Code)
		assertTotalCalls(t, 1, len(store.winCalls))
		assertCorrectWinner(t, "pepper", store.winCalls[0])
	})

	t.Run("handle store win concurrently", func(t *testing.T) {
		store := &StubPlayerStore{scores: map[string]int{}}
		server := &PlayerServer{store: store}
		player := "pepper"

		response := httptest.NewRecorder()
		totalReq := 3

		var wg sync.WaitGroup
		wg.Add(totalReq)

		for range totalReq {
			go func(){
				server.ServeHTTP(response, newPostWinRequest(player))
				wg.Done()
			}()
		}

		wg.Wait()

		assertStatus(t, http.StatusAccepted, response.Code)
		assertTotalCalls(t, totalReq, len(store.winCalls))
		assertCorrectWinner(t, player, store.winCalls[0])
	})
}

func TestLeague(t *testing.T) {
	store := &StubPlayerStore{}
	server := &PlayerServer{store: store}

	t.Run("It should return status 200", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/league", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		assertStatus(t, http.StatusOK, res.Code)		
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
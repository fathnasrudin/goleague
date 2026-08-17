package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
	league []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func TestGETPlayers(t *testing.T) {
	store := &StubPlayerStore{scores: map[string]int{
		"pepper": 20,
		"floyd": 10,
	}}

	server := NewPlayerServer(store)
	
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
		server := NewPlayerServer(store)
		player := "pepper"

		response := httptest.NewRecorder()
		server.ServeHTTP(response, newPostWinRequest(player))

		assertStatus(t, http.StatusAccepted, response.Code)
		assertTotalCalls(t, 1, len(store.winCalls))
		assertCorrectWinner(t, "pepper", store.winCalls[0])
	})

	t.Run("handle store win concurrently", func(t *testing.T) {
		store := &StubPlayerStore{scores: map[string]int{}}
		server := NewPlayerServer(store)
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
	wantedLeague := []Player{
		{Name: "John", Wins: 20},
		{Name: "Downey", Wins: 10},
	}
	store := &StubPlayerStore{league: wantedLeague}
	server := NewPlayerServer(store)
	

	t.Run("it returns league table as JSON", func(t *testing.T) {
		req := newGetLeagueRequest()
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		got := getLeagueFromResponse(t, res.Body)

		assertContentType(t, jsonContentType, res)
		assertStatus(t, http.StatusOK, res.Code)		
		assertLeague(t, wantedLeague, got)
	})
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, %v", body, err)
	}
	return
}

func assertLeague(t testing.TB, want, got []Player) {
	t.Helper()
	
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Want %v but got %v", want, got)
	}
}

func assertContentType(t testing.TB, want string, response *httptest.ResponseRecorder) {
	t.Helper()

	got :=  response.Result().Header.Get("content-type")
	if want != got {
		t.Errorf("Response did not have content-type of application/json, got %q", got)
	}
}

func newGetLeagueRequest() *http.Request {
	request, _:=  http.NewRequest(http.MethodGet, "/league", nil)
	return request
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

func assertScoreEquals(t testing.TB, want, got int) {
	t.Helper()

	if got != want {
		t.Errorf("Want %d but got %d", want, got)
	}
}
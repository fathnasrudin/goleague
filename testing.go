package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, wantWinner string) {
	t.Helper()

	got := store.winCalls[0]

	if len(store.winCalls) != 1 {
		t.Fatalf("want %d calls but got %d", 1, len(store.winCalls))
	}

	if got != wantWinner {
		t.Fatalf("want %q as winner but got %q", wantWinner, got)
	}
}


func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()

	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Player, %v", body, err)
	}
	return
}

func AssertLeague(t testing.TB, want, got []Player) {
	t.Helper()
	
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Want %v but got %v", want, got)
	}
}

func AssertContentType(t testing.TB, want string, response *httptest.ResponseRecorder) {
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

func AssertResponseBody(t testing.TB, got, want string, ) {
	t.Helper()

	if got != want {
		t.Errorf("want %q but got %q", want, got)
	}
}

func AssertStatus(t testing.TB, want, got int, ) {
	t.Helper()
	
	if got != want {
		t.Errorf("want status code %d but got %d", want, got)
	}
}

func AssertScoreEquals(t testing.TB, want, got int) {
	t.Helper()

	if got != want {
		t.Errorf("Want %d but got %d", want, got)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)	
	}
}
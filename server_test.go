package poker

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

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

		AssertStatus(t, http.StatusOK, response.Code)
		AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("floyd")
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)

		AssertStatus(t, http.StatusOK, response.Code)
		AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("udin")
		response := httptest.NewRecorder()
		
		server.ServeHTTP(response, request)
		want := http.StatusNotFound
		got := response.Code

		AssertStatus(t, want, got)
	})
}

func TestStoreWins(t *testing.T) {
	AssertTotalCalls := func(t testing.TB, want, got int) {
		t.Helper()
		if want != got {
			t.Errorf("Record win not working correctly. Want %d calls but got %d", want, got)
		}
	}

	AssertCorrectWinner := func(t testing.TB, want, got string) {
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

		AssertStatus(t, http.StatusAccepted, response.Code)
		AssertTotalCalls(t, 1, len(store.winCalls))
		AssertCorrectWinner(t, "pepper", store.winCalls[0])
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

		AssertStatus(t, http.StatusAccepted, response.Code)
		AssertTotalCalls(t, totalReq, len(store.winCalls))
		AssertCorrectWinner(t, player, store.winCalls[0])
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

		AssertContentType(t, jsonContentType, res)
		AssertStatus(t, http.StatusOK, res.Code)		
		AssertLeague(t, wantedLeague, got)
	})
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingsWinsAndRetrievingThem(t *testing.T) {
	database, removeFile := createTempFile(t, "")
	defer removeFile()
	
	store := &FileSystemPlayerStore{database: database}
	server := NewPlayerServer(store)
	player := "pepper"

	// call 3 times should have 3 score
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assertResponseBody(t, response.Body.String(), "3")
		assertStatus(t, http.StatusOK, response.Code)
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetLeagueRequest())
		assertStatus(t, http.StatusOK, response.Code)

		want := []Player{
			{Name: player, Wins: 3},
		}

		got := getLeagueFromResponse(t, response.Body)

		assertContentType(t, jsonContentType, response)
		assertLeague(t, want, got )
	})
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingsWinsAndRetrievingThem(t *testing.T) {
	server := NewPlayerServer(NewInMemoryPlayerStore())
	player := "pepper"

	// call 3 times should have 3 score
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))

	assertResponseBody(t, response.Body.String(), "3")
	assertStatus(t, http.StatusOK, response.Code)
}

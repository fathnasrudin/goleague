package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	t.Run("returns pepper score", func(t *testing.T) {
		request, _:= http.NewRequest(http.MethodGet, "/players/pepper", nil)
		response := httptest.NewRecorder()
		
		PlayerServer(response, request)

		got := response.Body.String()
		want := "20"

		if got != want {
			t.Errorf("want %q but got %q", want, got)
		}
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request, _:= http.NewRequest(http.MethodGet, "/players/floyd", nil)
		response := httptest.NewRecorder()
		
		PlayerServer(response, request)

		got := response.Body.String()
		want := "10"

		if got != want {
			t.Errorf("want %q but got %q", want, got)
		}
	})
}
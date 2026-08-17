package main

import (
	"log"
	"net/http"
	"os"
)
const dbFileName = "game.db.json"
func main() {
	database, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store := &FileSystemPlayerStore{database: database}
	server := NewPlayerServer(store)
	
	log.Fatal(http.ListenAndServe(":5000", server))
}
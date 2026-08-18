package main

import (
	"log"
	"net/http"
	"os"

	poker "github.com/fath-nasrudin/goleague"
)
const dbFileName = "game.db.json"
func main() {
	database, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(database)
	
	if err != nil {
		log.Fatalf("problem creating new file system player store, %v", err)
	}

	server := poker.NewPlayerServer(store)
	
	log.Fatal(http.ListenAndServe(":5000", server))
}
package main

import (
	"log"
	"net/http"
	"os"

	poker "github.com/fath-nasrudin/goleague"
)
const dbFileName = "game.db.json"
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("fail to get PORT from env var")
	}

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)
	log.Fatal(http.ListenAndServe(":"+port, server))
}
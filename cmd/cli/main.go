package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/fath-nasrudin/goleague"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	file, err := os.OpenFile(dbFileName, os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening file %s, %v", file.Name(), err)
	}

	store, err := poker.NewFileSystemPlayerStore(file)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v", err)
	}

	game := poker.NewCLI(store, os.Stdin)
	game.PlayPoker()
}
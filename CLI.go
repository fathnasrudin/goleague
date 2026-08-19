package poker

import (
	"bufio"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in io.Reader
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{store, in}
}

func (cli *CLI) PlayPoker() {
	scanner := bufio.NewScanner(cli.in)
	scanner.Scan()
	input := scanner.Text()
	name := extractWinner(input)
	cli.store.RecordWin(name)
}

func extractWinner(userInput string) string {
	return strings.TrimSuffix(userInput, " wins")
}
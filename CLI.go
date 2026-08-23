package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{store, bufio.NewScanner(in)}
}

func (cli *CLI) PlayPoker() {
	input := cli.readLine()
	name, err := extractWinner(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	cli.store.RecordWin(name)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

func extractWinner(userInput string) (string, error) {
	winSuffix := " wins"

	if  !strings.HasSuffix(userInput, winSuffix) {
		return "", errors.New("failed to record win. Invalid format. should use format {name} wins")
	}

	return strings.TrimSuffix(userInput, winSuffix), nil
}
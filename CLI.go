package poker

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CLI struct {
	store PlayerStore
	in io.Reader
}

func (cli *CLI) PlayPoker() error {
	// extract caller
	scanner := bufio.NewScanner(cli.in)

	for scanner.Scan() {
		input := scanner.Text()
		name := strings.TrimSuffix(input, " wins")
		cli.store.RecordWin(name)
	}

	if scanner.Err() != nil {
		return fmt.Errorf("failed when read input cli, %v", scanner.Err())
	}
	return nil
}
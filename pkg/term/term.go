package term

import (
	"fmt"
	"io"

	"github.com/moby/term"
)

func TerminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, fmt.Errorf("given writer is no terminal")
	}
	winSize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winSize.Width), int(winSize.Height), nil
}

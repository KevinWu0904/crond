package term

import (
	"errors"
	"io"

	"github.com/moby/term"
)

// ErrNoTerminal throws when input io.Writer has no terminal.
var ErrNoTerminal = errors.New("no terminal")

// TerminalSize tries to fetch this io.Writer's terminal window size.
func TerminalSize(w io.Writer) (int, int, error) {
	outFd, isTerminal := term.GetFdInfo(w)
	if !isTerminal {
		return 0, 0, ErrNoTerminal
	}
	winSize, err := term.GetWinsize(outFd)
	if err != nil {
		return 0, 0, err
	}
	return int(winSize.Width), int(winSize.Height), nil
}

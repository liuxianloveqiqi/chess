package handler

import (
	"fmt"
	"testing"
)

func NewInitialBoard(isWhite bool) State {
	board := Board{
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
	}

	if !isWhite {
		return State{
			board: board.Flip(),
		}
	}

	return State{
		board: board,
	}
}
func TestChess(t *testing.T) {
	b := NewInitialBoard(true)
	fmt.Printf("%v", b.Moves())
}

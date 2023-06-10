package handler

func initPostion() Position {
	board, _ := FEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKQBNR")
	return Position{
		board: board,
	}
}

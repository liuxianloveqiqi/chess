package handler

import (
	"fmt"
	"testing"
)

func TestNewInitialBoard(t *testing.T) {
	client := &GameClient{}

	// 测试白方初始化棋盘
	expectedBoard := Board{
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
	}

	actualState := client.NewInitialBoard(true)

	if actualState.board != expectedBoard {
		t.Errorf("白方初始化棋盘不匹配")
	}

	// 测试黑方初始化棋盘
	expectedFlippedBoard := Board{
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	}

	actualFlippedState := client.NewInitialBoard(false)

	if actualFlippedState.board != expectedFlippedBoard {
		t.Errorf("黑方初始化棋盘不匹配")
	}
}

func TestPieceValue(t *testing.T) {
	pieces := map[Piece]int{
		'P': 1,
		'N': 3,
		'B': 3,
		'R': 5,
		'Q': 10,
		'K': 10000,
	}

	for piece, expectedValue := range pieces {
		actualValue := piece.value()
		if actualValue != expectedValue {
			t.Errorf("棋子价值不匹配。棋子: '%c'，期望值: %d，实际值: %d", piece, expectedValue, actualValue)
		}
	}
}

func TestPieceMe(t *testing.T) {
	pieces := map[Piece]bool{
		'P': true,
		'N': true,
		'B': true,
		'R': true,
		'Q': true,
		'K': true,
		'p': false,
		'n': false,
		'b': false,
		'r': false,
		'q': false,
		'k': false,
	}

	for piece, expectedMe := range pieces {
		actualMe := piece.me()
		if actualMe != expectedMe {
			t.Errorf("棋子所有权不匹配。棋子: '%c'，期望值: %t，实际值: %t", piece, expectedMe, actualMe)
		}
	}
}

func TestPieceFlip(t *testing.T) {
	pieces := map[Piece]Piece{
		'P': 'p',
		'N': 'n',
		'B': 'b',
		'R': 'r',
		'Q': 'q',
		'K': 'k',
		'p': 'P',
		'n': 'N',
		'b': 'B',
		'r': 'R',
		'q': 'Q',
		'k': 'K',
		' ': ' ',
		'.': '.',
	}

	for piece, expectedFlippedPiece := range pieces {
		actualFlippedPiece := piece.Flip()
		if actualFlippedPiece != expectedFlippedPiece {
			t.Errorf("棋子翻转不匹配。棋子: '%c'，期望值: '%c'，实际值: '%c'", piece, expectedFlippedPiece, actualFlippedPiece)
		}
	}
}

func TestBoardString(t *testing.T) {
	board := Board{
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		' ', '.', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	}

	expectedString :=
		" -----------------\n" +
			"8|R N B Q K B N R|\n" +
			"7|P P P P P P P P|\n" +
			"6|         .     |\n" +
			"5|               |\n" +
			"4|               |\n" +
			"3|               |\n" +
			"2|p p p p p p p p|\n" +
			"1|r n b q k b n r|\n" +
			" -----------------\n" +
			"   a b c d e f g h\n"

	actualString := board.String()
	if actualString != expectedString {
		t.Errorf("棋盘字符串不匹配。期望值：\n%s\n实际值：\n%s", expectedString, actualString)

	}
}

func TestBoardFlip(t *testing.T) {
	board := Board{
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		' ', '.', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	}

	expectedFlippedBoard := Board{
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
	}

	actualFlippedBoard := board.Flip()

	for i := 0; i < len(expectedFlippedBoard); i++ {
		if actualFlippedBoard[i] != expectedFlippedBoard[i] {
			t.Errorf("在索引%d处翻转板不匹配。期望:'%c'，得到:'%c'", i, expectedFlippedBoard[i], actualFlippedBoard[i])
		}
	}
}

func TestAbs(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{5, 5},
		{-5, 5},
		{10, 10},
		{-10, 10},
		{100, 100},
		{-100, 100},
	}

	for _, testCase := range testCases {
		actual := abs(testCase.input)
		if actual != testCase.expected {
			t.Errorf("输入%d的绝对值不匹配期望:%d，得到:%d", testCase.input, testCase.expected, actual)
		}
	}
}

func TestSquareFlip(t *testing.T) {
	squares := []Square{A1, H1, A8, H8}

	for _, square := range squares {
		expectedFlippedSquare := 63 - square
		actualFlippedSquare := square.Flip()
		if actualFlippedSquare != expectedFlippedSquare {
			t.Errorf("方格翻转不匹配。方格：%s，期望值：%s，实际值：%s", square, expectedFlippedSquare, actualFlippedSquare)
		}
	}
}

func TestSquareString(t *testing.T) {
	squares := []Square{A1, H1, A8, H8}
	expectedStrings := []string{"a1", "h1", "a8", "h8"}

	for i, square := range squares {
		actualString := square.String()
		expectedString := expectedStrings[i]
		if actualString != expectedString {
			t.Errorf("方格字符串不匹配。方格：%s，期望值：%s，实际值：%s", square, expectedString, actualString)
		}
	}
}

func TestMoveString(t *testing.T) {
	moves := []Move{
		{from: A1, to: H8},
		{from: H1, to: A8},
		{from: A8, to: H1},
		{from: H8, to: A1},
	}
	expectedStrings := []string{"a1h8", "h1a8", "a8h1", "h8a1"}

	for i, move := range moves {
		actualString := move.String()
		expectedString := expectedStrings[i]
		if actualString != expectedString {
			t.Errorf("移动字符串不匹配。移动：%s，期望值：%s，实际值：%s", move, expectedString, actualString)
		}
	}
}
func TestStateFlip(t *testing.T) {
	board := Board{
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		' ', '.', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	}

	state := State{
		board: board,
		score: 10,
		wc:    [2]bool{true, false},
		bc:    [2]bool{false, true},
		ep:    A1 + 4*D + 4*W,
		kp:    H1,
	}

	expectedFlippedState := State{
		board: board.Flip(),
		score: -10,
		wc:    [2]bool{false, true},
		bc:    [2]bool{true, false},
		ep:    A1 + 4*D + 5*W,
		kp:    A8,
	}

	actualFlippedState := state.Flip()

	if actualFlippedState.board != expectedFlippedState.board {
		t.Errorf("翻转棋局状态的棋盘不匹配")
	}

	if actualFlippedState.score != expectedFlippedState.score {
		t.Errorf("翻转棋局状态的得分不匹配")
	}

	if actualFlippedState.wc != expectedFlippedState.wc {
		t.Errorf("翻转棋局状态的白方易位条件不匹配")
	}

	if actualFlippedState.bc != expectedFlippedState.bc {
		t.Errorf("翻转棋局状态的黑方易位条件不匹配")
	}

	if actualFlippedState.ep != expectedFlippedState.ep {
		t.Errorf("翻转棋局状态的过路兵位置不匹配")
	}

	if actualFlippedState.kp != expectedFlippedState.kp {
		t.Errorf("翻转棋局状态的易位王位置不匹配")
	}
}

func TestStateMove(t *testing.T) {
	board := Board{
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		' ', '.', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	}

	state := State{
		board: board,
		score: 10,
		wc:    [2]bool{true, false},
		bc:    [2]bool{false, true},
		ep:    A1 + 4*D + 3*W,
		kp:    H1,
	}

	move := Move{from: A1 + 4*D + W, to: A1 + 4*D + 2*W}

	expectedNewState := State{
		board: Board{
			'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
			'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
			' ', '.', ' ', ' ', ' ', ' ', ' ', ' ',
			' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			'p', 'p', 'p', 'p', ' ', 'p', 'p', 'p',
			'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
		},
		score: 11,
		wc:    [2]bool{true, false},
		bc:    [2]bool{false, true},
		ep:    A1 + 4*D + 3*W,
		kp:    H8,
	}

	actualNewState := state.Move(move)

	if actualNewState.board != expectedNewState.board {
		t.Errorf("移动后的棋局状态的棋盘不匹配")
	}

	if actualNewState.score != expectedNewState.score {
		t.Errorf("移动后的棋局状态的得分不匹配")
	}

	if actualNewState.wc != expectedNewState.wc {
		t.Errorf("移动后的棋局状态的白方易位条件不匹配")
	}

	if actualNewState.bc != expectedNewState.bc {
		t.Errorf("移动后的棋局状态的黑方易位条件不匹配")
	}

	if actualNewState.ep != expectedNewState.ep {
		t.Errorf("移动后的棋局状态的过路兵位置不匹配")
	}

	if actualNewState.kp != expectedNewState.kp {
		t.Errorf("移动后的棋局状态的易位王位置不匹配")
	}
}

func TestProcessInput(t *testing.T) {
	client := &GameClient{
		isBout:  true,
		isWhite: true,
		hub: &GameHub{
			clients: make(map[*GameClient]bool),
		},
	}

	// 设置初始状态
	board := Board{
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'.', '.', '.', '.', '.', '.', '.', '.',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	}
	state := State{
		board: board,
		score: 0,
		wc:    [2]bool{true, true},
		bc:    [2]bool{true, true},
		ep:    0,
		kp:    0,
	}
	fmt.Println(state)
	client.hub.clients[client] = true

	// 测试合法的移动
	validInput := "e2e4"
	expectedValidOutput := []byte(" -----------------\n8|r n b q k b n r|\n7|p p p p p p p p|\n6|               |\n5|               |\n4|           P   |\n3|               |\n2|P P P P P P P P|\n1|R N B Q K B N R|\n -----------------\n   a b c d e f g h\n")
	actualValidOutput := client.processInput(validInput)

	if string(actualValidOutput) != string(expectedValidOutput) {
		t.Errorf("合法移动处理不正确。期望输出：%s，实际输出：%s", string(expectedValidOutput), string(actualValidOutput))
	}

	// 测试非法的移动
	invalidInput := "e2e5"
	expectedInvalidOutput := []byte("输入非法，请重新输入")
	actualInvalidOutput := client.processInput(invalidInput)

	if string(actualInvalidOutput) != string(expectedInvalidOutput) {
		t.Errorf("非法移动处理不正确。期望输出：%s，实际输出：%s", string(expectedInvalidOutput), string(actualInvalidOutput))
	}

	// 测试非当前回合的移动
	client.isBout = false
	notBoutInput := "e7e5"
	expectedNotBoutOutput := []byte("现在不是白方的回合，请等待对方走完")
	actualNotBoutOutput := client.processInput(notBoutInput)

	if string(actualNotBoutOutput) != string(expectedNotBoutOutput) {
		t.Errorf("非当前回合移动处理不正确。期望输出：%s，实际输出：%s", string(expectedNotBoutOutput), string(actualNotBoutOutput))
	}
}

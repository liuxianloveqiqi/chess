package handler

import (
	"fmt"
	"math"
)

// 棋子
type Piece byte

// 计算棋子价值
func (p Piece) value() int {
	return map[Piece]int{'P': 1, 'N': 3, 'B': 3, 'R': 5, 'Q': 10, 'K': 10000}[p]
}

// 判断棋子是否属于我方
func (p Piece) me() bool {
	return p.value() > 0
}

// 翻转棋子
func (p Piece) Flip() Piece {
	return map[Piece]Piece{'P': 'p', 'N': 'n', 'B': 'b', 'R': 'r', 'Q': 'q', 'K': 'k', 'p': 'P', 'n': 'N', 'b': 'B', 'r': 'R', 'q': 'Q', 'k': 'K', ' ': ' ', '.': '.'}[p]
}

type Board [64]Piece

// 打印棋盘
func (a Board) String() string {
	s := ""
	s += " -----------------\n"
	for row := 0; row < 8; row++ {
		s += string('1'+row) + "|"
		for col := 0; col < 8; col++ {
			s += " " + string(a[row*8+col])
		}
		s += " |\n"
	}
	s += " -----------------\n"
	s += "   a b c d e f g h\n"
	return s
}

// 翻转棋盘
//
//	func (a Board) Flip() string {
//		s := "   8 7 6 5 4 3 2 1\n"
//		s += " -----------------\n"
//		for row := 9; row >= 2; row-- {
//			s += string('A'+row-2) + "|"
//			for col := 1; col < 9; col++ {
//				piece := a[row*10+col]
//				if piece.me() {
//					piece = piece.Flip()
//				}
//				s += " " + string(piece)
//			}
//			s += " |\n"
//		}
//		s += " -----------------\n"
//		s += "   A B C D E F G H\n"
//		return s
//	}
// 翻转棋盘

func (a Board) Flip() (b Board) {
	for i := len(a) - 1; i >= 0; i-- {
		b[i] = a[len(a)-i-1].Flip()
	}
	return b
}

type Square int

// 棋盘上的四个角的位置
const A1, H1, A8, H8 Square = 0, 7, 56, 63

func (s Square) Flip() Square {
	return 63 - s
}

// 位置转换为类似 "a1"、"b2" 的字符串形式
func (s Square) String() string {
	return string([]byte{" abcdefgh "[s%8], "  87654321  "[s/8]})
}

// 移动
type Move struct {
	from Square
	to   Square
}

// 返回移动字符串
func (m Move) String() string {
	return m.from.String() + m.to.String()
}

// 棋局状态
type State struct {
	board Board   // 棋盘
	score int     // 当前棋盘位置的得分
	wc    [2]bool // wc[0] 表示白方短易位是否可行，wc[1] 表示白方长易位是否可行
	bc    [2]bool // bc表示黑方短易位是否可行，bc[1] 表示黑方长易位是否可行
	ep    Square  // 可以进行吃过路兵的位置。
	kp    Square  // 易位过程中可以吃掉王的位置
}

// 翻转
func (s State) Flip() State {
	newState := State{
		score: -s.score,
		wc:    [2]bool{s.bc[0], s.bc[1]},
		bc:    [2]bool{s.wc[0], s.wc[1]},
		ep:    s.ep.Flip(),
		kp:    s.kp.Flip(),
	}
	newState.board = s.board.Flip()
	return newState
}

// 上下左右
const W, D, S, A = 8, 1, -8, -1

func (s State) Moves() (moves []Move) {
	// 所以棋子可以走的步
	fmt.Println("傻逼吧日")
	directions := map[Piece][]Square{
		'P': {W, W + W, W + A, W + D},
		'N': {W + W + D, D + W + D, D + S + D, S + S + D, S + S + A, A + S + A, A + W + A, W + W + A},
		'B': {W + D, S + D, S + A, W + A},
		'R': {W, D, S, A},
		'Q': {W, D, S, A, W + D, S + D, S + A, W + A},
		'K': {W, D, S, A, W + D, S + D, S + A, W + A},
	}

	// 遍历棋盘上的每个方格，并检查该方格上的棋子是否属于当我方
	for k, p := range s.board {
		if !p.me() {
			continue
		}
		i := Square(k)
		// 遍历该棋子可以移动的所有方向
		for _, d := range directions[p] {

			for j := i + d; ; j += d {
				// 遇到空白或者己方棋子，跳出
				q := s.board[j]

				if q == ' ' || (q.me() && q != '.') {
					break
				}

				// 检验兵的特殊情况
				if p == 'P' {

					if (d == W || d == W+W) && q != '.' {
						break
					}

					// 双步起步，兵不在起始位置（即 i < A1+W），或者在兵的前进方向上存在障碍物
					if d == W+W && (i < A1+W || s.board[i+W] != '.') {
						break
					}

					// 兵左右上方没棋子，不是吃过路兵的位置，不是吃易位王的位置
					if (d == W+A || d == W+D) && q == '.' && (j != s.ep && j != s.kp && j != s.kp-1 && j != s.kp+1) {
						break
					}
				}
				// 添加合法步
				fmt.Println("开始添加和法步")
				moves = append(moves, Move{from: i, to: j})
				// 不能越过对方的棋子
				if p == 'P' || p == 'N' || p == 'K' || (q != '.' && !q.me()) {
					break
				}
				// 判断白方王车易位
				// 短易位
				if i == A1 && s.board[j+D] == 'K' && s.wc[0] {
					moves = append(moves, Move{from: j + D, to: j + A})
				}
				// 长易位
				if i == H1 && s.board[j+A] == 'K' && s.wc[1] {
					moves = append(moves, Move{from: j + A, to: j + D})
				}
				// 判断黑方王车易位
				if i == A1 && s.board[j+D] == 'K' && s.bc[0] {
					moves = append(moves, Move{from: j + D, to: j + A})
				}
				if i == H1 && s.board[j+A] == 'K' && s.bc[1] {
					moves = append(moves, Move{from: j + A, to: j + D})
				}

			}
		}
	}
	fmt.Println("这就是moves", moves)
	return moves
}

func (s State) Move(m Move) (newState State) {
	// 初始化
	i, j, p := m.from, m.to, s.board[m.from]
	newState = s
	newState.ep = 0
	newState.kp = 0
	// 棋子易位
	newState.board[m.to] = s.board[m.from]
	newState.board[m.from] = '.'
	// 王车易位，车在角上就判断可以进行易位
	newState.wc[0] = newState.wc[0] || (i == A1)
	newState.wc[1] = newState.wc[1] || (i == H1)
	newState.bc[1] = newState.bc[1] || (j == A8)
	newState.bc[0] = newState.bc[0] || (j == H8)

	if p == 'K' {
		// 移动王就否定易位条件
		newState.wc[0], newState.wc[1] = false, false
		// 差两步判断为王车易位
		if abs(int(j-i)) == 2 {
			if j > i {
				// 王向右移动，进行短易位
				if s.board[H1+A] == '.' && s.board[H1+A+A] == '.' {
					newState.board[H1] = '.'
				}
			} else {
				// 王向左移动，进行长易位
				if s.board[A1+D] == '.' && s.board[A1+D+D] == '.' && s.board[A1+D+D+D] == '.' {
					newState.board[H1] = '.'
				}
			}
			newState.board[(i+j)/2] = 'R'
		}
	}
	if p == 'P' {
		// 判断是否双步起步
		if j-i == 2*S {
			newState.ep = i + S
		}

		// 判断是否能吃过路兵
		if j == s.ep {
			newState.board[j+S] = '.'
		}
		// 兵升后
		if A8 <= j && j <= H8 {
			newState.board[j] = 'Q'
		}

	}
	for _, p := range s.board {
		if !p.me() {
			continue
		} else {
			newState.score += p.value()
		}
	}
	return newState.Flip()
}

func abs(v int) int {
	return int(math.Abs(float64(v)))
}

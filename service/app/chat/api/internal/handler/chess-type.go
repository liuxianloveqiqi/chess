package handler

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

type Board [120]Piece

// 打印棋盘
func (a Board) String() string {
	s := ""
	s += " -----------------\n"
	for row := 2; row < 10; row++ {
		s += string('a'+row-2) + "|"
		for col := 1; col < 9; col++ {
			s += " " + string(a[row*10+col])
		}
		s += " |\n"
	}
	s += " -----------------\n"
	s += "   1 2 3 4 5 6 7 8\n"
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
const A1, H1, A8, H8 Square = 91, 98, 21, 28

func (s Square) Flip() Square {
	return 119 - s
}

// 位置转换为类似 "a1"、"b2" 的字符串形式
func (s Square) String() string {
	return string([]byte{" abcdefgh "[s%10], "  87654321  "[s/10]})
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

// 位置
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
	np := State{
		score: -s.score,
		wc:    [2]bool{s.bc[0], s.bc[1]},
		bc:    [2]bool{s.wc[0], s.wc[1]},
		ep:    s.ep.Flip(),
		kp:    s.kp.Flip(),
	}
	np.board = s.board.Flip()
	return np
}

// 上下左右
const W, D, S, A = -10, 1, 10, -1

func (s State) Moves() (moves []Move) {
	// 所以棋子可以走的步
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

					if (d == W+A || d == W+D) && q == '.' && (j != s.ep && j != s.kp && j != s.kp-1 && j != s.kp+1) {
						break
					}
				}
				moves = append(moves, Move{from: i, to: j})
				// Crawling pieces should stop after a single move
				if p == 'P' || p == 'N' || p == 'K' || (q != ' ' && q != '.' && !q.me()) {
					break
				}
				// Castling rules
				if i == A1 && s.board[j+D] == 'K' && s.wc[0] {
					moves = append(moves, Move{from: j + D, to: j + A})
				}
				if i == H1 && s.board[j+A] == 'K' && s.wc[1] {
					moves = append(moves, Move{from: j + A, to: j + D})
				}
			}
		}
	}
	return moves
}

package handler

import (
	"fmt"
	"strings"
)

func (c *GameClient) processInput(input string) []byte {
	var res []byte

	// 去除输入字符串首尾的空白字符。
	input = strings.TrimSpace(input)

	// 标志是否发现合法的移动
	validMove := false

	// 遍历当前位置的合法移动列表。
	for _, m := range state.Moves() {
		// 检查用户输入的移动是否与当前遍历的移动相匹配。
		if input == m.String() {
			// 如果输入的移动有效，则更新当前位置为移动后的位置
			state = state.Move(m)
			res = []byte(state.board.String())
			if state.score < 1000 {
				res = []byte(fmt.Sprintf("%v 输了，游戏结束", c.hub.stringWhiteOrBlack(c.isWhite)))
			}
			validMove = true
			break
		}
	}

	// 如果没有发现合法的移动，返回错误提示
	if !validMove {
		res = []byte("输入非法，请重新输入")
	}

	return res
}

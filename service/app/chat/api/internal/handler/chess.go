package handler

import (
	"fmt"
	"strings"
)

func (c *GameClient) processInput(input string) []byte {
	var res []byte
	if !c.isBout {
		return []byte(fmt.Sprintf("现在不是%v的回合，请等待对方走完", c.hub.stringWhiteOrBlack(c.isWhite)))
	}
	// 去除输入字符串首尾的空白字符。
	input = strings.TrimSpace(input)

	// 标志是否发现合法的移动
	validMove := false
	fmt.Println("开始操作")
	// 遍历当前位置的合法移动列表。
	c.hub.mutex.Lock()
	for _, m := range state.Moves() {
		fmt.Println("开始遍历")
		// 检查用户输入的移动是否与当前遍历的移动相匹配。
		if input == m.String() {
			// 如果输入的移动有效，则更新当前位置为移动后的位置
			state = state.Move(m)
			res = []byte(state.board.String())
			for client, _ := range c.hub.clients {
				if client != c {
					client.isBout, c.isBout = true, false
				}
			}
			if state.score < 1000 {
				res = []byte(fmt.Sprintf("%v 输了，游戏结束", c.hub.stringWhiteOrBlack(c.isWhite)))
			}
			validMove = true
			break
		}
	}
	c.hub.mutex.Unlock()
	// 如果没有发现合法的移动，返回错误提示
	if !validMove {
		res = []byte("输入非法，请重新输入")
	}
	fmt.Println(string(res))
	return res
}

package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func (c *GameClient) processInput(input []byte) []byte {
	// 处理用户输入的消息
	if string(input) == "start" {
		c.isReady = true
	}

	// 返回引擎的响应消息
	// ...
	pos := start()
	searcher := &Searcher{tp: map[Position]entry{}}
	r := bufio.NewReader(bytes.NewReader(input))
	var output bytes.Buffer

	for {
		output.WriteString(fmt.Sprintln(pos.board))
		valid := false
		for !valid {
			output.WriteString("请输入移动：")
			input, _ := r.ReadString('\n')
			input = strings.TrimSpace(input)
			valid = false
			for _, m := range pos.Moves() {
				if input == m.String() {
					pos = pos.Move(m)
					valid = true
					break
				}
			}
		}
		output.WriteString(fmt.Sprintln(pos.Flip().board))
		m := searcher.Search(pos, 10000)
		score := pos.value(m)
		if score < 1000{
			output.WriteString("你输了\n")
			return output.Bytes()
		}
		if score >= MateValue {
			output.WriteString("你输了\n")
			return output.Bytes()
		}
		pos = pos.Move(m)
	}
}
}

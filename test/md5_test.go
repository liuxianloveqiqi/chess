package test

import (
	"chess/service/common/utils"
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(utils.Md5Password("123456", "liuxian"))
}

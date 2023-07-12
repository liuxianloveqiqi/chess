package logic

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"user/rpc/types/user"
)

func TestLoginLogic_Login(t *testing.T) {
	// 创建测试用例的输入数据
	in := &user.LoginReq{
		PhoneOrEmail: "test@example.com",
		PassWord:     "password",
	}

	// 创建 LoginLogic 实例并注入所需的依赖项
	logic := NewLoginLogic(context.Background(), nil)
	// 执行被测试的函数
	resp, err := logic.Login(in)

	// 验证函数的返回结果和错误
	assert.NoError(t, err)
	assert.Equal(t, &user.CommonResp{UserId: 1}, resp)
}

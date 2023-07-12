package handler

import (
	"chat/common/initdb"
	"chat/common/utils"
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"testing"
	"time"
)

func TestMd5(t *testing.T) {
	fmt.Println(utils.Md5Password("123456", "liuxian"))
}

var accessSecret = []byte("liuxian123")
var refreshSecret = []byte("123456789")

func TestGetToken(t *testing.T) {
	id := int64(123)
	state := "example"
	accessToken, refreshToken := utils.GetToken(id, state)

	if accessToken == "" || refreshToken == "" {
		t.Errorf("GetToken 返回了空的令牌")
	}
}

func TestParseToken(t *testing.T) {
	id := int64(123)
	state := "example"
	accessToken, refreshToken := utils.GetToken(id, state)

	claims, isRefreshToken, err := utils.ParseToken(accessToken, refreshToken)

	if err != nil {
		t.Errorf("ParseToken 返回了错误：%v", err)
	}

	if claims.ID != id || claims.State != state {
		t.Errorf("ParseToken 没有正确解析声明")
	}

	if isRefreshToken {
		t.Errorf("ParseToken 错误地将访问令牌标识为刷新令牌")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	claims, isRefreshToken, err := utils.ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MCwic3RhdGUiOiIyNmVhMGJlYS01Y2Q4LTQwYzUtYjZlMS0zMDliZDBlN2ZkYWYiLCJleHAiOjE2ODg0NjUxOTQsImlhdCI6MTY4ODQ2NTAxNCwiaXNzIjoiQVIifQ.BnWmXU4tgoZ6GVJqZ2Pg4iQFbDejUCWuaFHtJ72SFe8", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6MCwic3RhdGUiOiIyNmVhMGJlYS01Y2Q4LTQwYzUtYjZlMS0zMDliZDBlN2ZkYWYiLCJleHAiOjE2OTEwNTcwMTQsImlhdCI6MTY4ODQ2NTAxNCwiaXNzIjoiUlQifQ.Ok1ivXtSpxl3ustqcf78whyzn9OFEId1JVqBWAnjFJc")

	if claims != nil || isRefreshToken || err == nil {
		t.Errorf("ParseToken 没有返回无效令牌的错误")
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	// Generate an expired access token
	expiredTime := time.Now().Add(-time.Minute)
	expiredClaims := utils.MyClaims{
		ID:    123,
		State: "example",
		StandardClaims: jwt.StandardClaims{
			Issuer:    "AR",
			IssuedAt:  expiredTime.Unix(),
			ExpiresAt: expiredTime.Unix(),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString(accessSecret)

	// Parse the expired token
	claims, isRefreshToken, err := utils.ParseToken(expiredTokenString, "validRefreshToken")

	if claims != nil || isRefreshToken || err == nil {
		t.Errorf("ParseToken 没有返回已过期令牌的错误")
	}
}
func TestSMS(t *testing.T) {
	phone := "123456789"
	// 密钥不能暴露
	secretId := "xxxxx"
	secretKey := "zzzzz"
	ctx := context.TODO()
	// 创建一个模拟的Redis客户端
	rdb := initdb.InitRedis("43.139.195.17:6333")
	code := utils.SMS(phone, secretId, secretKey, ctx, rdb)
	fmt.Println(code)
}

func TestDefaultGetValidParams_Valid(t *testing.T) {
	ctx := context.Background()
	params := struct {
		Phone string `validate:"required,phone"`
	}{
		Phone: "13812345678",
	}

	err := utils.DefaultGetValidParams(ctx, params)

	if err != nil {
		t.Errorf("期望无错误，实际获得：%v", err)
	}
}

func TestDefaultGetValidParams_Invalid(t *testing.T) {
	ctx := context.Background()
	params := struct {
		Phone string `validate:"required,phone"`
	}{
		Phone: "123456789",
	}

	err := utils.DefaultGetValidParams(ctx, params)

	expectedError := errors.New("phone格式不正确，必须为手机号码")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("期望错误：%v，实际获得：%v", expectedError, err)
	}
}

func TestDefaultGetValidParams_ValidatorNotFound(t *testing.T) {
	ctx := context.Background()
	params := struct {
		Phone string `validate:"required,phone"`
	}{
		Phone: "13812345678",
	}

	err := utils.DefaultGetValidParams(ctx, params)

	expectedError := errors.New("在上下文中未找到验证器")
	if err == nil || err.Error() != expectedError.Error() {
		t.Errorf("期望错误：%v，实际获得：%v", expectedError, err)
	}
}
func TestRandNickname(t *testing.T) {
	// 固定随机种子以确保可重复性
	rand.Seed(0)

	nickname := utils.RandNickname()

	expectedFormat := "user%08d"
	expectedLength := len(fmt.Sprintf(expectedFormat, 0))

	if len(nickname) != expectedLength {
		t.Errorf("生成的昵称长度不正确，期望长度：%d，实际长度：%d", expectedLength, len(nickname))
	}

	if nickname[:4] != "user" {
		t.Errorf("生成的昵称前缀不正确，期望前缀：%s，实际前缀：%s", "user", nickname[:4])
	}

	numStr := nickname[4:]
	_, err := fmt.Sscanf(numStr, "%08d", new(int))
	if err != nil {
		t.Errorf("生成的昵称数字部分不正确，无法解析为有效数字：%s", numStr)
	}
}
func TestInitGorm(t *testing.T) {
	// 使用测试数据库的连接字符串
	db := initdb.InitGorm("root:root@(43.139.195.17:6446)/gopan?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai")

	if db == nil {
		t.Errorf("无法连接到MySQL数据库")
	}

}

func TestInitRedis(t *testing.T) {
	// 使用测试Redis服务器的地址
	rdb := initdb.InitRedis("43.139.195.17:6333")

	if rdb == nil {
		t.Errorf("无法连接到Redis服务器")
	}

}

func TestInitNsqProduct(t *testing.T) {
	// 使用测试NSQ服务器的地址
	producer := initdb.InitNsqProduct("127.0.0.1:4150")

	if producer == nil {
		t.Errorf("无法连接到NSQ服务器")
	}

}

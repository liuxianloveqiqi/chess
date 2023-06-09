// Code generated by goctl. DO NOT EDIT.
package types

type RegisterByPhoneRep struct {
	UserPhone string `json:"userPhone" validate:"required,phone"`
}

type RegisterByPhoneResp struct {
	VeCode string `json:"veCode"`
}

type RegisterReq struct {
	UserPhone string `json:"userPhone" validate:"required,phone"`
	VeCode    string `json:"veCode" validate:"required,len=6"`
}

type TokenResp struct {
	UserId       int64  `json:"userId"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LoginReq struct {
	PhoneOrEmail string `json:"phoneOrEmail" validate:"required"` // 手机号或者邮箱
	PassWord     string `json:"PassWord"`                         // 用户密码，MD5加密
}

type UserInfoReq struct {
	UserId int64 `json:"userId"` // 用户id
}

type UserInfoResp struct {
	UserInfo *UserInfoItem `json:"Data"`
}

type UserInfoItem struct {
	Id         int64  `json:"userId"`     // 用户ID
	Password   string `json:"password"`   // 用户密码，MD5加密
	UserNick   string `json:"userNick"`   // 用户昵称
	UserSex    int64  `json:"userSex"`    // 用户性别：0男，1女，2保密
	UserEmail  string `json:"userEmail"`  // 用户邮箱
	UserPhone  string `json:"userPhone"`  // 手机号
	CreateTime string `json:"createTime"` // 创建时间
	UpdateTime string `json:"updateTime"` // 更新时间
	DeleteTime string `json:"deleteTime"` // 删除时间
}

type CommonResply struct {
	Code    int64  `json:"Code"`
	Message string `json:"Message"`
	Data    string `json:"Data"`
}

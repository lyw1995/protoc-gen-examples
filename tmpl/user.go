// Code generated by protoc-gen-example.
// source: user.proto

package user

import (
	"time"
)

// 用户权限
type Role int32

const (
	// 管理员
	Role_Admin Role = 0
	// 普通用户
	Role_User Role = 1
)

// 登录请求
type LoginReq struct {
	// 用户昵称
	Name string `json:"name,omitempty"`
	// 用户密码
	Password string `json:"password,omitempty"`
}

// 用户ID请求
type IdReq struct {
	// 用户id
	Id int64 `json:"id,omitempty"`
}

// 登录响应
type LoginResp struct {
	// 用户id
	UserId string `json:"userId,omitempty"`
	// 消息
	Msg string `json:"msg,omitempty"`
	// 列表
	List []interface{} `json:"list,omitempty"`
	// 创建时间
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	// 权限
	Role Role `json:"role,omitempty"`
}

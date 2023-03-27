package serializer

import (
	"main/Video_Web/model"
)

type User struct {
	ID       uint   `json:"id" form:"id" example:"1"`                  // 用户ID
	UserName string `json:"user_name" form:"user_name" example:"Fan1"` //用户名
	Status   string `json:"status" form:"status"`                      //用户状态
	CreateAt int64  `json:"create_at" form:"create_at"`                //创建时间
}

type ChangePassword struct {
	UserName    string `json:"id" `
	PrePassword string `json:"prepassword"`
	NewPassword string `json:"newpassword"`
}

type UserMsg struct {
	Uid      uint   `json:"uid"`
	Img      string `json:"img"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Gender   int    `json:"gender"`
	Birthday string `json:"birthday"`
	Sign     string `json:"sign"`
}

// 序列化用户
func BuildUser(user model.User) User {
	return User{
		ID:       user.ID,
		UserName: user.UserName,
		CreateAt: user.CreatedAt.Unix(),
	}
}

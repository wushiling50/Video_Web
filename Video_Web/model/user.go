package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户列表
type User struct {
	gorm.Model
	UserName       string `gorm:"type:varchar(20);not null;unique"`
	PasswordDigest string //存储的是密文
	Img            string `gorm:"size:1000"`
	Email          string `gorm:"type:varchar(30);not null;index"`
	Gender         int    `gorm:"default:0"` //0 为男性 , 1 为女性
	Birthday       string `gorm:"default:'1970年01月01日'"`
	Sign           string `gorm:"type:varchar(50);default:'这个人很懒，什么都没有留下'"`
	Role           int    `gorm:"size:1;default:0"` //0为普通用户 ，1为管理员
	State          int    `gorm:"size:1;default:0"` //0为正常状态，1为拉黑状态， 2为封禁状态
}

////-----------------------------------------------------------------////
//																	  //
//	拉黑为禁止该用户的视频上传，评论与弹幕发送；封禁则禁止用户使用一切功能  //
//																	 //
// //-------------------------------------------------------------////

// 加密
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// 验证密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}

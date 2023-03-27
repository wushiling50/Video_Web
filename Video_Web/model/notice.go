package model

import "gorm.io/gorm"

// 存放邮件模板
type Notice struct {
	gorm.Model
	Text string `gorm:"type:text"`
}

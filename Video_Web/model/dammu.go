package model

import "gorm.io/gorm"

//弹幕
type Danmu struct {
	gorm.Model
	Vid  string `gorm:"not null;index"`
	Type uint   `gorm:"not null"` //类型0滚动;1顶部;2底部
	Text string `gorm:"type:varchar(100);not null"`
	Uid  uint   `gorm:"not null"` //发送人的id
}

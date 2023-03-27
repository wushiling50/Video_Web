package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Vid     string `gorm:"not null;index"`             //视频ID
	Uid     uint   `gorm:"not null"`                   //评论用户的ID
	Level   uint   `gorm:"not null;default:0"`         //楼层
	Content string `gorm:"type:varchar(255);not null"` //内容

	ParentID uint `gorm:"default: 0"` //回复的 目标评论用户的ID
	State    uint `gorm:"not null"`   //1为评论 ，2 为回复
}

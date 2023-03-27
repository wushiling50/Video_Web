package model

import "gorm.io/gorm"

//收藏夹
type Collection struct {
	gorm.Model
	Uid          uint   `gorm:"not null"`          //所属用户
	CollectionID uint   `gorm:"default:0"`         //收藏夹ID
	Name         string `gorm:"type:varchar(20);"` //收藏夹名称
	Open         uint   `gorm:"default:0"`         //是否公开(0为私密，1为公开)
}

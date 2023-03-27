package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	User User `gorm:"ForeignKey:uid;AssociationForeignKey:id"`
	Uid  uint `gorm:"not null;index"` //所属用户

	Vid      string `gorm:"not null"`                          //视频号
	Title    string `gorm:"type:varchar(50);not null;index"`   //视频标题
	Resource string `gorm:"size:1000"`                         //视频源
	Desc     string `gorm:"type:varchar(200);default:'什么都没有'"` //视频简介
	Clicks   int    `gorm:"default:0"`                         //点击量
	Review   int    `gorm:"not null;default:0"`                //审核状态(0为待审，1为过审)
}

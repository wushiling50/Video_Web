package model

import "gorm.io/gorm"

type Collect struct {
	gorm.Model
	Cid uint `gorm:"not null"` //所属收藏夹ID

	Uid uint   `gorm:"not null"` //收藏者的ID
	Vid string `gorm:"not null"` //视频的ID

}

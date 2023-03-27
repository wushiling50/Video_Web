package model

import "gorm.io/gorm"

//点赞

type Liked struct {
	gorm.Model
	Uid     uint   `gorm:"not null"`
	Vid     string `gorm:"not null"`
	IsLiked bool   `gorm:"default:false"` //是否点赞

}

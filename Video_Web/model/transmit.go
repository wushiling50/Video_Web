package model

import "gorm.io/gorm"

type Transmit struct {
	gorm.Model
	Uid     uint   `gorm:"not null"`
	Vid     string `gorm:"not null;index"`
	Path    string `gorm:"size:1000"`
	Comment string `gorm:"type:varchar(100);not null"`
}

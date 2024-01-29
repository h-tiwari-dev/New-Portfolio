package models

import "app/db"

//type List struct {
//ID  uint   `gorm:"primaryKey"`
//	Url string `gorm:"not null; unique"`
//
//	UpdatedAt  int
//	Categories []Category `gorm:"foreignkey:Url;references:Url"`
//}

func Migrate() {
	db.DB.AutoMigrate()
}

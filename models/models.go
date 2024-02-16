package models

import (
	"time"

	"app/db"
)

type Blog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `                  json:"title"`
	Content   string    `                  json:"content"`
	CreatedAt time.Time `                  json:"created_at"`
	UpdatedAt time.Time `                  json:"updated_at"`
}

type Image struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Filename  string    `                  json:"filename"`
	Data      []byte    `                  json:"data"`
	CreatedAt time.Time `                  json:"created_at"`
}

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `                  json:"username"`
	Role     string `                  json:"role"`
	Token    string `                  json:"token"`
	Password string `                  json:"-"`
}

func Migrate() {
	db.DB.AutoMigrate(&Blog{}, &Image{}, &User{})
}

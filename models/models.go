package models

import (
	"time"

	"app/db"
)

type Blog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `                  json:"title"`
	Description string    `                  json:"description"`
	Content     string    `                  json:"content"`
	CreatedAt   time.Time `                  json:"created_at"`
	UpdatedAt   time.Time `                  json:"updated_at"`
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
	Email    string `                  json:"email"`
	Role     string `                  json:"role"`
	Password string `                  json:"password"`
}

func Migrate() {
	db.DB.AutoMigrate(&Blog{}, &Image{}, &User{})
}

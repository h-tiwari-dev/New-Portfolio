package models

import (
	"time"

	"app/db"

	"gorm.io/gorm"
)

type Blog struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `                  json:"title"`
	Description string    `                  json:"description"`
	ContentHTML string    `                  json:"content_html"`
	ContentMD   string    `                  json:"content_md"`
	CreatedAt   time.Time `                  json:"created_at"`
	UpdatedAt   time.Time `                  json:"updated_at"`
	Images      []Image   `gorm:"many2many:blog_images;" json:"images"`
}

type Image struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey" json:"id"`
	Filename  string    `json:"filename"`
	TopImage  uint      `gorm:"default:0;" json:"top_image"`
	Data      []byte    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `                  json:"username"`
	Email    string `                  json:"email"`
	Role     string `                  json:"role"`
	Password string `                  json:"password"`
}

func Migrate() {
	db.DB.AutoMigrate(&Blog{}, &Image{}, &User{})
}


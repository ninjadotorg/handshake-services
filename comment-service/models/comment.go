package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm"
	_ "encoding/gob"
	"time"
)

type Comment struct {
	DateCreated  time.Time
	DateModified time.Time
	ID           int64
	UserId       int64
	ObjectType   string
	ObjectId     int64
	Comment      string
	Image        string
	Status       int
}

func (Comment) TableName() string {
	return "comment"
}
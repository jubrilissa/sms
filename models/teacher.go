package models

import (
	"github.com/jinzhu/gorm"
)

// Teacher field (Model) defined
type Teacher struct {
	gorm.Model
	Name     string
	Address  string
	PhoneNo  string
	Religion string
}

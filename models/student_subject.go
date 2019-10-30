package models

import (
	"github.com/jinzhu/gorm"
)

type StudentSubjectClass struct {
	gorm.Model
	StudentID      uint
	SubjectClassID uint
	// TODO: Not really a todo. just setting this here because it should come in handy and we don't want to delete past records
	IsActive bool
}

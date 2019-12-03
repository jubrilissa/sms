package models

import "github.com/jinzhu/gorm"

type grade struct {
	gorm.Model
	StudentID uint
	SubjectClassID uint
	ExamScore float32
	CaScore float32
	Total float32
}

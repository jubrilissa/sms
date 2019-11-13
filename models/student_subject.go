package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type StudentSubjectClass struct {
	gorm.Model
	StudentID      uint
	SubjectClassID uint
	// TODO: Not really a todo. just setting this here because it should come in handy and we don't want to delete past records
	IsActive bool
}

func (studentSubjectClass *StudentSubjectClass) Create() {

	fmt.Println(studentSubjectClass, "the class object")

	fmt.Println("Just before running create")
	GetDB().Create(&studentSubjectClass)
}

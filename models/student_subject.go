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

func UpdateStudentSubjectClassTeacher(id uint, teacher uint) *SubjectClass {
	subjectClass := &SubjectClass{}
	err := GetDB().Table("student_subject_classes").Where("id = ?", id).First(subjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}

	GetDB().Model(&subjectClass).Update("user_id", teacher)
	return subjectClass

}

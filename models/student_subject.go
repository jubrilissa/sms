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
	IsActive    bool
	FirstCA     float32
	SecondCA    float32
	FirstExam   float32
	TotalFirst  float32
	GradeFirst  string
	RemarkFirst string
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

// GetStudentSubjectsClassBySubjectClassID - Return the students taken a specific subjct in a class
func GetStudentSubjectsClassBySubjectClassID(subjectClassID uint) []*StudentSubjectClass {
	studentSubjectClass := make([]*StudentSubjectClass, 0)
	err := GetDB().Table("student_subject_classes").Where("subject_class_id = ?", subjectClassID).Find(&studentSubjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return studentSubjectClass
}

// UpdateStudentscore - Update a student score whether CA or Exams
func UpdateStudentscore(id uint, fieldName string, score float32) *StudentSubjectClass {
	studentSubjectClass := &StudentSubjectClass{}
	err := GetDB().Table("student_subject_classes").Where("id = ?", id).First(studentSubjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}

	switch fieldName {
	case "firstCA":
		GetDB().Model(&studentSubjectClass).Update("first_ca", score)
	case "secondCA":
		GetDB().Model(&studentSubjectClass).Update("second_ca", score)
	case "firstExam":
		GetDB().Model(&studentSubjectClass).Update("first_exam", score)
	}

	return studentSubjectClass
}

// GetStudentSubjectsClassByStudentID - Return the student subject grade detail for a given student
func GetStudentSubjectsClassByStudentID(studentID uint) []*StudentSubjectClass {
	studentSubjectClass := make([]*StudentSubjectClass, 0)
	err := GetDB().Table("student_subject_classes").Where("student_id = ?", studentID).Find(&studentSubjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return studentSubjectClass
}

func UpdateStudentSubject(id uint, total float32, grade string, remark string) *StudentSubjectClass {
	studentSubjectClass := &StudentSubjectClass{}
	err := GetDB().Table("student_subject_classes").Where("id = ?", id).First(studentSubjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}

	GetDB().Model(&studentSubjectClass).Update(
		StudentSubjectClass{
			TotalFirst:  total,
			GradeFirst:  grade,
			RemarkFirst: remark,
		})
	return studentSubjectClass

}

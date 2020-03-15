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
	IsActive     bool
	FirstCA      float32
	SecondCA     float32
	FirstExam    float32
	TotalFirst   float32
	GradeFirst   string
	RemarkFirst  string
	SFirstCA     float32
	SSecondCA    float32
	SecondExam   float32
	TotalSecond  float32
	GradeSecond  string
	RemarkSecond string
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

type SubjectListRow struct {
	Name        string
	FirstCA     float64
	SecondCA    float64
	FirstExam   float64
	TotalFirst  float64
	GradeFirst  string
	RemarkFirst string
}

// GetStudentSubjectListRowBytID - Return the student subject grade detail for a given student
func GetStudentSubjectListRowByStudentID(studentID uint) []*SubjectListRow {
	// 	SELECT sb.name, ssc.first_ca, ssc.second_ca, ssc.first_exam, ssc.total_first, ssc.grade_first, ssc.remark_first
	// FROM students as s
	// INNER JOIN student_subject_classes as ssc ON s.id = ssc.student_id
	// INNER JOIN subject_classes as sc ON sc.id = ssc.subject_class_id
	// INNER JOIN subjects as sb on sb.id = sc.subject_id
	// WHERE s.id = 32;

	// SELECT * FROM students WHERE id = 32;
	var result []*SubjectListRow

	rows, _ := GetDB().Raw(`SELECT sb.name, ssc.first_ca, ssc.second_ca, ssc.first_exam, ssc.total_first, ssc.grade_first, ssc.remark_first
	FROM students as s
	INNER JOIN student_subject_classes as ssc ON s.id = ssc.student_id
	INNER JOIN subject_classes as sc ON sc.id = ssc.subject_class_id
	INNER JOIN subjects as sb on sb.id = sc.subject_id
	WHERE s.id = ?`, studentID).Rows()

	defer rows.Close()
	for rows.Next() {
		var row SubjectListRow
		db.ScanRows(rows, &row)

		result = append(result, &row)

	}

	return result
}

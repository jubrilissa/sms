package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Student field (Model) defined
type Student struct {
	gorm.Model
	Name     string
	Address  string
	PhoneNo  string
	Email    string
	Religion string
	ClassID  uint
	// Use real Id not text
	ClassText string
	Image     string
	Gender    string
	// TODO: Should note not be somethinglike text
	Notes       string
	DateOfBirth *time.Time
}

// Create a user object
func (student *Student) Create() uint {

	fmt.Println(student, "the student object")

	fmt.Println("Just before running create")
	// GetDB().Create(&student)
	GetDB().Create(&student).Scan(&student)
	fmt.Println("Student Id is ", student.ID)
	studentID := student.ID
	fmt.Println("After running before running create")
	return studentID

	// if err != nil {
	// 	fmt.Print(err)
	// 	panic(err.Error())
	// }

	// fmt.Println("resp:", &resp.)

	//Create new JWT token for the newly registered user
	// claims := GenerateUserClaims(user.ID, user.Email)

	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	// tokenString, _ := token.SignedString([]byte(os.Getenv("PASSPHRASE")))
	// user.Token = tokenString

	// user.Password = "" //delete password

	// response := utils.Message(true, "user has been created")
	// response["user"] = user

}

func GetAllStudents() []*Student {
	students := make([]*Student, 0)

	err := GetDB().Table("students").Find(&students).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return students

}

func GetSingleStudentById(id uint) *Student {
	student := &Student{}
	err := GetDB().Table("students").Where("id = ?", id).First(student).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return student
}

// GetSubjectsClassForStudentByID - Return the subject the student take in a given class
func GetSubjectsClassForStudentByID(id uint) []*StudentSubjectClass {
	studentSubjectClass := make([]*StudentSubjectClass, 0)
	err := GetDB().Table("student_subject_classes").Where("student_id = ?", id).Find(&studentSubjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return studentSubjectClass
}

func GetSubjectsClassForStudentByID2(id uint) []*Subject {
	studentSubjectClass := make([]*Subject, 0)

	err := GetDB().Preload("student_subject_classes").Where("id = ?", id).Find(&studentSubjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return studentSubjectClass
}

// func GetSubjectsDetailsForClass(currentClass string) []*Subject {
// 	subjects := make([]*Subject, 0)

// 	err := GetDB().Preload("SubjectClass").Where("class =?", currentClass).Find(&subjects).Error

// 	// err := GetDB().Table("subjects").Find(&subjects).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil
// 	}

// 	return subjects
// }

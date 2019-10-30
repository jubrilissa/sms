package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Subject field (Model) defined
type Subject struct {
	// TODO: View the original image and find better ways to display the teachers for the subjects and total no of subject taking them and classes
	gorm.Model
	Image string
	Name  string
	// TODO: Change this to a boolean
	IsCompulsory string
	// TODO: We should enforce that it does a cascasde delete
	SubjectClass []SubjectClass
}

// Create a user object
func (subject *Subject) Create() uint {

	fmt.Println(subject, "the subject object")

	fmt.Println("Just before running create")
	GetDB().Create(&subject).Scan(&subject)
	subjectID := subject.ID

	return subjectID

}

func GetAllSubjects() []*Subject {
	subjects := make([]*Subject, 0)

	err := GetDB().Table("subjects").Find(&subjects).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return subjects
}

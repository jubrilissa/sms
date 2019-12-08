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

// type Deployment struct {
//     gorm.Model
//     Name        string `gorm:"unique_index:idx_name"`
//     RestAPIUser string
//     RestAPIPass string
//     Servers     []Server
//     model       *Model
// }

// type Server struct {
//     gorm.Model
//     DeploymentID uint
//     Hostname     string `gorm:"unique_index:idx_hostname"`
//     RestPort     string
//     Version      string
// }

// func (m *Model) GetDeployments() ([]Deployment, error) {
//     deployments := []Deployment{}
//     return deployments, m.db.Preload("Servers").Find(&deployments).Error
// }

func GetAllSubjectsDetails() []*Subject {
	subjects := make([]*Subject, 0)

	err := GetDB().Preload("SubjectClass").Find(&subjects).Error

	// err := GetDB().Table("subjects").Find(&subjects).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return subjects
}

func GetSubjectsDetailsForClass(currentClass string) []*Subject {
	subjects := make([]*Subject, 0)

	err := GetDB().Preload("SubjectClass").Where("class =?", currentClass).Find(&subjects).Error

	// err := GetDB().Table("subjects").Find(&subjects).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return subjects
}

func GetSubjectById(id uint) *Subject {
	subject := &Subject{}
	err := GetDB().Table("subjects").Where("id = ?", id).First(subject).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return subject
}

func GetSubjectBySubjectClassId(id uint) *Subject {

	subjectClass := &SubjectClass{}
	err := GetDB().Table("subject_classes").Where("id = ?", id).First(subjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}

	subjects := GetSubjectById(subjectClass.SubjectID)

	return subjects
}

func GetSubjectBySubjectClassId2(id uint) *Subject {

	subject := &Subject{}
	// db.Preload("Courses", "domain_id = ?", domainID).Table("appointments").Where("id = ?", appointmentID)
	err := GetDB().Preload("SubjectClass", "id = ?", id).Find(&subject).Error
	// err := GetDB().Table("subject_classes").Where("id = ?", id).First(subjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}

	// subjects := GetSubjectById(subjectClass.SubjectID)

	return subject
}

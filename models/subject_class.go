package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Class field (Model) defined
type SubjectClass struct {
	gorm.Model
	Class        string
	Teacher      int
	IsCompulsory bool
	UserID       uint
	SubjectID    uint
}

// Create a user object
func (subjectClass *SubjectClass) Create() {

	fmt.Println(subjectClass, "the class object")

	fmt.Println("Just before running create")
	GetDB().Create(&subjectClass)

	// if err != nil {
	// 	fmt.Print(err)
	// 	panic(err.Error())
	// }

	// fmt.Println("resp:", &resp.)

	fmt.Println("After running before running create")

	//Create new JWT token for the newly registered user
	// claims := GenerateUserClaims(user.ID, user.Email)

	// token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	// tokenString, _ := token.SignedString([]byte(os.Getenv("PASSPHRASE")))
	// user.Token = tokenString

	// user.Password = "" //delete password

	// response := utils.Message(true, "user has been created")
	// response["user"] = user

}

// func GetSubjectsForClass(currentClass string) *SubjectClass {
// 	subjectClass := &SubjectClass{}
// 	err := GetDB().Table("subject_classes").Where("class = ?", currentClass).First(subjectClass).Error
// 	if err != nil || err == gorm.ErrRecordNotFound {
// 		return nil
// 	}
// 	return subjectClass
// }

func GetSubjectsForClass(currentClass string) []*SubjectClass {

	subjectClass := make([]*SubjectClass, 0)
	// db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
	// 	return db.Order("orders.amount DESC")
	// }).Find(&users)
	err := GetDB().Table("subject_classes").Where("class = ?", currentClass).Find(&subjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return subjectClass
}

func UpdateSubjectClassTeacher(id uint, teacher uint) *SubjectClass {
	subjectClass := &SubjectClass{}
	err := GetDB().Table("subject_classes").Where("id = ?", id).First(subjectClass).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}

	GetDB().Model(&subjectClass).Update("user_id", teacher)
	return subjectClass

}

// func GetSubjectClassDetails() []*SubjectClass {
// 	subjectClass := make([]*SubjectClass, 0)

// }

// err := GetDB().Preload("SubjectClass").Where("class =?", currentClass).Find(&subjects).Error

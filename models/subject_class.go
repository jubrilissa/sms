package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Class field (Model) defined
type SubjectClass struct {
	gorm.Model
	Subject      string
	Class        string
	Teacher      int
	IsCompulsory bool
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

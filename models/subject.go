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
}

// Create a user object
func (subject *Subject) Create() {

	fmt.Println(subject, "the subject object")

	fmt.Println("Just before running create")
	GetDB().Create(&subject)

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

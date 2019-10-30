package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Class field (Model) defined
type Class struct {
	gorm.Model
	Name             string
	ClassCoordinator string
	Students         []Student
}

// Create a user object
func (class *Class) Create() {

	fmt.Println(class, "the class object")

	fmt.Println("Just before running create")
	GetDB().Create(&class)

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

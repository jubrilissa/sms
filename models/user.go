package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name     string
	Email    string
	PhoneNo  string
	Password string
	Role     string
	// TODO: This is currently representing the arrays of subject in a class the user(teacher) teaches
	SubjectClass []SubjectClass
}

func (user *User) Create() uint {
	fmt.Println(user, " The user object")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	GetDB().Create(&user).Scan(&user)
	userID := user.ID
	return userID

}

func Login(email string, password string) bool {
	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		return false
	}

	fmt.Println("user hash", user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		fmt.Println("The password hash compairison failed")
		return false
	}

	return true

}

func GetAllUserByRole(role string) []*User {
	users := make([]*User, 0)

	err := GetDB().Table("users").Where("role = ?", role).Find(&users).Error

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return users

}

func GetAllUsers() []*User {
	users := make([]*User, 0)

	err := GetDB().Table("users").Find(&users).Error

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return users

}

// GetUserByID - Get the user object for a give id
func GetUserByID(id uint) *User {
	user := &User{}
	err := GetDB().Table("users").Where("id = ?", id).First(user).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return user
}

// GetUserByEmail - Get the user object for a given email
func GetUserByEmail(email string) *User {
	user := &User{}
	err := GetDB().Table("users").Where("email = ?", email).First(user).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return user
}

// Validate incoming user details...
// func (user *User) Validate() (map[string]interface{}, bool) {

// 	if !strings.Contains(user.Email, "@") {
// 		// return utils.Message(false, "Email address is required"), false
// 		fmt.Println("Email address is required")
// 	}

// 	if len(user.Password) < 6 {
// 		// return utils.Message(false, "Password is required"), false
// 		fmt.Println("Password is required")
// 	}

// 	//Email must be unique
// 	temp := &User{}

// 	//check for errors and duplicate emails
// 	err := GetDB().Table("users").Where("email = ?", user.Email).First(temp).Error
// 	if err != nil && err != gorm.ErrRecordNotFound {
// 		return utils.Message(false, "Connection error. Please retry"), false
// 	}
// 	if temp.Email != "" {
// 		return utils.Message(false, "Email address already in use by another user."), false
// 	}

// 	return utils.Message(false, "Requirement passed"), true
// }

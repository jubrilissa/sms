package models

import (
	"fmt"

	// "os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // get the gorm postgres dialect
)

const (
	port = 5432
)

var db *gorm.DB

func init() {

	// password := os.Getenv("PGPASSWORD")
	// host := os.Getenv("PGHOST")
	// user := os.Getenv("PGUSER")
	// dbname := os.Getenv("PGDBNAME")
	password := ""
	host := "localhost"
	user := "masterp"
	dbname := "sms"

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", host, port, user, dbname, password)

	fmt.Println(psqlInfo)
	conn, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Print(err)
		panic(err.Error())
	}

	db = conn
	// defer db.Close()

	db.Debug().AutoMigrate(&Student{}, &Class{}, &SubjectClass{}, &Subject{}, &User{}, &StudentSubjectClass{})
	// db.Create(&Student{Name: "isere", Address: "Test addres", PhoneNo: "08034246506", Religion: "Tester"})
}

// GetDB function defined to return DB instance
func GetDB() *gorm.DB {
	return db
}

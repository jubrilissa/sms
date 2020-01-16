package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type FeesPayment struct {
	gorm.Model
	StudentID  uint
	Amount     float64
	IsComplete bool
	// necessary to track who made entered the payment
	UserID uint
}

func (feesPayment *FeesPayment) Create() {

	fmt.Println(feesPayment, "the class object")

	fmt.Println("Just before running create")
	GetDB().Create(&feesPayment)
}

func GetPaymentsForStudent(StudentID uint) []*FeesPayment {

	studentFeePayments := make([]*FeesPayment, 0)
	// db.Preload("Orders", func(db *gorm.DB) *gorm.DB {
	// 	return db.Order("orders.amount DESC")
	// }).Find(&users)
	err := GetDB().Table("fees_payments").Where("student_id = ?", StudentID).Find(&studentFeePayments).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return studentFeePayments
}

func GetFeePaidById(id uint) *FeesPayment {
	feePaid := &FeesPayment{}
	err := GetDB().Table("fees_payments").Where("id = ?", id).First(feePaid).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		return nil
	}
	return feePaid
}

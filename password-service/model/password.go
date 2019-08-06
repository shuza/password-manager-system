package model

import "github.com/jinzhu/gorm"

type Password struct {
	gorm.Model
	UserId      uint   `json:"user_id" gorm:"column:user_id;not null"`
	AccountName string `json:"account_name" gorm:"column:account_name;not null"`
	Email       string `json:"email" gorm:"column:email"`
	Username    string `json:"username" gorm:"column:username"`
	Password    string `json:"password" gorm:"column:password;not null"`
}

package model

import "github.com/jinzhu/gorm"

type Password struct {
	gorm.Model
	UserId      uint   `json:"user_id"`
	AccountName string `json:"account_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

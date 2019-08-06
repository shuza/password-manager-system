package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Email        string `json:"email" gorm:"text,UNIQUE"`
	Fullname     string `json:"fullname"`
	Password     string `json:"password"`
	BusinessName string `json:"business_name"`
}

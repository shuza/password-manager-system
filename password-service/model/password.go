package model

import "github.com/jinzhu/gorm"

type Password struct {
	gorm.Model
	UserId      int    `json:"user_id" gorm:"not null"`
	AccountName string `json:"account_name" gorm:"not null"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	Password    string `json:"password" gorm:"not null"`
}

func (p *Password) IsValid() bool {
	return p.Password != "" && (p.Email != "" || p.Username != "") && p.AccountName != ""
}

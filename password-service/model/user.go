package model

type User struct {
	Id           int    `json:"ID"`
	Email        string `json:"email"`
	Fullname     string `json:"fullname"`
	Password     string `json:"password"`
	BusinessName string `json:"business_name"`
}

package model

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"full_name" db:"full_name"`
	Password     string `json:"password"`
	BusinessName string `json:"business_name" db:"business_name"`
}

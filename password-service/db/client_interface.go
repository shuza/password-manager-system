package db

import "password-service/model"

type IRepository interface {
	Init() error
	Save(model interface{}) error
	Delete(model interface{}) error
	GetByUserId(userId uint) ([]model.Password, error)
	GetById(id int, userId int) (model.Password, error)
	GetByUsernameAndPassword(username string, password string) (model.Password, error)
	Close()
}

var Client IRepository

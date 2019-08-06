package db

import "user-service/model"

type IRepository interface {
	Init() error
	Save(model interface{}) error
	GetByEmail(email string) (model.User, error)
	Close()
}

var Client IRepository

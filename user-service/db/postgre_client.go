package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
	"os"
	"user-service/model"
)

type UserRepository struct {
	conn *gorm.DB
}

func (repo *UserRepository) Init() error {
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	log.Infoln("connectionStr :  " + connectionStr)
	db, err := gorm.Open("postgres", connectionStr)
	repo.conn = db
	return err
}

func (repo *UserRepository) Save(model interface{}) error {
	repo.conn.AutoMigrate(model)
	db := repo.conn.Save(model)
	return db.Error
}

func (repo *UserRepository) GetByEmail(email string) (model.User, error) {
	var user model.User
	err := repo.conn.Where("email = ?", email).First(&user).Error
	return user, err
}
func (repo *UserRepository) Get(user model.User) (model.User, error) {
	var data model.User
	err := repo.conn.First(&data, user).Error
	return data, err
}

func (repo *UserRepository) Close() {
	repo.conn.Close()
}

package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"os"
	"password-service/model"
)

type PostgresClient struct {
	conn *gorm.DB
}

func (c *PostgresClient) Init() error {
	connectionStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	log.Infoln("connectionStr :  " + connectionStr)
	db, err := gorm.Open("postgres", connectionStr)
	c.conn = db
	return err
}

func (c *PostgresClient) Save(model interface{}) error {
	c.conn.AutoMigrate(&model)
	return c.conn.Save(model).Error
}

func (c *PostgresClient) GetByUserId(userId uint) ([]model.Password, error) {
	result := make([]model.Password, 0)
	if err := c.conn.Where("user_id = ?", userId).Find(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

func (c *PostgresClient) GetById(id uint) (model.Password, error) {
	var password model.Password
	err := c.conn.Where("id = ?", id).First(&password).Error
	return password, err
}

func (c *PostgresClient) GetByUsernameAndPassword(username string, password string) (model.Password, error) {
	var result model.Password
	err := c.conn.Where("username = ? AND password = ?", username, password).First(&result).Error
	return result, err
}

func (c *PostgresClient) Delete(model interface{}) error {
	return c.conn.Delete(model).Error
}

func (c *PostgresClient) Close() {
	c.conn.Close()
}

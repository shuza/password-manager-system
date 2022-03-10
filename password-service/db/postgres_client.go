package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	db.AutoMigrate(model.Password{})
	return err
}

func (c *PostgresClient) Save(data interface{}) error {
	c.conn.AutoMigrate(data)
	return c.conn.Save(data).Error
}

func (c *PostgresClient) GetByUserId(userId uint) ([]model.Password, error) {
	result := make([]model.Password, 0)
	if err := c.conn.Where("user_id = ?", userId).Find(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

func (c *PostgresClient) GetById(id int, userId int) (model.Password, error) {
	var password model.Password
	err := c.conn.Where("id = ? AND user_id = ?", id, userId).First(&password).Error
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

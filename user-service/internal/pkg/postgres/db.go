package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	defaultMacIleConn      = 10
	defaultMaxOpenConn     = 10
	defaultConnMaxLifetime = 30 * time.Minute
	defaultMigrationPath   = "file://./migrations"
)

type Config struct {
	Host                     string
	Port                     string
	User                     string
	Password                 string
	Name                     string
	MaxIdleConn              int
	MaxOpenConn              int
	ConnMacLifeTimeInMinutes int
	MigrationPath            string
}

func New(conf *Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", conf.Url())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(conf.maxIdleConn())
	db.SetMaxOpenConns(conf.maxOpenConn())
	db.SetConnMaxLifetime(conf.connMaxLifeTime())

	return db, nil
}

func (c *Config) migrationPath() string {
	if c.MigrationPath == "" {
		return defaultMigrationPath
	}
	return c.MigrationPath
}

func (c *Config) Url() string {
	connectionUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name)
	log.Println("DB connection URL = ", connectionUrl)
	return connectionUrl
}

func (c *Config) maxIdleConn() int {
	if c.MaxIdleConn == 0 {
		return defaultMacIleConn
	}
	return c.MaxIdleConn
}

func (c *Config) maxOpenConn() int {
	if c.MaxOpenConn == 0 {
		return defaultMaxOpenConn
	}
	return c.MaxOpenConn
}

func (c *Config) connMaxLifeTime() time.Duration {
	if c.ConnMacLifeTimeInMinutes == 0 {
		return defaultConnMaxLifetime
	}
	return time.Duration(c.ConnMacLifeTimeInMinutes) * time.Minute
}

func RunDatabaseMigration(conf *Config) error {
	m, err := migrate.New(conf.migrationPath(), conf.Url())
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}

func RollbackLatestMigration(config *Config) error {
	m, err := migrate.New(config.migrationPath(), config.Url())
	if err != nil {
		return err
	}

	if err = m.Steps(-1); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}

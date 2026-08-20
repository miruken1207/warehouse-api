package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (dbCfg *DBConfig) Load() error {
	godotenv.Load()

	dbCfg.Host = os.Getenv("DB_HOST")
	dbCfg.Port = os.Getenv("DB_PORT")
	dbCfg.User = os.Getenv("DB_USER")
	dbCfg.Password = os.Getenv("DB_PASSWORD")
	dbCfg.Name = os.Getenv("DB_NAME")

	if dbCfg.Host == "" || dbCfg.Port == "" || dbCfg.User == "" || dbCfg.Password == "" || dbCfg.Name == "" {
		return errors.New("empty config field")
	}

	return nil
}

func (dbCfg *DBConfig) DSN() string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbCfg.User, dbCfg.Password, dbCfg.Host, dbCfg.Port, dbCfg.Name)
	return dsn
}

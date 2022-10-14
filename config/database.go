package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var pg *gorm.DB // main postgres connection

func ConnectDB() (err error) {
	var (
		host         = os.Getenv("DB_HOST")
		port         = os.Getenv("DB_PORT")
		user         = os.Getenv("DB_USER")
		pswd         = os.Getenv("DB_PSWD")
		dbName       = os.Getenv("DB_NAME")
		dbDataSource = "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta"
	)

	dbDataSource = fmt.Sprintf(dbDataSource, host, user, pswd, dbName, port)
	pg, err = gorm.Open(postgres.Open(dbDataSource))

	return err
}

func CloneDB() *gorm.DB {
	return pg.Session(&gorm.Session{
		SkipHooks:              true,
		SkipDefaultTransaction: true,
	})
}

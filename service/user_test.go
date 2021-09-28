package service

import (
	"database/sql"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func init() {
	var db *sql.DB
	var err error
	var DB *gorm.DB
	db, _, err = sqlmock.New()
	if err != nil {
		log.Fatalf("sqlmock failed:", err)
	}

	DB, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Init DB with sqlmock failed:", err)
	}
}

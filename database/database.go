package database

import (
	"log"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/YagoCarballo/kumquat.academy.api/tools"
)

var (
	DB *gorm.DB
)

const DBUrl = "localhost"
const DBName = "golang_test"

func connectWithDB() error {
	var mysqlError, mysqlDownError error
	var rawDatabase gorm.DB
	var uri, dbType string

	dbSettings := tools.GetSettings().Database

	switch strings.ToLower(dbSettings.Type) {
	case "mysql":
		dbType = "mysql"
		uri = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			dbSettings.Mysql.Username,
			dbSettings.Mysql.Password,
			dbSettings.Mysql.Host,
			dbSettings.Mysql.Name,
		)
	default:
		dbType = "sqlite3"
		uri = fmt.Sprintf("%s", dbSettings.Sqlite.Path)
	}

	// Connects to the Database
	if rawDatabase, mysqlError = gorm.Open(dbType, uri); mysqlError != nil {
		return mysqlError
	}

	DB = &rawDatabase

	// Open doesn't open a connection. Validate DSN data:
	mysqlDownError = DB.DB().Ping()
	if mysqlDownError != nil {
		return mysqlDownError
	}

	// Prints the Connection Details
	switch strings.ToLower(dbSettings.Type) {
	case "mysql":
		log.Printf("Connected to MySQL { server: %s, db: %s }\n", dbSettings.Mysql.Host, dbSettings.Mysql.Name)
	default:
		log.Printf("Connected to SQLite { db: %s }\n", dbSettings.Sqlite.Path)
	}

	return nil
}

func InitDatabase() (error, *gorm.DB) {
	err := connectWithDB()

	return err, DB
}

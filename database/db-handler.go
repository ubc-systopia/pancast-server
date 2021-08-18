package database

/*
	This file contains database connection logic.
 */

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitDatabaseConnection() *sql.DB {
	if err := loadEnvironmentVariables(); err != nil {
		log.Fatal(err)
	}
	accessParams := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWD") +
		"@tcp(" + os.Getenv("DB_ADDR") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME")
	db, err := sql.Open("mysql", accessParams)
	if err != nil {
		log.Fatal(err)
	}
	// verify that the connection has actually been established
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func loadEnvironmentVariables() error {
	var err error
	if flag.Lookup("test.v") == nil {
		err = godotenv.Load()
	} else {
		err = godotenv.Load("./../.env")
	}
	if err != nil {
		return err
	}
	return nil
}

package server

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitDatabaseConnection() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	accessParams := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWD") +
		"@tcp(localhost:3306)/" + os.Getenv("DB_NAME")
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

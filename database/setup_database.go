package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var db *sql.DB

func SetupDatabase() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
		return err
	}

	password := os.Getenv("MYSQL_PASSWORD")

	db, err = sql.Open("mysql", "root:"+password+"@tcp(127.0.0.1:3306)/neolib")
	if err != nil {
		log.Fatal(err)
		return err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return nil
}

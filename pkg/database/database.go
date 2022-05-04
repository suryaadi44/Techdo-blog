package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectDB(address string, username string, password string, dbName string) *sql.DB {
	log.Printf("[DB] starting database connection process")

	db_string := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, address, dbName)

	db, err := sql.Open("mysql", db_string)
	if err != nil {
		log.Fatalf("[DB] sql open connection fatal error: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("[DB] db ping fatal error: %v", err)
	}

	log.Printf("[DB] database connectionn: established successfully with %s\n", db_string)
	return db
}

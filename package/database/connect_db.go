package database

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() {
	dsn:="root:11042005(127.0.0.1:3306)/nguoi_dung"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}
	log.Println("Connect to database successful")
}

package database

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func ConnectDb() (*sql.DB, error){
	dsn:="root:11042005@tcp(127.0.0.1:3306)/user_manage?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}
	log.Println("Connect to database successful")
	return db,nil
}

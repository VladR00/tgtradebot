package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"

)
const (
	path string = "./sql/sql.db"
)
func OpenDB() (*sql.DB, error){
	db, err := sql.Open("sqlite3", path) 
	if err != nil {
		return nil, fmt.Errorf("Can't open DB: %w", err)
	}
	if err = db.Ping(); err != nil{
		return nil, fmt.Errorf("Can't ping DB: %w", err)
	}
	fmt.Println("DB successfully connecting")
	return db, nil
}

func CreateDB(db *sql.DB) error{
	query, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (chat_id INTEGER PRIMARY KEY, linkname TEXT, username TEXT, balance DECIMAL(15,2), registration_time INTEGER)")
	if err != nil {
		return fmt.Errorf("Can't preparing query for creating table users: %w", err)
	}
	defer query.Close()

	if _, err = query.Exec(); err != nil{
		return fmt.Errorf("Can't executing create table users: %w", err)
	}

	return nil
}


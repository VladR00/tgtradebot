package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"

)
const (
	path string = "gogi"
)
func OpenDB(){
	db, err := sql.Open("sqlite3", path) 
}

func CreateDB(){
	db, err := sql.Open("sqlite3", path)
	query := db.Prepare("CREATE TABLE IF NOT EXIST users (chat_id INTEGER PRIMARY KEY, linkname TEXT, username TEXT, balance DECIMAL(15,2), registration_time INTEGER)")
}
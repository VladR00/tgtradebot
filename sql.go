package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)
const (
	path string = "./sql/sql.db"
)

type User struct{
	ChatID		int64
	LinkName	string
	UserName	string
	Balance 	int64
	Time 		string
}


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

func CreateDB() error{
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (chat_id INTEGER PRIMARY KEY, linkname TEXT, username TEXT, balance DECIMAL(15,2), registration_time INTEGER)")
	if err != nil {
		return fmt.Errorf("Can't preparing query for creating table users: %w", err)
	}
	defer query.Close()

	if _, err = query.Exec(); err != nil{
		return fmt.Errorf("Can't execute create table users: %w", err)
	}

	return nil
}

func InsertNewUsersDB(chatID int64, linkname string, username string) error{
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query, err := db.Prepare("INSERT INTO users (chat_id, linkname, username, balance, registration_time) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("Can't preparing query for insert new users into users: %w", err)
	}
	defer query.Close()

	if _, err = query.Exec(chatID, linkname, username, 0, time.Now().Unix()); err != nil {
		return fmt.Errorf("Can't execute inserting new users into users: %w", err)
	}
	return nil
}

func ReadUserByID(chatID int64) (*User, error){
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := ("SELECT * FROM users WHERE chat_id = ?")

	user := &User{}
	var registrationTime int64

	row := db.QueryRow(query, chatID)
	err = row.Scan(&user.ChatID, &user.LinkName, &user.UserName, &user.Balance, &registrationTime) 
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("User not found while reads user: %w", err)
		}
		return nil, fmt.Errorf("Undefined error while reads user")
	}
	user.Time = time.Unix(registrationTime, 0).Format("2006-01-02 15:04")
	return user, nil
}
//func UpdateUsersDB(chatID int64){}
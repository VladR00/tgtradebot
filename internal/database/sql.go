package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)
var (
	DBpath 		=	"../storage/sql.db"
	UserMap 	=	map[int64]User{}
	StaffMap 	=	map[int64]Staff{}
)

type User struct{
	ChatID			int64
	LinkName		string
	UserName		string
	Balance 		int64
	Time 			string
	CurrentTicket 	int64
}

type Staff struct{
	ChatID			int64
	Admin			bool	
	CurrentTicket 	int64
	LinkName		string
	UserName		string
	TicketClosed	int64
	Rating 			int64
	Time 			string
}

func IsTableExists(db *sql.DB, tableName string) bool {
	query := `SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?;`
	var count int
	err := db.QueryRow(query, tableName).Scan(&count)
	if err != nil {
		fmt.Println(err)
	}

	return count > 0
}

func OpenDB() (*sql.DB, error){
	db, err := sql.Open("sqlite3", DBpath) 
	if err != nil {
		return nil, fmt.Errorf("Can't open DB: %w", err)
	}
	if err = db.Ping(); err != nil{
		return nil, fmt.Errorf("Can't ping DB: %w", err)
	}
	fmt.Println("DB successfully connecting")
	return db, nil
}

func CreateTable(table string) error{
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	var q string
	switch table {
		case "users":
			q = `CREATE TABLE IF NOT EXISTS users (
					chat_id INTEGER PRIMARY KEY,
					linkname TEXT, 
					username TEXT,
					balance DECIMAL(15,2),
					registration_time INTEGER)`
		case "staff":
			q = `CREATE TABLE IF NOT EXISTS staff (
					chat_id INTEGER PRIMARY KEY,
					admin BOOL,
					current_ticket INTEGER,
					linkname TEXT, 
					username TEXT,
					ticket_closed INTEGER,
					rating INTEGER,
					registration_time INTEGER)`
		case "bookkeeping":
			q = `CREATE TABLE IF NOT EXISTS bookkeeping (
				`
	}

	query, err := db.Prepare(q)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Can't preparing query for creating table %s: %w", table, err)
	}
	defer query.Close()

	if _, err = query.Exec(); err != nil{
		fmt.Println(err)
		return fmt.Errorf("Can't execute create table %s: %w", table, err)
	}

	return nil
}

func InsertNewUser(chatID int64, linkname string, username string) error{
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

func InsertNewStaff(chatID int64, admin bool, linkname string, username string) error{
	db, err := OpenDB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	query, err := db.Prepare("INSERT INTO staff (chat_id, admin, current_ticket, linkname, username, ticket_closed, rating, registration_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Can't preparing query for insert new staff into staff: %w", err)
	}
	defer query.Close()

	if _, err = query.Exec(chatID, admin, 0, linkname, username, 0, 0, time.Now().Unix()); err != nil {
		fmt.Println(err)
		return fmt.Errorf("Can't execute inserting new staff into staff: %w", err)
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
		return nil, fmt.Errorf("Undefined error while reads user: %w", err)
	}
	user.Time = time.Unix(registrationTime, 0).Format("2006-01-02 15:04")
	return user, nil
}

func ReadStaffByID(chatID int64) (*Staff, error){
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := ("SELECT * FROM staff WHERE chat_id = ?")

	staff := &Staff{}
	var registrationTime int64

	row := db.QueryRow(query, chatID)
	err = row.Scan(&staff.ChatID, &staff.Admin, &staff.CurrentTicket, &staff.LinkName, &staff.UserName, &staff.TicketClosed, &staff.Rating, &registrationTime) 
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("Staff not found while reads staff: %w", err)
		}
		return nil, fmt.Errorf("Undefined error while reads staff: %w", err)
	}
	staff.Time = time.Unix(registrationTime, 0).Format("2006-01-02 15:04")
	return staff, nil
}

func UpdateUsersDB(chatID int64, topUp int64) error{
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := ("UPDATE users SET balance = balance + ? WHERE chat_id = ?")

	result, err := db.Exec(query, topUp, chatID)
	if err != nil {
		return fmt.Errorf("Can't update balance from users: %w", err)
	}

	rowsAffected, err := result.RowsAffected() 
	if err != nil {
		return fmt.Errorf("Can't update balance from users while checking RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ChatID: %d; Isn't updated, user not found (maybe).", chatID)
	}

	return nil
}
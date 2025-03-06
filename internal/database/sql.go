package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)
var (
	DBpath 		=	"../storage/sql.db"
	UserMap 	=	map[int64]User{}
	StaffMap 	=	map[int64]Staff{}
	TicketMap	=	map[int64]Ticket{}
)

type User struct{
	ChatID			int64
	LinkName		string
	UserName		string
	Balance 		int64
	Time 			int64
	CurrentTicket 	int64
	Language 		string
}

type Staff struct{
	ChatID			int64
	Admin			int32	
	CurrentTicket 	int64
	LinkName		string
	UserName		string
	TicketClosed	int64
	Rating 			int64
	Time 			int64
}
type Ticket struct{
	TicketID		int64
	ChatID			int64
	SupChatID		int64
	LinkName		string
	SupLinkName		string
	UserName		string
	SupUserName		string
	Time 			int64
	ClosingTime 	int64
	Language		string
	Status			string
}
type TicketMessage struct{
	TicketID	int64
	Support 	int32
	ChatID		int64
	UserName	string
	MessageID	int
	Time		int64
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
					current_ticket INTEGER,
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
		case "tickets":
			q = `CREATE TABLE IF NOT EXISTS tickets (
				id INTEGER PRIMARY KEY,
				registration_time INTEGER,
				closing_time INTEGER,
				chat_id INTEGER,
				sup_chat_id INTEGER,
				linkname TEXT, 
				username TEXT,
				sup_linkname TEXT, 
				sup_username TEXT,
				prefered_language TEXT,
				status TEXT)`
		case "tickets_messages":
			q = `CREATE TABLE IF NOT EXISTS tickets_messages (
				ticket_id INTEGER,	
				sup INTEGER,
				chat_id INTEGER,
				username TEXT,
				message_id INTEGER,
				time INTEGER)`
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
	if (table == "tickets_messages"){
		if _, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_ticket_id ON tickets_messages(ticket_id)"); err != nil{
			return fmt.Errorf("Can't execute create index in  %s: %w", table, err)
		}
	}

	return nil
}

func (s *User) InsertNew() error{
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query, err := db.Prepare("INSERT INTO users (chat_id, linkname, username, balance, current_ticket, registration_time) VALUES (?, ?, ?, ?, ? ,?)")
	if err != nil {
		return fmt.Errorf("Can't preparing query for insert new users into users: %w", err)
	}
	defer query.Close()

	if _, err = query.Exec(s.ChatID, s.LinkName, s.UserName, s.Balance, s.CurrentTicket, s.Time); err != nil { 
		return fmt.Errorf("Can't execute inserting new users into users: %w", err)
	}
	return nil
}

func (s *Staff) InsertNew() error{
	db, err := OpenDB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	query, err := db.Prepare("INSERT INTO staff (chat_id, admin, linkname, username, current_ticket, ticket_closed, rating, registration_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Can't preparing query for insert new staff into staff: %w", err)
	}
	defer query.Close()

	if _, err = query.Exec(s.ChatID, s.Admin, s.LinkName, s.UserName, s.CurrentTicket, s.TicketClosed, s.Rating, s.Time); err != nil { 
		fmt.Println(err)
		return fmt.Errorf("Can't execute inserting new staff into staff: %w", err)
	}
	return nil
}

func (s *Ticket) InsertNew() error{
	db, err := OpenDB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	query, err := db.Prepare("INSERT INTO tickets (chat_id, sup_chat_id, linkname, sup_linkname, username, sup_username, registration_time, closing_time, prefered_language, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Can't preparing query for insert new ticket into tickets: %w", err)
	}
	defer query.Close()

	r, err := query.Exec(s.ChatID, s.SupChatID, s.LinkName, s.SupLinkName, s.UserName, s.SupUserName, s.Time, s.ClosingTime, s.Language, s.Status)
	if err != nil { 
		fmt.Println(err)
		return fmt.Errorf("Can't execute inserting new ticket into tickets: %w", err)
	}
	fmt.Println(r.RowsAffected())
	return nil
}

func (s *TicketMessage) InsertNew() error{
	db, err := OpenDB()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	query, err := db.Prepare("INSERT INTO tickets_messages (ticket_id, sup, chat_id, username, message_id, time) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("Can't preparing query for insert new message into tickets_messages: %w", err)
	}
	defer query.Close()

	if _, err = query.Exec(s.TicketID, s.Support, s.ChatID, s.UserName, s.MessageID, s.Time); err != nil { 
		fmt.Println(err)
		return fmt.Errorf("Can't execute inserting new message into tickets_messages: %w", err)
	}
	return nil
}

func (s *User) Update() error{
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := (`UPDATE users 
			   SET balance = ?, current_ticket = ?
			   WHERE chat_id = ?`)

	result, err := db.Exec(query, s.Balance, s.CurrentTicket, s.ChatID)
	if err != nil {
		return fmt.Errorf("Can't update balance from users: %w", err)
	}

	rowsAffected, err := result.RowsAffected() 
	if err != nil {
		return fmt.Errorf("Can't update balance from users while checking RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ChatID: %d; Isn't updated, user not found (maybe).", s.ChatID)
	}

	return nil
}
func (s *Staff) Update() error{
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := (`UPDATE staff 
			   SET current_ticket = ?, username = ?, ticket_closed = ?
			   WHERE chat_id = ?`)

	result, err := db.Exec(query, s.CurrentTicket, s.UserName, s.TicketClosed, s.ChatID)
	if err != nil {
		return fmt.Errorf("Can't update balance from staff: %w", err)
	}

	rowsAffected, err := result.RowsAffected() 
	if err != nil {
		return fmt.Errorf("Can't update staff while checking RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ChatID: %d; Isn't updated, staff not found (maybe).", s.ChatID)
	}

	return nil
}

func (s *Ticket) Update() error{ //by ID
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := (`UPDATE tickets 
			   SET sup_chat_id = ?, sup_linkname = ?,sup_username = ?, closing_time = ?, status = ?
			   WHERE id = ?`)

	result, err := db.Exec(query, s.SupChatID, s.SupLinkName, s.SupUserName, s.ClosingTime, s.Status, s.TicketID)
	if err != nil {
		return fmt.Errorf("Can't update ticket from tickets: %w", err)
	}

	rowsAffected, err := result.RowsAffected() 
	if err != nil {
		return fmt.Errorf("Can't update ticket while checking RowsAffected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("TicketID: %d; Isn't updated, id not found (maybe).", s.TicketID)
	}

	return nil
}


func (s *Staff) MapUpdateOrCreate() {
	StaffMap[s.ChatID] = *s															
}
func (s *User) MapUpdateOrCreate() {
	UserMap[s.ChatID] = *s
}
// func (s *Ticket) MapUpdateOrCreate() {
// 	TicketMap[s.ChatID] = *s
// }
func (s *Staff) MapDelete() {
	delete(StaffMap, s.ChatID)														
}
func (s *User) MapDelete() {
	delete(UserMap, s.ChatID)	
}

func ReadTicketByID(ticketID int64) (*Ticket, error){
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := (`SELECT id, registration_time, closing_time, chat_id, sup_chat_id, linkname, username, sup_linkname, sup_username, prefered_language, status 
              FROM tickets WHERE id = ?`)
	
	ticket := &Ticket{}
	
	row := db.QueryRow(query, ticketID)
	err = row.Scan(
		&ticket.TicketID, 
		&ticket.Time, 
		&ticket.ClosingTime, 
		&ticket.ChatID, 
		&ticket.SupChatID, 
		&ticket.LinkName, 
		&ticket.UserName, 
		&ticket.SupLinkName, 
		&ticket.SupUserName, 
		&ticket.Language, 
		&ticket.Status,)
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("Ticket not found while reads tickets ReadTicketByID: %w", err)
		}
		return nil, fmt.Errorf("Undefined error while reads tickets ReadTicketByID: %w", err)
	}
	return ticket, nil 
}

func ReadOpenTicketByUserID(chatID int64) (*Ticket, error){
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := (`SELECT id, registration_time, closing_time, chat_id, sup_chat_id, linkname, username, sup_linkname, sup_username, prefered_language, status 
              FROM tickets WHERE chat_id = ? AND status != ?`)
	
	ticket := &Ticket{}
	
	row := db.QueryRow(query, chatID, "Closed")
	err = row.Scan(
		&ticket.TicketID, 
		&ticket.Time, 
		&ticket.ClosingTime, 
		&ticket.ChatID, 
		&ticket.SupChatID, 
		&ticket.LinkName, 
		&ticket.UserName, 
		&ticket.SupLinkName, 
		&ticket.SupUserName, 
		&ticket.Language, 
		&ticket.Status,)
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("Ticket not found while reads tickets ReadOpenTicketByUserID: %w", err)
		}
		return nil, fmt.Errorf("Undefined error while reads tickets ReadOpenTicketByUserID: %w", err)
	}
	return ticket, nil 
}

func ReadUserByID(chatID int64) (*User, error){
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := (`SELECT * FROM users 
				WHERE chat_id = ?`)

	user := &User{}

	row := db.QueryRow(query, chatID)
	err = row.Scan(&user.ChatID, &user.LinkName, &user.UserName, &user.Balance, &user.CurrentTicket, &user.Time) 
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("User not found while reads user: %w", err)
		}
		return nil, fmt.Errorf("Undefined error while reads user: %w", err)
	}
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

	row := db.QueryRow(query, chatID)
	err = row.Scan(&staff.ChatID, &staff.Admin, &staff.CurrentTicket, &staff.LinkName, &staff.UserName, &staff.TicketClosed, &staff.Rating, &staff.Time) 
	if err != nil {
		if err == sql.ErrNoRows{
			return nil, fmt.Errorf("Staff not found while reads staff: %w", err)
		}
		return nil, fmt.Errorf("Undefined error while reads staff: %w", err)
	}
	return staff, nil //staff.Time = time.Unix(registrationTime, 0).Format("2006-01-02 15:04")
}

func OutputStaffWithCurrTicketNull() ([]*Staff, error){
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := ("SELECT * FROM staff WHERE current_ticket = ?")

	var stafflist []*Staff

	rows, err := db.Query(query, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		staff := &Staff{}
		err := rows.Scan(&staff.ChatID, &staff.Admin, &staff.CurrentTicket, &staff.LinkName, &staff.UserName, &staff.TicketClosed, &staff.Rating, &staff.Time)
		if err != nil {
			if err == sql.ErrNoRows{
				fmt.Println("Staff not found while reads staff:",err)
			}
			fmt.Println("Undefined error while reads staff:",err)
		}
		stafflist = append(stafflist, staff)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stafflist, nil //staff.Time = time.Unix(registrationTime, 0).Format("2006-01-02 15:04")
}

func ReadAllMessagesByTicketID(ticketid int64) ([]*TicketMessage, error){
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := ("SELECT * FROM tickets_messages WHERE ticket_id = ?")

	var messagelist []*TicketMessage

	rows, err := db.Query(query, ticketid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next(){
		message := &TicketMessage{}
		err := rows.Scan(&message.TicketID, &message.Support, &message.ChatID, &message.UserName, &message.MessageID, &message.Time)
		if err != nil {
			if err == sql.ErrNoRows{
				fmt.Println("Message not found while reads tickets_messages:",err)
			}
			fmt.Println("Undefined error while reads tickets_messages:",err)
		}
		messagelist = append(messagelist, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messagelist, nil
}
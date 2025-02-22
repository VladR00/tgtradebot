package main

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error){
	db, err := sql.Open(driverName:"sqlite3", path)
	if err != nil {
		return nil, fmt.Erorrf("Can't open DataBase: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Erorrf("Can't CONNECT to DataBase: %w", err)
	} 
	return *Storage{db: db}, nil
}

func Save(s *Storage) Save(p *storage.Page) error {
	q := `INSERT INTO pages (name, cid) VALUES(?,?)`
	if _, err := s.db.ExecContext(ctx, q, p.Name, p,CID); err != nil{
		return fmt.Errorf("Can't save page: %w", err)
	}
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXIST pages (name TEXT, chatId INTEGER)`
	if _, err := s.db.ExecContext(ctx, q); err != nil{
		return fmt.Errorf("Can't create table: %w", err)
	}
}
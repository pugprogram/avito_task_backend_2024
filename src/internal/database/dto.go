package database

import "database/sql"

type Database struct {
	db *sql.DB
}

func NewDataBase(db *sql.DB) Database {
	return Database{
		db: db,
	}
}

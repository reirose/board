package main

import (
	"database/sql"
)

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	sqlStmt := `
    create table if not exists posts (id integer not null primary key autoincrement, 
		content text, 
		published_at text, 
		parent_id integer not null default 0);
    `

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	sqlStmt = `create table if not exists users (id integer not null primary key autoincrement,
		user_id text,
		password text, 
		role text,
		token text)`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}

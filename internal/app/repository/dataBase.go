package dataBase

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

const UserTable = `
	CREATE TABLE IF NOT EXISTS User(
	login TEXT UNIQUE,
	password TEXT
	);
	`
const SQLInsert = `INSERT INTO User(login, password) VALUES (?, ?)`

type DataBase struct {
	db *sql.DB
}

func New() *DataBase {
	return &DataBase{}
}

func (myDB *DataBase) StartDB() error {
	var err error
	myDB.db, err = sql.Open("sqlite", "./Users.db")
	if err != nil {
		return err
	}

	myDB.db.SetMaxOpenConns(1)

	myDB.db.Exec(UserTable)

	return nil
}

func (myDB *DataBase) InsertToDB(login string, password string) error {
	_, err := myDB.db.Exec(SQLInsert, login, password)
	if err != nil {
		return err
	}
	fmt.Println(login, password)
	return nil
}

func (myDB *DataBase) CloseDB() {
	myDB.db.Close()
}

package dataBase

import (
	"database/sql"
	"errors"
	_ "modernc.org/sqlite"
)

const UserTable = `
	CREATE TABLE IF NOT EXISTS User(
	id INTEGER PRIMARY KEY,
	login TEXT UNIQUE,
	password TEXT
	);
	`
const SQLInsert = `INSERT INTO User(login, password) VALUES (?, ?)`

const SQLSelectPassword = `SELECT password FROM User WHERE login = (?)`

const SQLSelectID = `SELECT id FROM User WHERE login = (?)`

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

func (myDB *DataBase) InsertInDB(login string, password string) error {
	_, err := myDB.db.Exec(SQLInsert, login, password)
	if err != nil {
		return err
	}
	return nil
}

func (myDB *DataBase) SearchInDB(login string) (string, error) {
	var DBpass string
	err := myDB.db.QueryRow(SQLSelectPassword, login).Scan(&DBpass)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return DBpass, nil
}

func (myDB *DataBase) GetIdFromLogin(login string) (int, error) {
	var id int
	err := myDB.db.QueryRow(SQLSelectID, login).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (myDB *DataBase) CloseDB() {
	myDB.db.Close()
}

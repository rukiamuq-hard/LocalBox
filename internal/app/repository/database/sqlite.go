package dataBase

import (
	"database/sql"
	"errors"

	"Umbrella/internal/app/repository/database/models"

	_ "modernc.org/sqlite"
)

const UserTable = `
	CREATE TABLE IF NOT EXISTS User(
	id INTEGER PRIMARY KEY,
	login TEXT UNIQUE,
	password TEXT
	);
`

const UploadTable = `
	CREATE TABLE IF NOT EXISTS Upload(
	id INTEGER PRIMARY KEY,
	original_name TEXT,
	stored_name TEXT,
	upload_date TEXT,
	size INTEGER,
	uploader_id TEXT
	)
`

const SQLInsert = `INSERT INTO User(login, password) VALUES (?, ?)`

const SQLSelectPassword = `SELECT password FROM User WHERE login = (?)`

const SQLSelectID = `SELECT id FROM User WHERE login = (?)`

const SQLInsertFile = `INSERT INTO Upload(original_name, stored_name, upload_date, size, uploader_id) VALUES (?, ?, ?, ?, ?)`

const SQLGetFiles = `SELECT id, original_name, stored_name, upload_date, size, uploader_id FROM Upload`

const SQLDownloadFiles = `SELECT id, original_name, stored_name, upload_date, size, uploader_id FROM Upload WHERE id = (?)`

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
	myDB.db.Exec(UploadTable)

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

func (myDB *DataBase) StoreFile(fileName string, storeFileName string, dateTime string, size int64, uploader_id string) error {
	_, err := myDB.db.Exec(SQLInsertFile, fileName, storeFileName, dateTime, size, uploader_id)
	if err != nil {
		return err
	}
	return nil
}

func (myDB *DataBase) GetFile() ([]models.UploadedFiles, error) {
	rows, err := myDB.db.Query(SQLGetFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []models.UploadedFiles{}

	for rows.Next() {
		var f models.UploadedFiles
		err := rows.Scan(&f.ID, &f.Original_name, &f.Stored_name, &f.Upload_date, &f.Size, &f.Uploader_id)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func (myDB *DataBase) DownloadFile(id string) (models.UploadedFiles, error) {
	var f models.UploadedFiles
	err := myDB.db.QueryRow(SQLDownloadFiles, id).Scan(&f.ID, &f.Original_name, &f.Stored_name, &f.Upload_date, &f.Size, &f.Uploader_id)
	if err != nil {
		return models.UploadedFiles{}, err
	}
	return f, nil
}

func (myDB *DataBase) CloseDB() {
	myDB.db.Close()
}

package service

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"Umbrella/internal/app/repository/SQLite/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserDB interface {
	InsertInDB(login string, password string) error
	SearchInDB(login string) (string, error)
	GetIdFromLogin(login string) (int, error)
	StoreFile(fileName string, storeFileName string, dateTime string, size int64, uploader_id string) error
	GetFile() ([]models.UploadedFiles, error)
	DownloadFile(id string) (models.UploadedFiles, error)
}

type RedisDB interface {
	SetKeyValue(key string, value int) error
	GetValue(key string) (string, error)
}

type Service struct {
	db  UserDB
	rdb RedisDB
}

func New(uDB UserDB, rDB RedisDB) *Service {
	return &Service{db: uDB, rdb: rDB}
}

func (s *Service) RegisterUser(login string, password string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.db.InsertInDB(login, string(hashedPass))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) LoginUser(login string, password string) error {
	hashedPass, err := s.db.SearchInDB(login)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	if err != nil {
		return errors.New("login: password is wrong")
	}
	return nil
}

func (s *Service) GetIdFromDB(login string) (int, error) {
	id, err := s.db.GetIdFromLogin(login)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) RedisSetKeyValue(key string, value int) error {
	err := s.rdb.SetKeyValue(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) RedisGetValue(key string) (string, error) {
	val, err := s.rdb.GetValue(key)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (s *Service) MakeUUID() string {
	return uuid.New().String()
}

func (s *Service) StoreFileToDB(r io.Reader, fileName string, storeFileName string, size int64, uploader_id string) error {
	Ext := filepath.Ext(fileName)
	ServCreateFile, err := os.Create("./uploads/" + storeFileName + Ext)
	if err != nil {
		return err
	}
	defer ServCreateFile.Close()

	if _, err = io.Copy(ServCreateFile, r); err != nil {
		return err
	}

	TimeIs := time.Now().String()
	err = s.db.StoreFile(fileName, storeFileName+Ext, TimeIs, size, uploader_id)
	if err != nil {
		return err
	}
	fmt.Println("File uploaded: ", fileName)
	return nil
}

func (s *Service) GetFilesFromDB() ([]map[string]any, error) {
	files, err := s.db.GetFile()
	if err != nil {
		return nil, err
	}

	result := make([]map[string]any, 0, len(files))
	for _, f := range files {
		result = append(result, map[string]any{
			"id":   f.ID,
			"name": f.Original_name,
			"type": "file",
			"date": f.Upload_date,
			"size": f.Size,
		})
	}
	return result, nil
}

func (s *Service) DownloadFile(id string) (models.UploadedFiles, error) {
	return s.db.DownloadFile(id)
}

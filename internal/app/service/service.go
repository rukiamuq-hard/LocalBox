package service

import (
	"Umbrella/internal/app/repository/database/models"
	"fmt"
	"github.com/google/uuid"
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
	fmt.Println("[SERVICE]: Initialized")
	return &Service{db: uDB, rdb: rDB}
}

func (s *Service) MakeUUID() string {
	return uuid.New().String()
}

func (s *Service) DownloadFile(id string) (models.UploadedFiles, error) {
	return s.db.DownloadFile(id)
}

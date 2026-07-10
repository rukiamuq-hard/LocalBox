package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserDB interface {
	InsertInDB(login string, password string) error
	SearchInDB(login string) (string, error)
	GetIdFromLogin(login string) (int, error)
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

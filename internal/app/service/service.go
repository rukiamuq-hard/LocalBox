package service

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserDB interface {
	InsertInDB(login string, password string) error
	SearchInDB(login string) (string, error)
}

type Service struct {
	db UserDB
}

func New(uDB UserDB) *Service {
	return &Service{db: uDB}
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

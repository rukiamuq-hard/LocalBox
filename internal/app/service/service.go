package service

import (
	"golang.org/x/crypto/bcrypt"
)

type UserDB interface {
	InsertToDB(login string, password string) error
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

	err = s.db.InsertToDB(login, string(hashedPass))
	if err != nil {
		return err
	}
	return nil
}

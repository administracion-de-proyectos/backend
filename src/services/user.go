package services

import (
	"backend-admin-proyect/src/db"
	"fmt"
)

type userService struct {
	db db.DB[UserState]
}

func (s *userService) CheckCredentials(u UserState) (UserState, error) {
	if userCreated, err := s.db.Get(u.GetPrimaryKey()); err != nil || userCreated.Password != u.Password {
		return u, fmt.Errorf("wrong password or email")
	} else {
		return userCreated, nil
	}
}

func (s *userService) CreateUser(u UserState) error {
	if userCreated, err := s.db.Get(u.GetPrimaryKey()); err == nil {
		return fmt.Errorf("error with primary key: %s already created", userCreated.GetPrimaryKey())
	}
	s.db.Insert(u)
	return nil
}

func (s *userService) GetUser(userId string) (UserState, error) {
	return s.db.Get(userId)
}


func CreateUserService(db db.DB[UserState]) UserService {
	return &userService{db: db}
}

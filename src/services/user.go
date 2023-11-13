package services

import (
	"backend-admin-proyect/src/db"
	"fmt"
	"strings"
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

func (s *userService) UpdateUser(u UserState) (UserState, error) {
	s.db.Update(u)
	return s.db.Get(u.GetPrimaryKey())
}

type FilterValuesUser struct {
	Email   string
	Profile string
	Name    string
}

func (s *userService) FindUser(fv FilterValuesUser) []UserState {
	uss, _ := s.db.GetAll()
	response := make([]UserState, 0)
	for _, us := range uss {
		if us.isOkay(fv) {
			response = append(response, us)
		}
	}
	return response
}

func (us UserState) isOkay(fv FilterValuesUser) bool {
	ok := true
	ok = ok && strings.Contains(us.Email, fv.Email)
	ok = ok && strings.Contains(us.Name, fv.Name)
	ok = ok && strings.Contains(strings.ToLower(us.Profile), strings.ToLower(fv.Profile))
	return ok
}
func CreateUserService(db db.DB[UserState]) UserService {
	return &userService{db: db}
}

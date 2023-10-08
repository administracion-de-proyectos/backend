package services

type UserState struct {
	Email    string
	Password string
	Name     string
	Profile  string
}

func (us UserState) GetPrimaryKey() string {
	return us.Email
}

type UserService interface {
	CreateUser(u UserState) error
	CheckCredentials(u UserState) (UserState, error)
}

func CreateUserState(email, password, name, profile string) UserState {
	return UserState{
		email,
		password,
		name,
		profile,
	}
}

func CreateBasicUserState(email, password string) UserState {
	return UserState{
		Email:    email,
		Password: password,
	}
}

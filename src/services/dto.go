package services

import "fmt"

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
	GetUser(userId string) (UserState, error)
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

type CourseState struct {
	CreatorEmail string
	CourseTitle  string
	Classes      []string
}

func (cs CourseState) GetPrimaryKey() string {
	return cs.CourseTitle
}

type Class struct {
	Id          string
	CourseTitle string
	Metadata    interface{} // For any shit that the frontend wants to throw at us
}

func (cs Class) GetPrimaryKey() string {
	return getClassId(cs.CourseTitle, cs.Id)
}

func getClassId(courseTitle, classTitle string) string {
	return fmt.Sprintf("%s-%s", courseTitle, classTitle)
}

type CourseService interface {
	SetClassInPlaceN(n int, classTitle string, courseTitle string) CourseState
	AddClass(class Class, shouldEditCourse bool) (CourseState, error)
	AddCourse(course CourseState) CourseState
	GetCourse(courseId string) (CourseState, error)
	RemoveClass(courseId, classId string) error
	GetClass(courseId, classId string) (Class, error)
}

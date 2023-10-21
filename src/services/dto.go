package services

import (
	"fmt"
	"strings"
)

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
	UpdateUser(u UserState) (UserState, error)
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
	Category     string
	Metadata     interface{}
	AgeFiltered  bool `json:"age_filtered,omitempty"`
	MinAge       int  `json:"min_age,omitempty"`
	MaxAge       int  `json:"max_age,omitempty"`
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

type FilterValues struct {
	Title      string
	OwnerEmail string
	Category   string
	DesiredAge *int
}

func (cs CourseState) isOkayWithFilter(f FilterValues) bool {
	ok := strings.Contains(cs.CourseTitle, f.Title)
	ok = ok && strings.Contains(cs.CreatorEmail, f.OwnerEmail)
	ok = ok && strings.Contains(cs.Category, f.Category)
	ok = ok && strings.Contains(cs.Category, f.Category)
	ok = ok && (!(cs.AgeFiltered && f.DesiredAge != nil) || (cs.MinAge < *f.DesiredAge && cs.MaxAge > *f.DesiredAge))
	return ok
}

type CourseService interface {
	SetClassInPlaceN(n int, classTitle string, courseTitle string) CourseState
	AddClass(class Class, shouldEditCourse bool) (CourseState, error)
	AddCourse(course CourseState) CourseState
	GetCourse(courseId string) (CourseState, error)
	RemoveClass(courseId, classId string) error
	GetClass(courseId, classId string) (Class, error)
	GetCourses(values FilterValues) []CourseState
}

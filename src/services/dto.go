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
	Metadata interface{}
	HasPaid  bool `json:"has_paid"`
}

func (us UserState) GetPrimaryKey() string {
	return us.Email
}

func CreateUserState(email, password, name, profile string, metadata interface{}) UserState {
	return UserState{
		email,
		password,
		name,
		profile,
		metadata,
		false,
	}
}

func CreateBasicUserState(email, password string) UserState {
	return UserState{
		Email:    email,
		Password: password,
	}
}

type CourseState struct {
	CreatorEmail     string
	CourseTitle      string
	Classes          []string
	Category         string
	Metadata         interface{}
	AgeFiltered      bool `json:"age_filtered,omitempty"`
	MinAge           int  `json:"min_age,omitempty"`
	MaxAge           int  `json:"max_age,omitempty"`
	IsSchoolOriented bool
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
	Title            string
	OwnerEmail       string
	Category         string
	DesiredAge       *int
	IsSchoolOriented *bool
}

func (cs CourseState) isOkayWithFilter(f FilterValues) bool {
	ok := strings.Contains(cs.CourseTitle, f.Title)
	ok = ok && strings.Contains(cs.CreatorEmail, f.OwnerEmail)
	ok = ok && strings.Contains(cs.Category, f.Category)
	ok = ok && strings.Contains(cs.Category, f.Category)
	ok = ok && (!(cs.AgeFiltered && f.DesiredAge != nil) || (cs.MinAge < *f.DesiredAge && cs.MaxAge > *f.DesiredAge))
	ok = ok && (f.IsSchoolOriented == nil || *f.IsSchoolOriented == cs.IsSchoolOriented)
	return ok
}

type Subscription struct {
	UserId   string
	CourseId string
	Metadata interface{}
}

func (s Subscription) GetPrimaryKey() string {
	return s.UserId
}

func (s Subscription) GetSecondaryKey() string {
	return s.CourseId
}

type Point struct {
	Question string
	Options  []string
	Answer   string
}

type Exam struct {
	Points []Point `json:"points"`
	Class  string  `json:"class"`
	Course string  `json:"course"`
}

func (e Exam) GetPrimaryKey() string {
	return getClassId(e.Course, e.Class)
}

type StudentExam struct {
	StudentEmail string
	Course       string
	Class        string
	Points       []Point // Done like this because I am lazy, but this should be just question and answer
}
type Score struct {
	TotalAmount   int    `json:"total_amount"`
	CorrectAmount int    `json:"correct_amount"`
	Email         string `json:"email"`
	ClassId       string `json:"classId"`
}

func (s StudentExam) GetPrimaryKey() string {
	return s.StudentEmail
}

func (s StudentExam) GetSecondaryKey() string {
	return getClassId(s.Course, s.Class)
}

type Comment struct {
	CreatedAt  int64  `json:"created_at"`
	UserId     string `json:"user_id"`
	Commentary string `json:"comment"`
}

type Comments struct {
	CourseId string    `json:"course_id"`
	Data     []Comment `json:"comments"`
}

type Rate struct {
	CourseId  string `json:"course_id"`
	Score     int    `json:"score"`
	UserEmail string `json:"user_email"`
}

type GroupDto struct {
	OwnerEmail    string   `json:"owner_email"`
	StudentsGroup []string `json:"students_group"`
}

func (g GroupDto) GetPrimaryKey() string {
	return g.OwnerEmail
}

func (c Comments) GetPrimaryKey() string {
	return c.CourseId
}

func (r Rate) GetPrimaryKey() string {
	return r.UserEmail
}

func (r Rate) GetSecondaryKey() string {
	return r.CourseId
}

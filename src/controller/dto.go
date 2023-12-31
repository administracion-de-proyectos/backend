package controller

import "backend-admin-proyect/src/services"

type UserRequest struct {
	Email    string      `json:"email"`
	Password string      `json:"password,omitempty"`
	Name     string      `json:"name"`
	Profile  string      `json:"profile"`
	Metadata interface{} `json:"metadata"`
}

type UserResponse struct {
	Email    string      `json:"email"`
	Name     string      `json:"name"`
	Profile  string      `json:"profile"`
	Metadata interface{} `json:"metadata"`
	HasPaid  bool        `json:"has_paid"`
}

type Token struct {
	Token string `json:"token" example:"asdasfasd"`
}

type ErrorMsg struct {
	Reason string `json:"reason" example:"mensaje de error"`
}

type CourseRequest struct {
	Title            string      `json:"title"`
	Classes          []Class     `json:"classes"`
	Metadata         interface{} `json:"metadata"`
	Category         string      `json:"category"`
	MinAge           *int        `json:"min_age,omitempty"`
	MaxAge           *int        `json:"max_age,omitempty"`
	IsSchoolOriented bool        `json:"is_school_oriented"`
}

type Class struct {
	Title       string      `json:"title"`
	CourseTitle string      `json:"course_title"`
	Metadata    interface{} `json:"metadata"`
}

type SubscriptionRequest struct {
	UserId      string      `json:"user_id"`
	CourseTitle string      `json:"course_title"`
	Metadata    interface{} `json:"metadata"`
}

type Point struct {
	Question      string   `json:"question"`
	Answer        string   `json:"answer"`
	Possibilities []string `json:"possibilities"`
}

type CreateExamRequest struct {
	Points []Point `json:"points" binding:"required"`
}

type SubmissionPoint struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type Submission struct {
	Course string            `json:"course" binding:"required"`
	Class  string            `json:"class" binding:"required"`
	Points []SubmissionPoint `json:"points" binding:"required"`
}

type CommentRequest struct {
	Course  string `json:"course" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}

type RateDTO struct {
	Course string `json:"course" binding:"required"`
	Rate   int    `json:"rate" binding:"required"`
}

// CourseState Only for docs
type CourseState struct {
	CreatorEmail     string      `json:"creatorEmail"`
	CourseTitle      string      `json:"courseTitle"`
	Classes          []string    `json:"classes"`
	Category         string      `json:"category"`
	Metadata         interface{} `json:"metadata"`
	AgeFiltered      bool        `json:"age_filtered,omitempty"`
	MinAge           int         `json:"min_age,omitempty"`
	MaxAge           int         `json:"max_age,omitempty"`
	IsSchoolOriented bool        `json:"isSchoolOriented"`
	IsSubscribed     bool        `json:"isSubscribed"`
}

type CourseStateResponse struct {
	Courses []CourseState
	Amount  int
}

type Exam struct {
	Points []Point `json:"points"`
	Class  string  `json:"class"`
	Course string  `json:"course"`
}

type Score struct {
	TotalAmount   int    `json:"total_amount"`
	CorrectAmount int    `json:"correct_amount"`
	Email         string `json:"email"`
	CourseId      string `json:"course_id"`
}

type Comment struct {
	CreatedAt  int    `json:"created_at"`
	UserId     string `json:"user_id"`
	Commentary string `json:"comment"`
}

type Comments struct {
	CourseId string    `json:"course_id"`
	Data     []Comment `json:"comments"`
}

type RateResponse struct {
	CourseId string          `json:"course_id"`
	RateAvg  float64         `json:"rate_avg"`
	RateArr  []services.Rate `json:"rate_arr"`
}

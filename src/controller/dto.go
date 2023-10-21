package controller

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

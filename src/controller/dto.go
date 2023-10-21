package controller

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Profile  string `json:"profile"`
}

type UserResponse struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
}

type Token struct {
	Token string `json:"token" example:"asdasfasd"`
}

type ErrorMsg struct {
	Reason string `json:"reason" example:"mensaje de error"`
}

type CourseRequest struct {
	Title    string      `json:"title"`
	Classes  []Class     `json:"classes"`
	Metadata interface{} `json:"metadata"`
	Category string      `json:"category"`
	MinAge   *int        `json:"min_age,omitempty"`
	MaxAge   *int        `json:"max_age,omitempty"`
}

type Class struct {
	Title       string      `json:"title"`
	CourseTitle string      `json:"course_title"`
	Metadata    interface{} `json:"metadata"`
}

// CourseState Only for docs
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

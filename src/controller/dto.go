package controller

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Profile  string `json:"profile"`
}

type Token struct {
	Token string `json:"token" example:"asdasfasd"`
}

type ErrorMsg struct {
	Reason string `json:"reason" example:"mensaje de error"`
}

type CourseRequest struct {
	Title   string  `json:"title"`
	Classes []Class `json:"classes"`
}

type Class struct {
	Title       string      `json:"title"`
	CourseTitle string      `json:"course_title"`
	Metadata    interface{} `json:"metadata"`
}

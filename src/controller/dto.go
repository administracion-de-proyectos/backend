package controller

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
	Profile  string `json:"profile"`
}

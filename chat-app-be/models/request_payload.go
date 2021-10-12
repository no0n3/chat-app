package models

type ChangeProfileImagePayload struct {
	MediaId string
}

type LoginPayload struct {
	Email    string
	Password string
}

type SignupPayload struct {
	Email    string
	Password string
	Name     string
	Username string
}

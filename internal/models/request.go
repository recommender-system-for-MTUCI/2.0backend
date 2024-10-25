package models

type RequestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

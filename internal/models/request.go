package models

type RequestLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type RequestRegister struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type RequestAcceptEmail struct {
	Code int `json:"code"`
}
type RequestChangePassword struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

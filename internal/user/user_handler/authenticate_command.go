package user_handler

type AuthenticateCommand struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

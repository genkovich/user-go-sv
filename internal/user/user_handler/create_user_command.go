package user_handler

type CreateUserCommand struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

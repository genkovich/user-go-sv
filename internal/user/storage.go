package user

type Storage interface {
	GetList() []User
	Add(user User)
	Remove(userId string)
}

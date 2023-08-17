package user

type Storage interface {
	GetList(limit int, offset int) ([]User, error)
	Add(user User)
	Remove(userId string)
}

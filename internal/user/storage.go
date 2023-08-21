package user

type Storage interface {
	GetList(limit int, offset int) ([]User, error)
	Add(user User) ([]User, error)
	Remove(userId string)
	GetByLogin(login string) (*User, error)
	UpdatePassword(userId string, newPassword string) error
	UpdateUserRole(userId string, newRole string) error
}

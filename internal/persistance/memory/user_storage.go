package memory

import "user-service/internal/user"

type Storage struct {
	users map[string]user.User
}

func NewUserMemoryStorage() *Storage {
	testUser := user.CreateUser("test", "password")
	secondUser := user.CreateUser("second", "password")

	return &Storage{
		users: map[string]user.User{
			testUser.GetId().String():   *testUser,
			secondUser.GetId().String(): *secondUser,
		},
	}
}

func (ms *Storage) GetList() []user.User {
	var users []user.User

	for _, userEntity := range ms.users {
		users = append(users, userEntity)
	}

	return users
}

func (ms *Storage) Add(user user.User) {
	ms.users[user.GetId().String()] = user
}

func (ms *Storage) Remove(userId string) {
	delete(ms.users, userId)
}

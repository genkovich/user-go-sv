package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"time"
	"user-service/internal/user"
	"user-service/pkg/database"
)

type userFields struct {
	Id           string
	Login        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Storage struct {
	logger     *zap.Logger
	connection *pgxpool.Pool
}

func NewUserPsqlStorage(connector database.Connector, logger *zap.Logger) *Storage {
	return &Storage{
		connection: connector.GetConnection(),
		logger:     logger,
	}
}

func (s *Storage) GetList(limit int, offset int) ([]user.User, error) {
	rows, err := s.connection.Query(context.Background(), "SELECT * FROM users LIMIT $1 OFFSET $2", limit, offset)

	if err != nil {
		s.logger.Error("failed get list users", zap.Error(err))
		return nil, err
	}

	var userList []user.User

	for rows.Next() {
		var u userFields

		if err := rows.Scan(&u.Id, &u.Login, &u.Role, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			s.logger.Error("failed to scan row", zap.Error(err))
			return nil, err
		}

		tempUser, err := user.MapFromData(u.Id, u.Login, u.PasswordHash, u.Role)
		if err != nil {
			s.logger.Error("failed mapping users", zap.Error(err))
			return nil, err
		}

		userList = append(userList, *tempUser)
	}

	return userList, nil

}

func (s *Storage) Add(user user.User) ([]user.User, error) {
	roleStr := user.GetRole().String()

	_, err := s.connection.Exec(context.Background(),
		"INSERT INTO users (id, login, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.GetId(), user.GetLogin(), user.GetHashPassword(), roleStr, time.Now(), time.Now())

	if err != nil {
		s.logger.Error("failed to insert user into the database", zap.Error(err))
		return nil, err
	}

	userList, err := s.GetList(10, 0)
	if err != nil {
		s.logger.Error("failed to fetch updated user list", zap.Error(err))
		return nil, err
	}

	return userList, nil
}

func (s *Storage) Remove(userId string) {

}

func (s *Storage) GetByLogin(login string) (*user.User, error) {
	var u userFields

	err := s.connection.QueryRow(context.Background(), "SELECT * FROM users WHERE login = $1", login).
		Scan(&u.Id, &u.Login, &u.Role, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		s.logger.Error("failed to scan row", zap.Error(err))
		return nil, err
	}

	tempUser, err := user.MapFromData(u.Id, u.Login, u.PasswordHash, u.Role)
	if err != nil {
		s.logger.Error("failed mapping users", zap.Error(err))
		return nil, err
	}

	return tempUser, nil

}

func (s *Storage) UpdatePassword(userId string, newPassword string) error {
	hashedPassword, err := s.hashAndSalt([]byte(newPassword))
	if err != nil {
		s.logger.Error("password hashing error", zap.Error(err))
		return err
	}

	_, err = s.connection.Exec(context.Background(),
		"UPDATE users SET password_hash = $1 WHERE id = $2",
		hashedPassword, userId)

	if err != nil {
		s.logger.Error("failed to update password in the database", zap.Error(err))
		return err
	}

	return nil
}

func (s *Storage) UpdateUserRole(userId string, newRole string) error {
	allowedRoles := map[string]string{
		"ROLE_USER":  "ROLE_USER",
		"ROLE_ADMIN": "ROLE_ADMIN",
	}

	if _, ok := allowedRoles[newRole]; !ok {
		return fmt.Errorf("invalid role: %s", newRole)
	}

	_, err := s.connection.Exec(context.Background(),
		"UPDATE users SET role = $1 WHERE id = $2",
		newRole, userId)

	if err != nil {
		s.logger.Error("failed to update role in the database", zap.Error(err))
		return err
	}

	return nil
}

func (s *Storage) hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

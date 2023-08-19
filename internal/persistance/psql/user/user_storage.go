package user

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
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

func (s *Storage) Add(user user.User) {
	_, err := s.connection.Exec(context.Background(),
		"INSERT INTO users (id, login, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
		user.GetId(), user.GetLogin(), user.GetHashPassword(), user.GetRole(), time.Now(), time.Now())

	if err != nil {
		s.logger.Error("failed to insert user into the database", zap.Error(err))
	}

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

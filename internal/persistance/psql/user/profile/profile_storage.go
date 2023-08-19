package profile

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
	"user-service/internal/user/profile"
	"user-service/pkg/database"
)

type profileFields struct {
	Id        string
	UserID    string
	FirstName string
	LastName  string
	DOB       time.Time
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Storage struct {
	logger     *zap.Logger
	connection *pgxpool.Pool
}

func NewProfilePsqlStorage(connector database.Connector, logger *zap.Logger) *Storage {
	return &Storage{
		connection: connector.GetConnection(),
		logger:     logger,
	}
}

func (s *Storage) GetByUserId(userId uuid.UUID) (*profile.Profile, error) {
	var p profileFields

	err := s.connection.QueryRow(context.Background(), "SELECT * FROM profiles WHERE user_id = $1", userId).
		Scan(&p.Id, &p.UserID, &p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.DOB, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			s.logger.Error("failed get profile by user id", zap.Error(err))
			return nil, err
		}
	}

	profileEntity, err := profile.MapFromData(uuid.MustParse(p.Id), p.FirstName, p.LastName, p.DOB, p.Email, p.Phone)

	if err != nil {
		s.logger.Error("failed mapping profile", zap.Error(err))
		return nil, err
	}

	return profileEntity, nil
}

func (s *Storage) Upsert(userId uuid.UUID, profile *profile.Profile) error {
	now := time.Now()

	_, err := s.connection.Exec(context.Background(), "INSERT INTO profiles (id, user_id, first_name, last_name, dob, email, phone, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT (user_id) DO UPDATE SET first_name = $3, last_name = $4, dob = $5, email = $6, phone = $7, updated_at = $9",
		profile.GetId(), userId, profile.GetFirstName(), profile.GetLastName(), profile.GetDob(), profile.GetEmail(), profile.GetPhone(), now, now)

	if err != nil {
		s.logger.Error("failed upsert profile", zap.Error(err))
		return err
	}

	return nil
}

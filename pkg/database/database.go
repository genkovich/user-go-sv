package database

import (
	"context"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	connection *pgxpool.Pool
}

func NewConnection(dsn string) (*Database, error) {
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	return &Database{connection: conn}, nil
}

func (d *Database) GetConnection() *pgxpool.Pool {
	return d.connection
}

type Connector interface {
	GetConnection() *pgxpool.Pool
}

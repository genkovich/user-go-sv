package main

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PsqlDSN string
}

func ReadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	psqlDsn := os.Getenv("PSQL_DSN")

	return &Config{PsqlDSN: psqlDsn}, nil

}

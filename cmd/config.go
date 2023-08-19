package main

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	PsqlDSN   string
	JwtSecret []byte
	RedisAddr string
}

func readConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, err
	}

	psqlDsn := os.Getenv("PSQL_DSN")
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	redisAddr := os.Getenv("REDIS_ADDR")

	return &Config{
		PsqlDSN:   psqlDsn,
		JwtSecret: jwtSecret,
		RedisAddr: redisAddr,
	}, nil

}
